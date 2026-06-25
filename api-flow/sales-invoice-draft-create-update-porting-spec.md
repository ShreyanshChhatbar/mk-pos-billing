# Sales Invoice Draft Create/Update Flow

Endpoints:
- `POST /api/v1/sales/sales-invoice/draft/create`
- `PUT /api/v1/sales/sales-invoice/draft/{id}`

Scope:
- Supports only the default Medkart POS organization flow.
- Server forces `organization_id = DEFAULT_ORGANIZATION_ID`.
- Client `organization_id`, `store_id`, `user_id`, `device_master_id`, `till_id`, and `till_transaction_id` are not trusted.
- ONGC/external-organization prescription, substitution, and organization-credit flows are excluded.

Flow
1. Validate `device` header against Laravel-compatible device cache:
   - remember key: `device_remember_token_`
   - detail key: `device_token_{token}`
   - reject missing/invalid/expired/inactive device
   - set `device_master_id` from cache
   - require `store` header unless device type is `TAB`
2. Validate `Authorization: Bearer <token>` against POS auth cache:
   - token-to-user key: `pos_auth_user_cache_{token}`
   - user-detail key: `pos_auth_token_{user_id}`
   - reject missing/invalid user or empty permissions
   - set `user_id` from cache
3. Map store/till context:
   - parse `store` header as positive integer
   - set `store_id` from header
   - load till cache `till_cache_data_{store_id}` under tag `till_cache_tags`
   - set `till_id` when present
4. Check store POS permissions:
   - load store cache with prefix `STORE_CACHE_KEY_`
   - reject invalid store
   - check module `sales`, submodule `sales-invoice`
   - require action `CREATE` for POST and `UPDATE` for PUT
5. Check till status:
   - reject missing till with `400 Please Generate Till Number for this store till`
   - reject closed/missing transaction with `307 Please open the till`
   - allow `PUT /draft/{id}` when the open till transaction is from a previous day, matching Laravel
   - set `till_transaction_id` when the till transaction is open
6. Run duplicate request check on both POST and PUT:
   - key: `{user_id}_duplicate_check_cache_key__sales_sales-invoice`
   - value: SHA-256 of request body plus middleware context fields
   - TTL: `DUPLICATE_REQUEST_CHECK_EXPIRY_IN_SECONDS`, default `10`
   - reject same key/hash with `429 Duplicate Request Received`
7. Bind request body and route id:
   - `PUT` id must be a positive number
   - `items` omitted is distinct from `items: []`
   - payment `id` marks an existing draft payment and is not reinserted
8. Validate default-org draft rules:
   - `is_confirmed` required
   - `billing_user_id` required and active
   - `customer_id` required when confirmed or payments exist
   - `patient_id` required when confirmed and must belong to customer
   - `doctor_id` required when payments exist
   - home-delivery payment requires valid customer address with pin code
   - item product/batch/quantity rules include sales-unit, MSP, delivery-charge, and max batch-code length
   - H1 narcotics confirmation requires `course_days`
   - payment method must be active for store and default organization
9. Create or load draft:
   - draft cache key: `draft_invoices_{draft_id}`
   - draft cache tag: `draft_invoices_tags_{store_id}`
   - update loads only `DRAFT` or `PAYMENT_PENDING` draft for same store/default organization
10. If update omits `items`:
   - clear draft products from `draft_json`
   - reset product/tax/amount totals
   - append any new payments
   - preserve existing non-deleted draft payments and payment received total
   - return updated draft
11. If items are present:
   - fetch batch metadata from `batches`
   - fetch stock from `store_inventories` + `store_batches` + `batches`
   - require `batches.expiry_date >= now + DEFAULT_MINIMUM_DAYS_BATCH_EXPIRY`
   - reject insufficient stock
   - reject sales rate above MRP
   - calculate line totals, discount, invoice total, and round-off
12. Persist draft:
   - save `sales_invoice_draft_jsons`
   - append only new `sales_invoice_draft_payments`
   - do not soft-delete existing draft payments during create/update
   - set payment status from total active non-advance-refund payments
13. If `is_confirmed=false`:
   - cache draft result for 30 minutes
   - return draft response with items and active payments
14. If `is_confirmed=true`:
   - require payment total equals invoice total plus prepaid amount unless an `ADVANCE_REFUND` draft payment exists
   - create `sales_invoices`
   - create `sales_invoice_details`
   - copy draft payments into `sales_invoice_payments`
   - decrement `store_inventories` by exact `store_batch_id`
   - mark draft status `INVOICED`
   - cache final draft state
   - return created invoice id

Failure Branches
1. Missing/invalid device returns `401`.
2. Missing/invalid bearer token or permissions returns `401`.
3. Missing/invalid store returns `400`.
4. Missing till returns `400`.
5. Closed/stale till returns `307` except draft update stale-day bypass.
6. Duplicate request returns `429`.
7. Validation, stock, batch, product, payment, or draft lookup failures return `400`.

Tables Touched
- `sales_invoice_draft_jsons`
- `sales_invoice_draft_payments`
- `sales_invoices`
- `sales_invoice_details`
- `sales_invoice_payments`
- `store_inventories`
- `store_batches`
- `batches`
- validation lookup tables: `users`, `customers`, `patients`, `doctors`, `customer_addresses`, `products`, `store_payment_methods`

Verification
- Run `GOCACHE=/tmp/go-build go test ./...`
- Manually verify with Redis cache data shaped like Laravel `pos-store-api`
- Compare draft/payment/invoice rows against Laravel for default organization only
