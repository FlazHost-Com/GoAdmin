package api

import (
	"github.com/gin-gonic/gin"

	apperr "goadmin/internal/errors"
	"goadmin/internal/helpers"
	"goadmin/internal/modules/access/dto"
	"goadmin/internal/modules/access/service"
)

// RoleController = REST CRUD role.
type RoleController struct {
	roles service.IRoleService
}

// NewRoleController merakit controller.
func NewRoleController(roles service.IRoleService) *RoleController {
	return &RoleController{roles: roles}
}

// Index → GET /api/v1/roles.
func (ctl *RoleController) Index(c *gin.Context) {
	var q dto.ListQuery
	_ = c.ShouldBindQuery(&q)
	res, err := ctl.roles.Index(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "OK", res)
}

// Show → GET /api/v1/roles/:id.
func (ctl *RoleController) Show(c *gin.Context) {
	role, err := ctl.roles.Show(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "OK", role)
}

// Store → POST /api/v1/roles.
func (ctl *RoleController) Store(c *gin.Context) {
	var in dto.CreateRoleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(apperr.Validation("Input tidak valid", nil))
		return
	}
	role, err := ctl.roles.Store(c.Request.Context(), in)
	if err != nil {
		c.Error(err)
		return
	}
	helpers.Created(c, "Role dibuat", role)
}

// Update → PUT /api/v1/roles/:id.
func (ctl *RoleController) Update(c *gin.Context) {
	var in dto.UpdateRoleInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(apperr.Validation("Input tidak valid", nil))
		return
	}
	role, err := ctl.roles.Update(c.Request.Context(), c.Param("id"), in)
	if err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "Role diperbarui", role)
}

// Destroy → DELETE /api/v1/roles/:id.
func (ctl *RoleController) Destroy(c *gin.Context) {
	if err := ctl.roles.Destroy(c.Request.Context(), c.Param("id")); err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "Role dihapus", nil)
}
