package migration

import (
	"gorm.io/gorm"

	"goadmin/internal/auth"
	"goadmin/internal/helpers"
	"goadmin/internal/modules/access/model"
)

// CorePermissions = izin granular bawaan untuk modul inti.
var CorePermissions = []string{
	"user.view", "user.create", "user.update", "user.delete",
	"role.view", "role.create", "role.update", "role.delete",
	"permission.view", "permission.create", "permission.update", "permission.delete",
}

// Seed membuat data awal: permission inti, role Administrator (semua izin),
// dan user admin default. Idempoten (aman dijalankan berulang).
func Seed(db *gorm.DB, adminEmail, adminPassword string, bcryptRounds int) error {
	// 1. Permissions.
	perms := make([]model.Permission, 0, len(CorePermissions))
	for _, name := range CorePermissions {
		p := model.Permission{ID: helpers.NewID(), Name: name, GuardName: "web"}
		if err := db.Where("name = ?", name).FirstOrCreate(&p, model.Permission{Name: name}).Error; err != nil {
			return err
		}
		perms = append(perms, p)
	}

	// 2. Role Administrator + semua permission.
	admin := model.Role{ID: helpers.NewID(), Name: model.RoleAdministrator, GuardName: "web"}
	if err := db.Where("name = ?", model.RoleAdministrator).FirstOrCreate(&admin, model.Role{Name: model.RoleAdministrator}).Error; err != nil {
		return err
	}
	if err := db.Model(&admin).Association("Permissions").Replace(perms); err != nil {
		return err
	}

	// 3. User admin default.
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
