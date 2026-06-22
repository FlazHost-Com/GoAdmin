package web

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"goadmin/internal/middleware"
	"goadmin/internal/modules/home/fetemplate"
	settingdto "goadmin/internal/modules/setting/dto"
	settingsvc "goadmin/internal/modules/setting/service"
	"goadmin/internal/view"
)

const appearancePerPage = 12

// AppearanceController = halaman admin pemilih template landing (frontend
// switcher). Katalog dari fetemplate.Service; aktif disimpan di Setting.fe_template.
type AppearanceController struct {
	settings func() settingsvc.ISettingService
	fe       *fetemplate.Service
}

// NewAppearanceController merakit controller.
func NewAppearanceController(settings func() settingsvc.ISettingService, fe *fetemplate.Service) *AppearanceController {
	return &AppearanceController{settings: settings, fe: fe}
}

// Index → GET /admin/v1/appearance (katalog: paginasi + search, aktif disematkan).
func (ctl *AppearanceController) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	search := c.Query("search")
	category := c.Query("category")

	active := fetemplate.DefaultSlug
	if svc := ctl.settings(); svc != nil {
		if s, err := svc.Get(c.Request.Context()); err == nil {
			active = fetemplate.ResolveActive(s.FeTemplate)
		}
	}

	items, total := ctl.fe.Paginate(c.Request.Context(), search, category, page, appearancePerPage, active)
	lastPage := (total + appearancePerPage - 1) / appearancePerPage
	if lastPage < 1 {
		lastPage = 1
	}
	view.RenderView(c, "home/appearance", gin.H{
		"title": "Tampilan", "active": "appearance",
		"templates": items, "activeSlug": active, "search": search, "category": category,
		"categories": ctl.fe.Categories(c.Request.Context()),
		"page":       page, "lastPage": lastPage, "total": total,
	})
}

// Apply → POST /admin/v1/appearance. Validasi + unduh (Ensure) + set aktif.
func (ctl *AppearanceController) Apply(c *gin.Context) {
	slug := c.PostForm("template")
	sess := sessions.Default(c)

	if !fetemplate.IsValidSlug(slug) {
		flash(sess, middleware.FlashErrorKey, "Template tidak dikenali.")
		c.Redirect(http.StatusFound, "/admin/v1/appearance")
		return
	}
	// Unduh on-demand (builtin → no-op). Gagal unduh → jangan ganti.
	if err := ctl.fe.Ensure(c.Request.Context(), slug); err != nil {
		flash(sess, middleware.FlashErrorKey, "Gagal menyiapkan template: "+slug)
		c.Redirect(http.StatusFound, "/admin/v1/appearance")
		return
	}

	svc := ctl.settings()
	if svc == nil {
		flash(sess, middleware.FlashErrorKey, "Layanan setting belum siap.")
		c.Redirect(http.StatusFound, "/admin/v1/appearance")
		return
	}
	if _, err := svc.Update(c.Request.Context(), settingdto.UpdateSettingInput{FeTemplate: slug}, ""); err != nil {
		flash(sess, middleware.FlashErrorKey, "Gagal mengganti template.")
		c.Redirect(http.StatusFound, "/admin/v1/appearance")
		return
	}
	flash(sess, middleware.FlashSuccessKey, "Template landing diganti ke '"+slug+"'.")
	c.Redirect(http.StatusFound, "/admin/v1/appearance")
}

// flash menulis pesan one-shot ke sesi (middleware.Flash memindahkannya ke view).
func flash(sess sessions.Session, key, msg string) {
	sess.Set(key, msg)
	_ = sess.Save()
}
