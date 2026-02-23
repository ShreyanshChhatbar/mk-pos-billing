package service

import "time"

type CachedRBAC struct {
	Tenants                map[uint64]TenantRBAC    `json:"tenants"`
	CurrentFacilityID      *uint64                  `json:"current_facility_id"`
	CurrentFacilityCheckIn *FacilityCheckInSnapshot `json:"current_facility_check_in"`
}

type TenantRBAC struct {
	Roles      map[uint64]RoleRBAC            `json:"roles"`
	Facilities map[uint64]FacilityRoleBinding `json:"facilities"`
}

type RoleRBAC struct {
	IsFacilityLoginRequired bool            `json:"is_facility_login_required"`
	PermissionIndex         map[string]bool `json:"permission_index"`
}

type FacilityRoleBinding struct {
	RoleID uint64 `json:"role_id"`
}

type FacilityCheckInSnapshot struct {
	CheckOutTime *time.Time `json:"check_out_time"`
}
