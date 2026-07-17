package cache

type BankMaster struct {
	ID            int    `php:"id"`
	BankName      string `php:"bank_name"`
	AccountNumber string `php:"account_number"`
	IfscCode      string `php:"ifsc_code"`
}

type StoreSettings struct {
	FloatCashAmount                  float64      `php:"float_cash_amount"`
	TillMismatchTaskTrigger          bool         `php:"till_mismatch_task_trigger"`
	TillMismatchTaskTriggerMinAmount float64      `php:"till_mismatch_task_trigger_min_amount"`
	IsOtpRequiredSalesReturn         bool         `php:"is_otp_required_sales_return"`
	IsDeliveryChargeApplicable       bool         `php:"is_delivery_charge_applicable"`
	DeliveryChargeMinimumMov         float64      `php:"delivery_charge_minimum_mov"`
	TransferInMinimumAmount          float64      `php:"transfer_in_minimum_amount"`
	TransferOutMinimumAmount         float64      `php:"transfer_out_minimum_amount"`
	IsUrgentOrderEnabled             bool         `php:"is_urgent_order_enabled"`
	IsAdvanceOrderEnabled            bool         `php:"is_advance_order_enabled"`
	IsMinMaxOrderEnabled             bool         `php:"is_min_max_order_enabled"`
	IsMinMaxDropshipOrderEnabled     bool         `php:"is_min_max_dropship_order_enabled"`
	BankMasters                      []BankMaster `php:"bank_masters"`
	UrgentOrderTime                  *string      `php:"urgent_order_time"`
}

type ProductCategoryDiscount struct {
	ID                   int     `php:"id"`
	B2CPricingTemplateID int     `php:"b_2_c_pricing_template_id"`
	CategoryID           int     `php:"category_id"`
	PricingCategory      string  `php:"pricing_category"`
	Mode                 string  `php:"mode"`
	Value                float64 `php:"value"`
	Operator             string  `php:"operator"`
	// PromoCodes           []string `php:"promo_codes"` // Optional based on mapping
}

type B2BProductCategoryDiscount struct {
	ID                     int    `php:"id"`
	B2BPricingTemplateID   int    `php:"b_2_b_pricing_template_id"`
	B2BProductCategoryName string `php:"b_2_b_product_category_name"`
	B2BPricingCategory     string `php:"b_2_b_pricing_category"`
	Mode                   string `php:"mode"`
	Value                  string `php:"value"`
	Operator               string `php:"operator"`
}

type OrgCategoryDiscount struct {
	CategoryName    string `php:"category_name"`
	Mode            string `php:"mode"`
	Operator        string `php:"operator"`
	PricingCategory string `php:"pricing_category"`
	Value           string `php:"value"`
}

type Organization struct {
	Name                  string  `php:"name"`
	OrganizationCode      string  `php:"organization_code"`
	IsCreditSystemEnabled bool    `php:"is_credit_system_enabled"`
	CreditLimit           float64 `php:"credit_limit"`
}

type StoreDrugLicense struct {
	ID             int    `php:"id"`
	LicenseNumber  string `php:"license_number"`
	LicenseTagName string `php:"license_tag_name"`
}

type StoreManager struct {
	ID                 int    `php:"id"`
	ManagerName        string `php:"manager_name"`
	ManagerEmail       string `php:"manager_email"`
	ManagerContact     string `php:"manager_contact"`
	PrimaryContactName string `php:"primary_contact_name"`
	StorePhoneNumber   string `php:"store_phone_number"`
}

type Store struct {
	ID                                     int                                   `php:"id"`
	StoreName                              string                                `php:"store_name"`
	StoreType                              string                                `php:"store_type"`
	StoreClusters                          []int                                 `php:"store_clusters"`
	GstNumber                              string                                `php:"gst_number"`
	IsReReleased                           bool                                  `php:"is_re_released"`
	StoreGroup                             string                                `php:"store_group"`
	PlaceOfSupplyCode                      string                                `php:"place_of_supply_code"`
	IsCustomLocation                       bool                                  `php:"is_custom_location"`
	WsAlternateCode                        string                                `php:"ws_alternate_code"`
	IsLoyaltyEnrolled                      bool                                  `php:"is_loyalty_enrolled"`
	WmsStoreID                             int                                   `php:"wms_store_id"`
	WsStoreID                              string                                `php:"ws_store_id"`
	ProductCategoriesDiscounts             map[int]ProductCategoryDiscount       `php:"product_categories_discounts"`
	B2BProductCategoriesDiscounts          map[string]B2BProductCategoryDiscount `php:"b2b_product_categories_discounts"`
	OrganizationProductCategoriesDiscounts map[int]map[int]OrgCategoryDiscount   `php:"organization_product_categories_discounts"`
	Organizations                          map[int]Organization                  `php:"organizations"`
	Address1                               string                                `php:"address_1"`
	Address2                               string                                `php:"address_2"`
	City                                   string                                `php:"city"`
	State                                  string                                `php:"state"`
	Pincode                                string                                `php:"pincode"`
	DrugLicenseNumber                      []StoreDrugLicense                    `php:"drug_license_number"`
	StoreManagerDetails                    []StoreManager                        `php:"store_manager_details"`
	IsPosApplicable                        bool                                  `php:"is_pos_applicable"`
	GstTreatment                           string                                `php:"gst_treatment"`
	CinNumber                              *string                               `php:"cin_number"`
	StoreSettings                          StoreSettings                         `php:"store_settings"`
	IsWhatsappBillApplicable               bool                                  `php:"is_whatsapp_bill_applicable"`
	IsActive                               bool                                  `php:"is_active"`
	BillingOrganizationID                  int                                   `php:"billing_organization_id"`
}
