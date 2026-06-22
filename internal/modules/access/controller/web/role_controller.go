package web

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"goadmin/internal/modules/access/dto"
	"goadmin/internal/modules/access/model"
	"goadmin/internal/modules/access/service"
	"goadmin/internal/view"
)

// RoleController menyajikan CRUD role (web). Memakai PermissionService untuk
// menyediakan daftar permission pada form (checkbox).
type RoleController struct {
	roles service.IRoleService
	perms service.IPermissionService
}

// NewRoleController merakit controller (service di-inject).
func NewRoleController(roles service.IRoleService, perms service.IPermissionService) *RoleController {
	return &RoleController{roles: roles, perms: perms}
}

// Index → GET /admin/v1/roles.
func (ctl *RoleController) Index(c *gin.Context) {
	var q dto.ListQuery
	_ = c.ShouldBindQuery(&q)
	res, err := ctl.roles.Index(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "roles/index", gin.H{
		"title": "Manajemen Role", "active": "roles",
		"roles": res.Data, "meta": res.Meta, "search": q.Search,
	})
}

// Create → GET /admin/v1/roles/create.
func (ctl *RoleController) Create(c *gin.Context) {
	perms, err := ctl.allPermissions(c)
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "roles/form", gin.H{
		"title": "Tambah Role", "active": "roles",
		"action": "/admin/v1/roles", "role": nil,
		"permissions": perms, "selected": map[string]bool{},
	})
}

// Store → POST /admin/v1/roles.
func (ctl *RoleController) Store(c *gin.Context) {
	var in dto.CreateRoleInput
	_ = c.ShouldBind(&in)
	if _, err := ctl.roles.Store(c.Request.Context(), in); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/roles/create")
		return
	}
	setFlashSuccess(sessions.Default(c), "Role berhasil dibuat.")
	c.Redirect(http.StatusFound, "/admin/v1/roles")
}

// Edit → GET /admin/v1/roles/:id/edit.
func (ctl *RoleController) Edit(c *gin.Context) {
	role, err := ctl.roles.Show(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	perms, err := ctl.allPermissions(c)
	if err != nil {
		c.Error(err)
		return
	}
	selected := make(map[string]bool, len(role.Permissions))
	for _, p := range role.Permissions {
		selected[p.ID] = true
	}
	view.RenderView(c, "roles/form", gin.H{
		"title": "Ubah Role", "active": "roles",
		"action": "/admin/v1/roles/" + role.ID, "role": role,
		"permissions": perms, "selected": selected,
	})
}

// Update → POST /admin/v1/roles/:id.
func (ctl *RoleController) Update(c *gin.Context) {
	id := c.Param("id")
	var in dto.UpdateRoleInput
	_ = c.ShouldBind(&in)
	if _, err := ctl.roles.Update(c.Request.Context(), id, in); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/roles/"+id+"/edit")
		return
	}
	setFlashSuccess(sessions.Default(c), "Role berhasil diperbarui.")
	c.Redirect(http.StatusFound, "/admin/v1/roles")
}

// Destroy → POST /admin/v1/roles/:id/delete.
func (ctl *RoleController) Destroy(c *gin.Context) {
	if err := ctl.roles.Destroy(c.Request.Context(), c.Param("id")); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
	} else {
		setFlashSuccess(sessions.Default(c), "Role berhasil dihapus.")
	}
	c.Redirect(http.StatusFound, "/admin/v1/roles")
}

// allPermissions mengambil seluruh permission (untuk pilihan di form).
func (ctl *RoleController) allPermissions(c *gin.Context) ([]model.Permission, error) {
	res, err := ctl.perms.Index(c.Request.Context(), dto.ListQuery{Page: 1, PerPage: 100})
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}
