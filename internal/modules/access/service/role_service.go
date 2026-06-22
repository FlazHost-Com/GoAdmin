package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	apperr "goadmin/internal/errors"
	"goadmin/internal/helpers"
	"goadmin/internal/modules/access/dto"
	"goadmin/internal/modules/access/model"
)

// RoleService mengimplementasi IRoleService.
type RoleService struct {
	db *gorm.DB
}

var _ IRoleService = (*RoleService)(nil)

// NewRoleService merakit service.
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) Index(ctx context.Context, q dto.ListQuery) (helpers.Paginated[model.Role], error) {
	query := s.db.WithContext(ctx).Model(&model.Role{})
	if q.Search != "" {
		query = helpers.CiLike(query, "name", q.Search)
	}
	query = query.Order("name ASC").Preload("Permissions")

	var roles []model.Role
	meta, err := helpers.Paginate(query, q.Page, q.PerPage, &roles)
	if err != nil {
		return helpers.Paginated[model.Role]{}, apperr.Internal(err.Error())
	}
	return helpers.Paginated[model.Role]{Data: roles, Meta: meta}, nil
}

func (s *RoleService) Show(ctx context.Context, id string) (*model.Role, error) {
	var role model.Role
	if err := s.db.WithContext(ctx).Preload("Permissions").First(&role, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("Role tidak ditemukan")
		}
		return nil, apperr.Internal(err.Error())
	}
	return &role, nil
}

func (s *RoleService) Store(ctx context.Context, in dto.CreateRoleInput) (*model.Role, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Role{}).Where("name = ?", in.Name).Count(&count).Error; err != nil {
		return nil, apperr.Internal(err.Error())
	}
	if count > 0 {
		return nil, apperr.Conflict("Nama role sudah terpakai")
	}

	role := model.Role{
		ID:        helpers.NewID(),
		Name:      in.Name,
		GuardName: "web",
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return syncRolePermissions(tx, &role, in.PermissionIDs)
	})
	if err != nil {
		return nil, apperr.Internal(err.Error())
	}
	return s.Show(ctx, role.ID)
}

func (s *RoleService) Update(ctx context.Context, id string, in dto.UpdateRoleInput) (*model.Role, error) {
	role, err := s.Show(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != role.Name {
		var count int64
		if err := s.db.WithContext(ctx).Model(&model.Role{}).
			Where("name = ? AND id <> ?", in.Name, id).Count(&count).Error; err != nil {
			return nil, apperr.Internal(err.Error())
		}
		if count > 0 {
			return nil, apperr.Conflict("Nama role sudah terpakai")
		}
	}
	role.Name = in.Name

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(role).Error; err != nil {
			return err
		}
		return syncRolePermissions(tx, role, in.PermissionIDs)
	})
	if err != nil {
		return nil, apperr.Internal(err.Error())
	}
	return s.Show(ctx, id)
}

func (s *RoleService) Destroy(ctx context.Context, id string) error {
	role, err := s.Show(ctx, id)
	if err != nil {
		return err
	}
	if role.Name == model.RoleAdministrator {
		return apperr.Forbidden("Role Administrator tak boleh dihapus")
	}
	if err := s.db.WithContext(ctx).Select("Permissions", "Users").Delete(role).Error; err != nil {
		return apperr.Internal(err.Error())
	}
	return nil
}

func syncRolePermissions(tx *gorm.DB, role *model.Role, permIDs []string) error {
	if permIDs == nil {
		return nil
	}
	var perms []model.Permission
	if len(permIDs) > 0 {
		if err := tx.Where("id IN ?", permIDs).Find(&perms).Error; err != nil {
			return err
		}
	}
	return tx.Model(role).Association("Permissions").Replace(perms)
}
