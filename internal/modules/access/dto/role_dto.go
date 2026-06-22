package dto

// CreateRoleInput = payload buat role.
type CreateRoleInput struct {
	Name          string   `json:"name" form:"name" binding:"required,max=50"`
	PermissionIDs []string `json:"permission_ids" form:"permission_ids" binding:"omitempty,dive,max=36"`
}

// UpdateRoleInput = payload ubah role.
type UpdateRoleInput struct {
	Name          string   `json:"name" form:"name" binding:"required,max=50"`
	PermissionIDs []string `json:"permission_ids" form:"permission_ids" binding:"omitempty,dive,max=36"`
}

// CreatePermissionInput = payload buat permission.
type CreatePermissionInput struct {
	Name string `json:"name" form:"name" binding:"required,max=100"`
}

// UpdatePermissionInput = payload ubah permission.
type UpdatePermissionInput struct {
	Name string `json:"name" form:"name" binding:"required,max=100"`
}
