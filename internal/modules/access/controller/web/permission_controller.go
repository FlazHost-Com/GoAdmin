package web

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"goadmin/internal/modules/access/dto"
	"goadmin/internal/modules/access/service"
	"goadmin/internal/view"
)

// PermissionController menyajikan CRUD permission (web) — hanya nama.
type PermissionController struct {
	perms service.IPermissionService
}

// NewPermissionController merakit controller (service di-inject).
func NewPermissionController(perms service.IPermissionService) *PermissionController {
	return &PermissionController{perms: perms}
}

// Index → GET /admin/v1/permissions.
func (ctl *PermissionController) Index(c *gin.Context) {
	var q dto.ListQuery
	_ = c.ShouldBindQuery(&q)
	res, err := ctl.perms.Index(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "permissions/index", gin.H{
		"title": "Manajemen Permission", "active": "permissions",
		"permissions": res.Data, "meta": res.Meta, "search": q.Search,
	})
}

// Create → GET /admin/v1/permissions/create.
func (ctl *PermissionController) Create(c *gin.Context) {
	view.RenderView(c, "permissions/form", gin.H{
		"title": "Tambah Permission", "active": "permissions",
		"action": "/admin/v1/permissions", "permission": nil,
	})
}

// Store → POST /admin/v1/permissions.
func (ctl *PermissionController) Store(c *gin.Context) {
	var in dto.CreatePermissionInput
	_ = c.ShouldBind(&in)
	if _, err := ctl.perms.Store(c.Request.Context(), in); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/permissions/create")
		return
	}
	setFlashSuccess(sessions.Default(c), "Permission berhasil dibuat.")
	c.Redirect(http.StatusFound, "/admin/v1/permissions")
}

// Edit → GET /admin/v1/permissions/:id/edit.
func (ctl *PermissionController) Edit(c *gin.Context) {
	perm, err := ctl.perms.Show(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "permissions/form", gin.H{
		"title": "Ubah Permission", "active": "permissions",
		"action": "/admin/v1/permissions/" + perm.ID, "permission": perm,
	})
}

// Update → POST /admin/v1/permissions/:id.
func (ctl *PermissionController) Update(c *gin.Context) {
	id := c.Param("id")
	var in dto.UpdatePermissionInput
	_ = c.ShouldBind(&in)
	if _, err := ctl.perms.Update(c.Request.Context(), id, in); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/permissions/"+id+"/edit")
		return
	}
	setFlashSuccess(sessions.Default(c), "Permission berhasil diperbarui.")
	c.Redirect(http.StatusFound, "/admin/v1/permissions")
}

// Destroy → POST /admin/v1/permissions/:id/delete.
func (ctl *PermissionController) Destroy(c *gin.Context) {
	if err := ctl.perms.Destroy(c.Request.Context(), c.Param("id")); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
	} else {
		setFlashSuccess(sessions.Default(c), "Permission berhasil dihapus.")
	}
	c.Redirect(http.StatusFound, "/admin/v1/permissions")
}
