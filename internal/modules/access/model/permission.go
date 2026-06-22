// Package model berisi entity GORM modul access (User, Role, Permission).
// Tipe kolom dijaga PORTABEL (string/text/int/bool/timestamp abstrak) — tanpa
// tipe vendor (longtext/datetime) atau collation hardcoded. Checker menolak
// pelanggaran ini agar app tetap multi-DB (MySQL/Postgres/SQLite).
package model

import "time"

// Permission = satu izin granular (mis. "user.create").
type Permission struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	GuardName string    `gorm:"type:varchar(50);index" json:"guard_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Roles []Role `gorm:"many2many:roles_permissions;" json:"-"`
}

// TableName memetakan ke tabel 'permissions'.
func (Permission) TableName() string { return "permissions" }
