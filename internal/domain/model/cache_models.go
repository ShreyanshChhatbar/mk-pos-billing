package model

// ==========================================
// TILL CACHE MODELS
// ==========================================
type TillTransactionCache struct {
	ID     int    `php:"id"`
	Status string `php:"status"`
	Date   string `php:"date"`
}

type TillCache struct {
	ID              int                   `php:"id"`
	TillNumber      string                `php:"till_number"`
	IsActive        bool                  `php:"is_active"`
	StoreID         int                   `php:"store_id"`
	OpenedDatetime  string                `php:"opened_datetime"`
	TillTransaction *TillTransactionCache `php:"till_transaction"`
}

// ==========================================
// USER AUTH CACHE MODELS
// ==========================================
type DesignationCache struct {
	ID   int    `php:"id"`
	Name string `php:"name"`
}

type UserRoleCache struct {
	RoleName  string `php:"role_name"`
	StoreName string `php:"store_name"`
	RoleID    int    `php:"role"`
	StoreID   int    `php:"store"`
}

type UserCache struct {
	ID                int                            `php:"id"`
	Designation       *DesignationCache              `php:"designation"`
	Name              string                         `php:"name"`
	MobileNumber      string                         `php:"mobile_number"`
	IsActive          bool                           `php:"is_active"`
	IsSuperAdmin      bool                           `php:"is_super_admin"`
	PassKey           string                         `php:"pass_key"`
	LastLogin         string                         `php:"last_login"`
	DeviceMasterID    int                            `php:"device_master_id"`
	Permissions       []UserRoleCache                `php:"permissions"`
	ModulePermissions map[string]map[string][]string `php:"module_permissions"`
}

// Used for the specific POS auth token payload
type UserAuthCache struct {
	UserID      int                            `php:"user_id"`
	AuthToken   string                         `php:"auth_token"`
	DeviceToken string                         `php:"device_token"`
	StoreID     string                         `php:"store_id"`
	Permissions map[string]map[string][]string `php:"permissions"`
}

// ==========================================
// STORE CACHE MODELS
// ==========================================
type BankMasterCache struct {
	ID            int    `php:"id"`
	BankName      string `php:"bank_name"`
	AccountNumber string `php:"account_number"`
	IfscCode      string `php:"ifsc_code"`
}

type StoreSettingsCache struct {
	FloatCashAmount                  float64           `php:"float_cash_amount"`
	TillMismatchTaskTrigger          bool              `php:"till_mismatch_task_trigger"`
	TillMismatchTaskTriggerMinAmount float64           `php:"till_mismatch_task_trigger_min_amount"`
	IsOtpRequiredSalesReturn         bool              `php:"is_otp_required_sales_return"`
	IsDeliveryChargeApplicable       bool              `php:"is_delivery_charge_applicable"`
	DeliveryChargeMinimumMov         float64           `php:"delivery_charge_minimum_mov"`
	TransferInMinimumAmount          float64           `php:"transfer_in_minimum_amount"`
	TransferOutMinimumAmount         float64           `php:"transfer_out_minimum_amount"`
	IsUrgentOrderEnabled             bool              `php:"is_urgent_order_enabled"`
	IsAdvanceOrderEnabled            bool              `php:"is_advance_order_enabled"`
	IsMinMaxOrderEnabled             bool              `php:"is_min_max_order_enabled"`
	IsMinMaxDropshipOrderEnabled     bool              `php:"is_min_max_dropship_order_enabled"`
	BankMasters                      []BankMasterCache `php:"bank_masters"`
	UrgentOrderTime                  *string           `php:"urgent_order_time"`
}

type ProductCategoryDiscountCache struct {
	ID                   int     `php:"id"`
	B2CPricingTemplateID int     `php:"b_2_c_pricing_template_id"`
	CategoryID           int     `php:"category_id"`
	PricingCategory      string  `php:"pricing_category"`
	Mode                 string  `php:"mode"`
	Value                float64 `php:"value"`
	Operator             string  `php:"operator"`
	// PromoCodes           []string `php:"promo_codes"` // Optional based on mapping
}

