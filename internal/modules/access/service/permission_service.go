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

// PermissionService mengimplementasi IPermissionService.
type PermissionService struct {
	db *gorm.DB
}

var _ IPermissionService = (*PermissionService)(nil)

// NewPermissionService merakit service.
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

func (s *PermissionService) Index(ctx context.Context, q dto.ListQuery) (helpers.Paginated[model.Permission], error) {
	query := s.db.WithContext(ctx).Model(&model.Permission{})
	if q.Search != "" {
		query = helpers.CiLike(query, "name", q.Search)
	}
	query = query.Order("name ASC")

	var perms []model.Permission
	meta, err := helpers.Paginate(query, q.Page, q.PerPage, &perms)
	if err != nil {
		return helpers.Paginated[model.Permission]{}, apperr.Internal(err.Error())
	}
	return helpers.Paginated[model.Permission]{Data: perms, Meta: meta}, nil
}

func (s *PermissionService) Show(ctx context.Context, id string) (*model.Permission, error) {
	var perm model.Permission
	if err := s.db.WithContext(ctx).First(&perm, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("Permission tidak ditemukan")
		}
		return nil, apperr.Internal(err.Error())
	}
	return &perm, nil
}

func (s *PermissionService) Store(ctx context.Context, in dto.CreatePermissionInput) (*model.Permission, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.Permission{}).Where("name = ?", in.Name).Count(&count).Error; err != nil {
		return nil, apperr.Internal(err.Error())
	}
	if count > 0 {
		return nil, apperr.Conflict("Nama permission sudah terpakai")
	}
	perm := model.Permission{
		ID:        helpers.NewID(),
		Name:      in.Name,
		GuardName: "web",
	}
	if err := s.db.WithContext(ctx).Create(&perm).Error; err != nil {
		return nil, apperr.Internal(err.Error())
	}
	return &perm, nil
}

func (s *PermissionService) Update(ctx context.Context, id string, in dto.UpdatePermissionInput) (*model.Permission, error) {
	perm, err := s.Show(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != perm.Name {
		var count int64
		if err := s.db.WithContext(ctx).Model(&model.Permission{}).
			Where("name = ? AND id <> ?", in.Name, id).Count(&count).Error; err != nil {
			return nil, apperr.Internal(err.Error())
		}
		if count > 0 {
			return nil, apperr.Conflict("Nama permission sudah terpakai")
		}
	}
	perm.Name = in.Name
	if err := s.db.WithContext(ctx).Save(perm).Error; err != nil {
		return nil, apperr.Internal(err.Error())
	}
	return perm, nil
}

func (s *PermissionService) Destroy(ctx context.Context, id string) error {
	perm, err := s.Show(ctx, id)
	if err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Select("Roles").Delete(perm).Error; err != nil {
		return apperr.Internal(err.Error())
	}
	return nil
}
