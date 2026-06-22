package model

import "time"

// Role mengelompokkan permission dan diberikan ke user (RBAC).
type Role struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(50);uniqueIndex" json:"name"`
	GuardName string    `gorm:"type:varchar(50);index" json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Permissions []Permission `gorm:"many2many:roles_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:users_roles;" json:"-"`
}

func (Role) TableName() string { return "roles" }

// HasPermission true bila role memuat permission bernama tertentu.
func (r *Role) HasPermission(name string) bool {
	for _, p := range r.Permissions {
		if p.Name == name {
			return true
		}
	}
	return false
}
