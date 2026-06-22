package home

import (
	"goadmin/internal/router"

	accessmw "goadmin/internal/modules/access/middleware"
	accesssvc "goadmin/internal/modules/access/service"
	webctl "goadmin/internal/modules/home/controller/web"
	"goadmin/internal/modules/home/fetemplate"
	"goadmin/internal/modules/home/service"
	settingsvc "goadmin/internal/modules/setting/service"
)

// registerWebRoutes memasang landing publik + halaman admin pemilih template.
func registerWebRoutes(ctx *router.RegistrationContext, svc service.IHomeService, fe *fetemplate.Service) {
	c := ctx.Container

	// --- Landing publik ('/' render langsung; '/home' alias) ---
	ctl := webctl.NewHomeController(svc, fe)
	ctx.Web.GET("/", ctl.Index)
	router.Register("web.home.root", "/")
	ctx.Web.GET("/home", ctl.Index)
	router.Register("web.home.index", "/home")

	// --- Admin: pemilih template (membaca/menyetel Setting.fe_template, lazy) ---
	authSvc, _ := c.Resolve("access.IAuthService").(accesssvc.IAuthService)
	settingsProvider := func() settingsvc.ISettingService {
		s, _ := c.Resolve("setting.ISettingService").(settingsvc.ISettingService)
		return s
	}
	appCtl := webctl.NewAppearanceController(settingsProvider, fe)

	jwtless := accessmw.NewGuardWebOnly(authSvc)
	admin := ctx.Web.Group("/admin/v1")
	admin.Use(jwtless.EnsureAuthenticatedWeb("/auth/login"))

	admin.GET("/appearance", accessmw.AuthorizeWeb("setting.view"), appCtl.Index)
	router.Register("admin.v1.appearance.index", "/admin/v1/appearance")
	admin.POST("/appearance", accessmw.AuthorizeWeb("setting.update"), appCtl.Apply)
	// Proxy pratinjau template (iframe thumbnail/modal).
	admin.GET("/appearance/preview/:slug", accessmw.AuthorizeWeb("setting.view"), ctl.Preview)
	router.Register("admin.v1.appearance.preview", "/admin/v1/appearance/preview/:slug")
}
