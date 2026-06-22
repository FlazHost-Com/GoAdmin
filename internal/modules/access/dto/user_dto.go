// Package dto berisi struct input tervalidasi (padanan validator Joi
// stripUnknown di NodeAdmin). Hanya field di DTO yang diterima → anti
// mass-assignment (whitelist). Validasi lewat tag go-playground/validator.
package dto

// CreateUserInput = payload buat user baru.
type CreateUserInput struct {
	Name     string   `json:"name" form:"name" binding:"required,max=50"`
	Email    string   `json:"email" form:"email" binding:"required,email,max=255"`
	Phone    string   `json:"phone" form:"phone" binding:"omitempty,max=15"`
	Password string   `json:"password" form:"password" binding:"required,min=8,max=72"`
	Status   string   `json:"status" form:"status" binding:"omitempty,oneof=Active Inactive"`
	Timezone string   `json:"timezone" form:"timezone" binding:"omitempty,max=64"`
	Picture  string   `json:"picture" form:"picture" binding:"omitempty,max=255"`
	RoleIDs  []string `json:"role_ids" form:"role_ids" binding:"omitempty,dive,max=36"`
}

// UpdateUserInput = payload ubah user. Password opsional (kosong = tak diubah).
type UpdateUserInput struct {
	Name     string   `json:"name" form:"name" binding:"required,max=50"`
	Email    string   `json:"email" form:"email" binding:"required,email,max=255"`
	Phone    string   `json:"phone" form:"phone" binding:"omitempty,max=15"`
	Password string   `json:"password" form:"password" binding:"omitempty,min=8,max=72"`
	Status   string   `json:"status" form:"status" binding:"omitempty,oneof=Active Inactive"`
	Timezone string   `json:"timezone" form:"timezone" binding:"omitempty,max=64"`
	Picture  string   `json:"picture" form:"picture" binding:"omitempty,max=255"`
	RoleIDs  []string `json:"role_ids" form:"role_ids" binding:"omitempty,dive,max=36"`
}

// ListQuery = parameter list (paginasi + search).
type ListQuery struct {
	Page    int    `form:"page"`
	PerPage int    `form:"per_page"`
	Search  string `form:"search"`
}
