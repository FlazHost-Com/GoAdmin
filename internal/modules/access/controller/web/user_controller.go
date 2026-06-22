// Package web berisi controller HTML modul access (jalur sesi). Render lewat
// view.RenderView (bukan c.HTML path mentah) — di-enforce checker.
package web

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	accessmw "goadmin/internal/modules/access/middleware"
	"goadmin/internal/modules/access/dto"
	"goadmin/internal/modules/access/model"
	"goadmin/internal/modules/access/service"
	"goadmin/internal/storage"
	"goadmin/internal/view"
)

// UserController menyajikan CRUD user (web). Memakai RoleService untuk pilihan
// role pada form (checkbox) + storage untuk upload foto.
type UserController struct {
	users   service.IUserService
	roles   service.IRoleService
	storage storage.Storage
}

// NewUserController merakit controller (service + storage di-inject).
func NewUserController(users service.IUserService, roles service.IRoleService, store storage.Storage) *UserController {
	return &UserController{users: users, roles: roles, storage: store}
}

// pictureURL memproses file "picture" (opsional): validasi + simpan → URL.
func (ctl *UserController) pictureURL(c *gin.Context) (string, error) {
	fh, err := c.FormFile("picture")
	if err != nil || fh == nil {
		return "", nil // tak ada file
	}
	f, oerr := fh.Open()
	if oerr != nil {
		return "", oerr
	}
	defer f.Close()
	return storage.ValidateAndSave(c.Request.Context(), ctl.storage, f)
}

// Index → GET /admin/v1/users (daftar user + search + paginasi).
func (ctl *UserController) Index(c *gin.Context) {
	var q dto.ListQuery
	_ = c.ShouldBindQuery(&q)
	res, err := ctl.users.Index(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "access/users/index", gin.H{
		"title": "Manajemen User", "active": "users",
		"users": res.Data, "meta": res.Meta, "search": q.Search,
	})
}

// Create → GET /admin/v1/users/create.
func (ctl *UserController) Create(c *gin.Context) {
	roles, err := ctl.allRoles(c)
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "users/form", gin.H{
		"title": "Tambah Pengguna", "active": "users",
		"action": "/admin/v1/users", "user": nil,
		"roles": roles, "selected": map[string]bool{},
	})
}

// Store → POST /admin/v1/users.
func (ctl *UserController) Store(c *gin.Context) {
	var in dto.CreateUserInput
	_ = c.ShouldBind(&in)
	if url, perr := ctl.pictureURL(c); perr != nil {
		setFlashError(sessions.Default(c), errMessage(perr))
		c.Redirect(http.StatusFound, "/admin/v1/users/create")
		return
	} else if url != "" {
		in.Picture = url
	}
	if _, err := ctl.users.Store(c.Request.Context(), in, actorID(c)); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/users/create")
		return
	}
	setFlashSuccess(sessions.Default(c), "Pengguna berhasil dibuat.")
	c.Redirect(http.StatusFound, "/admin/v1/users")
}

// Edit → GET /admin/v1/users/:id/edit.
func (ctl *UserController) Edit(c *gin.Context) {
	user, err := ctl.users.Show(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	roles, err := ctl.allRoles(c)
	if err != nil {
		c.Error(err)
		return
	}
	selected := make(map[string]bool, len(user.Roles))
	for _, r := range user.Roles {
		selected[r.ID] = true
	}
	view.RenderView(c, "users/form", gin.H{
		"title": "Ubah Pengguna", "active": "users",
		"action": "/admin/v1/users/" + user.ID, "user": user,
		"roles": roles, "selected": selected,
	})
}

// Update → POST /admin/v1/users/:id.
func (ctl *UserController) Update(c *gin.Context) {
	id := c.Param("id")
	var in dto.UpdateUserInput
	_ = c.ShouldBind(&in)
	if url, perr := ctl.pictureURL(c); perr != nil {
		setFlashError(sessions.Default(c), errMessage(perr))
		c.Redirect(http.StatusFound, "/admin/v1/users/"+id+"/edit")
		return
	} else if url != "" {
		in.Picture = url
	}
	if _, err := ctl.users.Update(c.Request.Context(), id, in, actorID(c)); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/users/"+id+"/edit")
		return
	}
	setFlashSuccess(sessions.Default(c), "Pengguna berhasil diperbarui.")
	c.Redirect(http.StatusFound, "/admin/v1/users")
}

// Destroy → POST /admin/v1/users/:id/delete.
func (ctl *UserController) Destroy(c *gin.Context) {
	if err := ctl.users.Destroy(c.Request.Context(), c.Param("id")); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
	} else {
		setFlashSuccess(sessions.Default(c), "Pengguna berhasil dihapus.")
	}
	c.Redirect(http.StatusFound, "/admin/v1/users")
}

// allRoles mengambil seluruh role (untuk pilihan di form).
func (ctl *UserController) allRoles(c *gin.Context) ([]model.Role, error) {
	res, err := ctl.roles.Index(c.Request.Context(), dto.ListQuery{Page: 1, PerPage: 100})
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

// actorID mengambil id user terautentikasi dari context (untuk audit created_by).
func actorID(c *gin.Context) string {
	if u := accessmw.UserFrom(c); u != nil {
		return u.ID
	}
	return ""
}
