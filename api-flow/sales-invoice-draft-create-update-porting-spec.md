# Sales Invoice Draft Create/Update Porting Spec

This document captures full behavior parity for these two routes so they can be rebuilt in another stack (for example Go/Gin) without changing business logic:

- `POST /api/v1/sales/sales-invoice/draft/create`
- `PUT /api/v1/sales/sales-invoice/draft/{id}`

Both call:
- `SalesInvoiceController::createOrUpdate(SalesInvoiceRequest $request, $id = null)`

## 1. Route Contract

### 1.1 Route definitions
- `Route::post('create', [SalesInvoiceController::class, 'createOrUpdate'])->name('create');`
- `Route::put('{id}', [SalesInvoiceController::class, 'createOrUpdate'])->name('update')->whereNumber('id');`
- Group prefix path: `/api/v1/sales/sales-invoice/draft`
- Effective full endpoints:
  - `POST /api/v1/sales/sales-invoice/draft/create`
  - `PUT /api/v1/sales/sales-invoice/draft/{id}`

### 1.2 Middleware chain affecting behavior
These routes are inside global API groups plus `duplicate-check` route group:
- `device-authentication`
- `user-auth`
- `map-store-till`
- `store-pos-permissions-check`
- `pagination`
- `check-till-status`
- `duplicate-check`

The business logic depends on middleware-populated context fields like:
- `user_id`
- `store_id`
- `till_id`
- `till_transaction_id`
- `device_master_id`

## 2. High-Level Behavior

Both endpoints are unified into one flow:
- Create (`POST create`, `id = null`): create a new draft invoice.
- Update (`PUT {id}`): update an existing draft invoice.
- If `is_confirmed = false`: save/update draft only.
- If `is_confirmed = true`: finalize draft into invoice (stock movement, invoice rows, payment rows, status change).

Controller is thin and delegates to repository; repository delegates to `SalesInvoiceDraftService::createOrUpdate`.

## 3. Request Validation (Exact Laravel Semantics)

Validation is route-aware (`SalesInvoiceRequest::rules`) and applies this block only when route is:
- `sales.sales-invoice.draft.create`
- `sales.sales-invoice.draft.update`

Before validation, request is mutated:
- `id` = route param `id`
- `organization_id` = `config('mk-psa.default_organization_id')`

### 3.1 Core fields and conditional required rules
- `customer_id`
  - required if `is_confirmed=true` OR payments exist
  - nullable otherwise
  - must exist in `customers` active and not soft-deleted
- `billing_user_id`
  - required
  - must exist in active `users`
- `patient_id`
  - required if `is_confirmed=true`
  - patient must belong to `customer_id`, active, not deleted
- `customer_address_id`
  - required when payments exist AND `is_home_delivery=true`
  - must belong to customer, active, not deleted
  - when home delivery, address must have `pin_code`
- `doctor_id`
  - required when payments exist
  - must exist in active `doctors`
- `is_home_delivery`
  - boolean

### 3.2 Items validation
- `items`
  - array
  - required if `is_confirmed=true`
  - `LoyaltyProductValidation` applied
- `items.*.product_id`
  - required integer > 0
  - product must exist and be active
- `items.*.batch_code`
  - required string
  - max length `config('mk-psa.batch_code_length')`
  - must exist in `batches.batch_number`
- `items.*.quantity`
  - required integer > 0
  - `CheckQuantity(route_id, store_id)`
  - additional custom checks:
    - must be multiple of product `sales_unit`
    - MSP product quantity must be `1`
    - delivery-charge product quantity must be `1`

### 3.3 Special medical rule
- `course_days`
  - required numeric > 0 when:
    - `is_confirmed=true`
    - and any item is H1 narcotics (schedule type check using cached product data)

### 3.4 Promo and payments
- `promo_code`
  - nullable string
  - `IsPromoCodeApplicable(store_id)`
- `payments`
  - array
  - conditional required when confirmed, with special exception for advance/drop shipping with zero due amount
- `payments.*.store_payment_method_id`
  - required integer > 0
  - must exist in active `store_payment_methods` for this store/org
  - `StorePaymentMethodValidation` rule
- `payments.*.amount`
  - required numeric > 0
- `payments.*.voucher_code`
  - nullable string
  - custom `IsVoucherCodeApplied` check
- `payments.*.is_advance_refund`
  - optional boolean
  - custom `AdvanceToSalesRefundValidation`

### 3.5 Confirmation and ASM passkey
- `is_confirmed`: required boolean
- `asm_user_id`: optional integer user id
  - at repository layer, if provided, ASM passkey is validated (`validateASMApprovePermission`), otherwise 400 `Invalid passkey!`

### 3.6 Validation custom messages
Draft create/update messages include custom text for:
- item product validation failures
- payment method/amount validation failures

## 4. Request Payload Structure (Porting DTO)

Use this JSON shape for both endpoints.

