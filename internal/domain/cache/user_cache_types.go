package cache

type Designation struct {
	ID   int    `php:"id"`
	Name string `php:"name"`
}

type UserRole struct {
	RoleName  string `php:"role_name"`
	StoreName string `php:"store_name"`
	RoleID    int    `php:"role"`
	StoreID   int    `php:"store"`
}

type User struct {
	ID                int                            `php:"id"`
	Designation       *Designation                   `php:"designation"`
	Name              string                         `php:"name"`
	MobileNumber      string                         `php:"mobile_number"`
	IsActive          bool                           `php:"is_active"`
	IsSuperAdmin      bool                           `php:"is_super_admin"`
	PassKey           string                         `php:"pass_key"`
	LastLogin         string                         `php:"last_login"`
	DeviceMasterID    int                            `php:"device_master_id"`
	Permissions       []UserRole                     `php:"permissions"`
	ModulePermissions map[string]map[string][]string `php:"module_permissions"`
}
