// Package web berisi controller HTML modul profile. Render lewat
// view.RenderView (bukan c.HTML path mentah) — di-enforce checker.
package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	accessmw "goadmin/internal/modules/access/middleware"
	"goadmin/internal/modules/profile/dto"
	"goadmin/internal/modules/profile/service"
	"goadmin/internal/storage"
	"goadmin/internal/view"
)

// ProfileController menyajikan halaman profil milik-sendiri.
type ProfileController struct {
	profiles service.IProfileService
	storage  storage.Storage
}

// NewProfileController merakit controller (service + storage di-inject).
func NewProfileController(profiles service.IProfileService, store storage.Storage) *ProfileController {
	return &ProfileController{profiles: profiles, storage: store}
}

// Index → GET /admin/v1/profile (form profil).
func (ctl *ProfileController) Index(c *gin.Context) {
	user := accessmw.UserFrom(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}
	profile, err := ctl.profiles.Get(c.Request.Context(), user.ID)
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "profile/index", gin.H{
		"title":   "Profil Saya",
		"active":  "profile",
		"profile": profile,
	})
}

// Update → POST /admin/v1/profile (PRG: simpan lalu redirect balik).
func (ctl *ProfileController) Update(c *gin.Context) {
	user := accessmw.UserFrom(c)
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}
	var in dto.UpdateProfileInput
	_ = c.ShouldBind(&in)

	// Avatar opsional: validasi (magic-byte) + simpan → URL.
	if fh, ferr := c.FormFile("picture"); ferr == nil && fh != nil {
		f, oerr := fh.Open()
		if oerr != nil {
			c.Error(oerr)
			return
		}
		defer f.Close()
		url, uerr := storage.ValidateAndSave(c.Request.Context(), ctl.storage, f)
		if uerr != nil {
			c.Error(uerr)
			return
		}
		in.Picture = url
	}

	if _, err := ctl.profiles.Update(c.Request.Context(), user.ID, in); err != nil {
		c.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/v1/profile")
}