type B2BProductCategoryDiscountCache struct {
	ID                     int    `php:"id"`
	B2BPricingTemplateID   int    `php:"b_2_b_pricing_template_id"`
	B2BProductCategoryName string `php:"b_2_b_product_category_name"`
	B2BPricingCategory     string `php:"b_2_b_pricing_category"`
	Mode                   string `php:"mode"`
	Value                  string `php:"value"`
	Operator               string `php:"operator"`
}

type OrgCategoryDiscountCache struct {
	CategoryName    string `php:"category_name"`
	Mode            string `php:"mode"`
	Operator        string `php:"operator"`
	PricingCategory string `php:"pricing_category"`
	Value           string `php:"value"`
}

type OrganizationCache struct {
	Name                  string  `php:"name"`
	OrganizationCode      string  `php:"organization_code"`
	IsCreditSystemEnabled bool    `php:"is_credit_system_enabled"`
	CreditLimit           float64 `php:"credit_limit"`
}

type StoreDrugLicenseCache struct {
	ID             int    `php:"id"`
	LicenseNumber  string `php:"license_number"`
	LicenseTagName string `php:"license_tag_name"`
}

type StoreManagerCache struct {
	ID                 int    `php:"id"`
	ManagerName        string `php:"manager_name"`
	ManagerEmail       string `php:"manager_email"`
	ManagerContact     string `php:"manager_contact"`
	PrimaryContactName string `php:"primary_contact_name"`
	StorePhoneNumber   string `php:"store_phone_number"`
}

type StoreCache struct {
	ID                                     int                                        `php:"id"`
	StoreName                              string                                     `php:"store_name"`
	StoreType                              string                                     `php:"store_type"`
	StoreClusters                          []int                                      `php:"store_clusters"`
	GstNumber                              string                                     `php:"gst_number"`
	IsReReleased                           bool                                       `php:"is_re_released"`
	StoreGroup                             string                                     `php:"store_group"`
	PlaceOfSupplyCode                      string                                     `php:"place_of_supply_code"`
	IsCustomLocation                       bool                                       `php:"is_custom_location"`
	WsAlternateCode                        string                                     `php:"ws_alternate_code"`
	IsLoyaltyEnrolled                      bool                                       `php:"is_loyalty_enrolled"`
	WmsStoreID                             int                                        `php:"wms_store_id"`
	WsStoreID                              string                                     `php:"ws_store_id"`
	ProductCategoriesDiscounts             map[int]ProductCategoryDiscountCache       `php:"product_categories_discounts"`
	B2BProductCategoriesDiscounts          map[string]B2BProductCategoryDiscountCache `php:"b2b_product_categories_discounts"`
	OrganizationProductCategoriesDiscounts map[int]map[int]OrgCategoryDiscountCache   `php:"organization_product_categories_discounts"`
	Organizations                          map[int]OrganizationCache                  `php:"organizations"`
	Address1                               string                                     `php:"address_1"`
	Address2                               string                                     `php:"address_2"`
	City                                   string                                     `php:"city"`
	State                                  string                                     `php:"state"`
	Pincode                                string                                     `php:"pincode"`
	DrugLicenseNumber                      []StoreDrugLicenseCache                    `php:"drug_license_number"`
	StoreManagerDetails                    []StoreManagerCache                        `php:"store_manager_details"`
	IsPosApplicable                        bool                                       `php:"is_pos_applicable"`
	GstTreatment                           string                                     `php:"gst_treatment"`
	CinNumber                              *string                                    `php:"cin_number"`
	StoreSettings                          StoreSettingsCache                         `php:"store_settings"`
	IsWhatsappBillApplicable               bool                                       `php:"is_whatsapp_bill_applicable"`
	IsActive                               bool                                       `php:"is_active"`
	BillingOrganizationID                  int                                        `php:"billing_organization_id"`
}

