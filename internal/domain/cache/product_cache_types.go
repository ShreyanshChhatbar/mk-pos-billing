package cache

type ComboProduct struct {
	ProductID int     `php:"product_id"`
	Quantity  int     `php:"quantity"`
	SalesRate float64 `php:"sales_rate"`
}

type Product struct {
	ID                            int            `php:"id"`
	ProductName                   string         `php:"product_name"`
	ProductType                   string         `php:"product_type"`
	WsCode                        int            `php:"ws_code"`
	ProductCode                   int            `php:"product_code"`
	Mrp                           float64        `php:"mrp"`
	OldMrp                        float64        `php:"old_mrp"`
	ScheduledTypeCode             string         `php:"scheduled_type_code"`
	MisReportingCategory          string         `php:"mis_reporting_category"`
	IsActive                      bool           `php:"is_active"`
	DosageForm                    string         `php:"dosage_form"`
	PackageType                   string         `php:"package_type"`
	Uom                           string         `php:"uom"`
	PackageSizeOld                int            `php:"package_size_old"`
	PackageSize                   string         `php:"package_size"`
	SalesUnit                     int            `php:"sales_unit"`
	GstType                       string         `php:"gst_type"`
	OldGstType                    string         `php:"old_gst_type"`
	WmsProductID                  int            `php:"wms_product_id"`
	IsHiddenFromAlternateProducts bool           `php:"is_hidden_from_alternate_products"`
	IsAssured                     bool           `php:"is_assured"`
	IsDiscontinued                bool           `php:"is_discontinued"`
	IsBanned                      bool           `php:"is_banned"`
	IsAlternateAvailable          bool           `php:"is_alternate_available"`
	CombinationsString            string         `php:"combinations_string"`
	CombinationsStringSlug        string         `php:"combinations_string_slug"`
	Pack                          string         `php:"pack"`
	IsRefrigerated                bool           `php:"is_refrigerated"`
	IsChronic                     bool           `php:"is_chronic"`
	IsRxRequired                  bool           `php:"is_rx_required"`
	B2CProductCategoryID          int            `php:"b_2_c_product_category_id"`
	OrganizationCategoryID        int            `php:"organization_category_id"`
	IsGeneric                     bool           `php:"is_generic"`
	HsnCode                       string         `php:"hsn_code"`
	B2BProductType                string         `php:"b_2_b_product_type"`
	B2CProductType                string         `php:"b_2_c_product_type"`
	ManufacturerName              string         `php:"manufacturer_name"`
	VendorReturnType              string         `php:"vendor_return_type"`
	IsComboProduct                bool           `php:"is_combo_product"`
	IsMspProduct                  bool           `php:"is_msp_product"`
	IsSpecialMolecule             bool           `php:"is_special_molecule"`
	IsFreeCombo                   bool           `php:"is_free_combo"`
	ComboMov                      float64        `php:"combo_mov"`
	ComboStoreIds                 []int          `php:"combo_store_ids"`
	ComboProducts                 []ComboProduct `php:"combo_products"`
	ApplicableType                string         `php:"applicable_type"`
	ClusterMasterID               int            `php:"cluster_master_id"`
	ComboSalesPrice               float64        `php:"combo_sales_price"`
	MisReportingCategoryID        int            `php:"mis_reporting_category_id"`
	IsOngcRestricted              bool           `php:"is_ongc_restricted"`
	TransferIn                    float64        `php:"transfer_in"`
	TransferOut                   float64        `php:"transfer_out"`
	FranchiseIn                   float64        `php:"franchise_in"`
	FranchiseOut                  float64        `php:"franchise_out"`
}
