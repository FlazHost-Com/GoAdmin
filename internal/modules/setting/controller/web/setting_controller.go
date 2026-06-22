// Package web berisi controller HTML modul setting. Render lewat
// view.RenderView (bukan c.HTML path mentah) — di-enforce checker.
package web

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	apperr "goadmin/internal/errors"
	"goadmin/internal/middleware"
	accessmw "goadmin/internal/modules/access/middleware"
	"goadmin/internal/modules/setting/dto"
	"goadmin/internal/modules/setting/service"
	"goadmin/internal/storage"
	"goadmin/internal/view"
)

// SettingController menyajikan halaman pengaturan global + theme switcher +
// upload logo.
type SettingController struct {
	settings service.ISettingService
	storage  storage.Storage
}

// NewSettingController merakit controller (service + storage di-inject).
func NewSettingController(settings service.ISettingService, store storage.Storage) *SettingController {
	return &SettingController{settings: settings, storage: store}
}

// Index → GET /admin/v1/setting (form pengaturan + pilihan tema).
func (ctl *SettingController) Index(c *gin.Context) {
	setting, err := ctl.settings.Get(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	view.RenderView(c, "setting/index", gin.H{
		"title":   "Pengaturan",
		"active":  "setting",
		"setting": setting,
		"themes":  ctl.settings.Themes(),
	})
}

// Update → POST /admin/v1/setting (PRG: simpan lalu redirect balik). Bila ada
// file logo, divalidasi (magic-byte) + disimpan, lalu URL-nya dipakai.
func (ctl *SettingController) Update(c *gin.Context) {
	var in dto.UpdateSettingInput
	_ = c.ShouldBind(&in)

	if fh, err := c.FormFile("logo"); err == nil && fh != nil {
		url, uerr := ctl.saveLogo(c, fh)
		if uerr != nil {
			setFlashError(sessions.Default(c), errMessage(uerr))
			c.Redirect(http.StatusFound, "/admin/v1/setting")
			return
		}
		in.Logo = url
	}

	if _, err := ctl.settings.Update(c.Request.Context(), in, actorID(c)); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/setting")
		return
	}
	setFlashSuccess(sessions.Default(c), "Pengaturan disimpan.")
	c.Redirect(http.StatusFound, "/admin/v1/setting")
}

// saveLogo membuka file upload lalu validasi (magic-byte) + simpan → URL publik.
func (ctl *SettingController) saveLogo(c *gin.Context, fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	return storage.ValidateAndSave(c.Request.Context(), ctl.storage, f)
}

// --- helper flash & actor (paket web access dipakai lewat re-deklarasi lokal) ---

func setFlashSuccess(sess sessions.Session, msg string) {
	sess.Set(middleware.FlashSuccessKey, msg)
	_ = sess.Save()
}

func setFlashError(sess sessions.Session, msg string) {
	sess.Set(middleware.FlashErrorKey, msg)
	_ = sess.Save()
}

func errMessage(err error) string {
	if ae, ok := apperr.As(err); ok {
		return ae.Message
	}
	return "Terjadi kesalahan."
}

func actorID(c *gin.Context) string {
	if u := accessmw.UserFrom(c); u != nil {
		return u.ID
	}
	return ""
}