// ==========================================
// PRODUCT CACHE MODELS
// ==========================================
type ComboProductCache struct {
	ProductID int     `php:"product_id"`
	Quantity  int     `php:"quantity"`
	SalesRate float64 `php:"sales_rate"`
}

type ProductCache struct {
	ID                            int                 `php:"id"`
	ProductName                   string              `php:"product_name"`
	ProductType                   string              `php:"product_type"`
	WsCode                        int                 `php:"ws_code"`
	ProductCode                   int                 `php:"product_code"`
	Mrp                           float64             `php:"mrp"`
	OldMrp                        float64             `php:"old_mrp"`
	ScheduledTypeCode             string              `php:"scheduled_type_code"`
	MisReportingCategory          string              `php:"mis_reporting_category"`
	IsActive                      bool                `php:"is_active"`
	DosageForm                    string              `php:"dosage_form"`
	PackageType                   string              `php:"package_type"`
	Uom                           string              `php:"uom"`
	PackageSizeOld                int                 `php:"package_size_old"`
	PackageSize                   string              `php:"package_size"`
	SalesUnit                     int                 `php:"sales_unit"`
	GstType                       string              `php:"gst_type"`
	OldGstType                    string              `php:"old_gst_type"`
	WmsProductID                  int                 `php:"wms_product_id"`
	IsHiddenFromAlternateProducts bool                `php:"is_hidden_from_alternate_products"`
	IsAssured                     bool                `php:"is_assured"`
	IsDiscontinued                bool                `php:"is_discontinued"`
	IsBanned                      bool                `php:"is_banned"`
	IsAlternateAvailable          bool                `php:"is_alternate_available"`
	CombinationsString            string              `php:"combinations_string"`
	CombinationsStringSlug        string              `php:"combinations_string_slug"`
	Pack                          string              `php:"pack"`
	IsRefrigerated                bool                `php:"is_refrigerated"`
	IsChronic                     bool                `php:"is_chronic"`
	IsRxRequired                  bool                `php:"is_rx_required"`
	B2CProductCategoryID          int                 `php:"b_2_c_product_category_id"`
	OrganizationCategoryID        int                 `php:"organization_category_id"`
	IsGeneric                     bool                `php:"is_generic"`
	HsnCode                       string              `php:"hsn_code"`
	B2BProductType                string              `php:"b_2_b_product_type"`
	B2CProductType                string              `php:"b_2_c_product_type"`
	ManufacturerName              string              `php:"manufacturer_name"`
	VendorReturnType              string              `php:"vendor_return_type"`
	IsComboProduct                bool                `php:"is_combo_product"`
	IsMspProduct                  bool                `php:"is_msp_product"`
	IsSpecialMolecule             bool                `php:"is_special_molecule"`
	IsFreeCombo                   bool                `php:"is_free_combo"`
	ComboMov                      float64             `php:"combo_mov"`
	ComboStoreIds                 []int               `php:"combo_store_ids"`
	ComboProducts                 []ComboProductCache `php:"combo_products"`
	ApplicableType                string              `php:"applicable_type"`
	ClusterMasterID               int                 `php:"cluster_master_id"`
	ComboSalesPrice               float64             `php:"combo_sales_price"`
	MisReportingCategoryID        int                 `php:"mis_reporting_category_id"`
	IsOngcRestricted              bool                `php:"is_ongc_restricted"`
	TransferIn                    float64             `php:"transfer_in"`
	TransferOut                   float64             `php:"transfer_out"`
	FranchiseIn                   float64             `php:"franchise_in"`
	FranchiseOut                  float64             `php:"franchise_out"`
}

// ==========================================
// DEVICE CACHE MODELS
// ==========================================
type DeviceCache struct {
	DeviceID      int     `php:"device_id"`
	StoreID       *int    `php:"store_id"`
	IsActive      bool    `php:"is_active"`
	IsExpired     bool    `php:"is_expired"`
	DeviceType    string  `php:"device_type"`
	UserAuthToken *string `php:"user_auth_token"`
}