```json
{
  "billing_user_id": 123,
  "customer_id": 456,
  "patient_id": 789,
  "doctor_id": 999,
  "customer_address_id": 111,
  "is_home_delivery": false,
  "is_confirmed": false,
  "promo_code": null,
  "course_days": 7,
  "notes": "optional",
  "asm_user_id": null,
  "items": [
    {
      "product_id": 1001,
      "batch_code": "BATCH123",
      "quantity": 2,
      "is_free_product": false,
      "combo_product_id": null,
      "sales_rate": 120.5,
      "priceDelta": null,
      "best_alternate": {},
      "per_tab_frontend": null
    }
  ],
  "payments": [
    {
      "store_payment_method_id": 55,
      "amount": 241.0,
      "voucher_code": null,
      "is_advance_refund": false
    }
  ]
}
```

### 4.1 Context fields required from middleware/request context
Even if not explicitly validated in this rule block, logic expects:
- `store_id`
- `user_id`
- `device_master_id`
- `till_id`
- `till_transaction_id`
- `organization_id` (forced to default org in request)

## 5. Core Business Flow (Step-by-Step)

### 5.1 Controller -> Repository
- Controller calls repository `createOrUpdate($request, $id)`.
- Repository opens transaction and calls `SalesInvoiceDraftService::createOrUpdate(...)`.

### 5.2 Draft service orchestration
`SalesInvoiceDraftService::createOrUpdate(request, id, isConfirmed, ignoreStockCheck)`:
1. Read/seed draft cache (`draft_invoices_{id}` tagged by `draft_invoices_tags_{store_id}`).
2. Initialize draft model:
   - from cache if available
   - else DB by id for update
   - else new draft model
3. Populate basic draft fields and save.
4. If update request has no `items`, clear draft products from draft JSON.
5. Process free products and combo products (expands `items` list when applicable).
6. Fetch batch stocks (`SalesInvoiceService::getAllBatches`).
7. Calculate pricing/tax/totals/product lines (`SalesInvoiceService::calculateSalesInvoice`).
8. Persist draft line items into draft JSON products.
9. Store draft payments and update payment status.
10. Validate prescription doctor rule for prescription-required items.
11. If confirmed:
   - ensure payment integrity (`total_invoice_amount == total_amount_received + prepaid_amount`) unless advance refund logic applies
   - create final invoice + details + tax + payments
   - move inventory stock
   - set draft status to invoiced
12. Return success payload (`success=true,data=...`) to repository.

### 5.3 Pricing/stock calculation highlights (SalesInvoiceService)
- Batch fetch with expiry cutoff (`default_minimum_days_for_batch`).
- Promo code minimum-order validation (`validatePromoCode`).
- Stock checks for non-service products.
- Combo product stock checks and possible expansion logic.
- Delivery-charge product enforcement.
- Free combo minimum order value constraints.
- Narcotics conflict checks for customer.
- Computation outputs:
  - product-level rates/discount/gst fields
  - invoice-level totals: bill amount, total amount, taxes, quantity, round-off, invoice amount

### 5.4 Confirm-only side effects
On `is_confirmed=true`:
- Insert into `sales_invoices`
- Insert into `sales_invoice_details`
- Insert into `sales_invoice_tax_details`
- Insert into `sales_invoice_payments`
- Update voucher status when voucher used
- Move stock via inventory transaction entries
- Create COD/Credit payment rows when payment method flags require it
- Update draft status to `INVOICED`
- Async job: `ProcessCustomerMetaDataAfterInvoice`

## 6. Data Stores / Tables Touched

At minimum across draft+confirm path:
- `sales_invoice_draft_jsons`
- `sales_invoice_draft_payments`
- `sales_invoices`
- `sales_invoice_details`
- `sales_invoice_tax_details`
- `sales_invoice_payments`
- `store_inventories`
- `store_batches`
- `batches`
- plus voucher/loyalty/COD/credit related tables based on payment methods

## 7. Cache Structures (Must Preserve for Parity)

### 7.1 Duplicate request cache (middleware-level)
Purpose: prevent replay/rapid duplicate POST/PUT payloads.

- Key format:
  - `{user_id}_{duplicate_check_cache_key}_{segment3}_{segment4}`
  - with defaults this becomes like: `{user_id}_duplicate_check_cache_key_{sales}_{sales-invoice}`
- Value: `sha256(json_encode(request_body))`
- TTL: `config('mk-psa.duplicate_request_check_expiry_in_seconds')` (default 10s)
- Behavior:
  - if same key and same hash exists: return `429 Duplicate Request Received`
  - on validation failure/exception: key is cleared so caller can retry

### 7.2 Draft invoice computation cache (service-level)
Used to accelerate repeated create/update calculations for same draft/store.

- Tag format:
  - `{prefix_draft_bill_cache_tags}{store_id}`
  - default: `draft_invoices_tags_{store_id}`
