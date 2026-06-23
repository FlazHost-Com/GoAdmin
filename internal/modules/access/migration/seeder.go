package migration

import (
	"errors"

	"gorm.io/gorm"

	"goadmin/internal/auth"
	"goadmin/internal/helpers"
	"goadmin/internal/modules/access/model"
)

// Seed membuat data awal: role Administrator + user admin default. Idempoten.
//
// CATATAN RBAC route-driven (a la NodeAdmin): permission TIDAK lagi di-seed dari
// daftar tetap — diturunkan dari named-route registry lewat
// bootstrap.SyncPermissions (dipanggil setelah app.Build) + lazy saat buka
// halaman Permission. Administrator BYPASS RBAC (IsAdministrator), jadi tak perlu
// di-assign permission apa pun di sini.
func Seed(db *gorm.DB, adminEmail, adminPassword string, bcryptRounds int) error {
	// 1. Role Administrator (guard web). Tanpa assignment permission — bypass.
	var admin model.Role
	err := db.Where("name = ? AND guard_name = ?", model.RoleAdministrator, "web").First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		admin = model.Role{ID: helpers.NewID(), Name: model.RoleAdministrator, GuardName: "web", Status: model.StatusActive}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 2. User admin default.
	var count int64
	if err := db.Model(&model.User{}).Where("email = ?", adminEmail).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		hash, err := auth.HashPassword(adminPassword, bcryptRounds)
		if err != nil {
			return err
		}
		user := model.User{
			ID:       helpers.NewID(),
			Code:     helpers.NewCode("U"),
			Name:     "Administrator",
			Email:    adminEmail,
			Password: hash,
			Status:   model.StatusActive,
			Timezone: "UTC",
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
		if err := db.Model(&user).Association("Roles").Append(&admin); err != nil {
			return err
		}
	}
	return nil
}