- Key format:
  - `{prefix_draft_bill_cache}{draft_id}`
  - default: `draft_invoices_{draft_id}`
- TTL:
  - 30 minutes (`SalesInvoiceDraftService::CACHE_TTL_MINUTES`)
- Value structure:

```json
{
  "salesInvoiceDraft": "SalesInvoiceDraftJson model snapshot",
  "calculatedProductData": {
    "products": {
      "{product_id}_{batch_code}": {
        "product_id": 1001,
        "batch_code": "BATCH123",
        "sales_rate": 120.5,
        "...": "computed product fields"
      }
    }
  }
}
```

- Cache updates happen:
  - immediately after draft initialization
  - after recalculation
  - after confirm (cached model becomes finalized invoice model)

### 7.3 Config keys affecting cache behavior
- `mk-psa.duplicate_check_cache_key` (default `duplicate_check_cache_key_`)
- `mk-psa.duplicate_request_check_expiry_in_seconds` (default `10`)
- `mk-psa.prefix_draft_bill_cache` (default `draft_invoices_`)
- `mk-psa.prefix_draft_bill_cache_tags` (default `draft_invoices_tags_`)

## 8. Response Contract

### 8.1 Success envelope (repository response)

```json
{
  "code": 200,
  "data": {},
  "message": "...",
  "meta": {}
}
```

For these two routes, `meta` is typically absent.

### 8.2 Success cases for create/update

#### A) Draft save success (`is_confirmed=false`)
- `message`:
  - create: `Sales Invoice (Draft) Created Successfully`
  - update: `Sales Invoice (Draft) Updated Successfully`
- `data`: `SalesInvoiceDraftDetailsResource`

Important top-level fields from this resource:
- `id`
- `organization_id`
- `prescription_id`
- `invoice_number`
- `date_and_timestamp`, `timestamp`
- `is_home_delivery`
- `total_products`, `total_items`, `total_quantity`
- `prepaid_amount`
- `total_invoice_amount`
- `total_amount_received`
- `amount_due`
- `delivery_charges`
- `total_mrp`
- `total_savings`
- `taxable_amount`
- `status`, `payment_status`
- `round_off`
- `store`, `billing_user`, `customer`, `customer_address`, `doctor`, `patient`
- `applicable_discount`
- `applicable_taxes`
- `total_bill_amount_before_promo`
- `promo_code`
- `items[]` (`SalesInvoiceItemResource`)
- `payments[]` (`SalesInvoicePaymentResource`)
- `is_advance_order`, `is_dropship_order`

#### B) Confirmed invoice success (`is_confirmed=true`)
- `message`: `Sales Invoice Created Successfully`
- `data`: integer invoice id (not full invoice resource)

### 8.3 Error envelopes

#### Validation failure (`BaseRequest`)
```json
{
  "code": 400,
  "type": "Bad Request",
  "message": "Validation Error",
  "errors": {"field": ["message"]}
}
```

#### Duplicate request
```json
{
  "code": 429,
  "type": "Not Acceptable",
  "message": "Duplicate Request Received"
}
```

#### Business rule failure (thrown as 400)
- Typical message directly in `message`, OR
- If JSON error map was thrown: `message = "Validation Error"` and `errors` object present.

#### Server failure
```json
{
  "code": 500,
  "data": [],
  "message": "Something went wrong"
}
```

## 9. Create vs Update Differences

- `POST create`
  - `id` route param absent
  - new draft instance initialized
  - success message says Created
- `PUT {id}`
  - `id` must be numeric route param
  - existing draft loaded with store/org/status constraints
  - if `items` omitted: draft products are cleared
  - success message says Updated

Everything else (validation logic, computation, cache keys, confirm flow) is shared.

## 10. Go/Gin Porting Checklist (No Logic Drift)

Recreate exactly:
1. Keep a single service function handling both create/update.
2. Preserve conditional validation semantics (`is_confirmed`, payments/address/doctor/patient dependencies).
3. Preserve duplicate-request cache key format + TTL + clear-on-failure behavior.
4. Preserve draft-calculation cache key/tag structure + 30-minute TTL + cached object shape.
5. Preserve same payment integrity check before confirming invoice.
6. Preserve same data writes and order inside one DB transaction.
7. Preserve confirm-only side effects (stock movement, payment rows, status transitions, async customer metadata job equivalent).
8. Preserve response shape differences:
  - draft mode returns detailed draft resource
  - confirmed mode returns invoice id only

## 11. Notes for Implementers

- `organization_id` is forcibly merged to default org in request rules for these routes.
- Internal calculations rely on product/store cache lookups (`CacheMaster`) and DB inventory joins; parity implementation should keep equivalent source-of-truth order.
- Some rules/flows branch for external organizations and ONGC flows; for these exact routes default org path is primary.
- Route names drive validation branching in Laravel. In Go/Gin, represent this by binding logic to endpoint handlers explicitly.
