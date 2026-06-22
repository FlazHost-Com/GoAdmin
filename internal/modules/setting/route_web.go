package setting

import (
	"goadmin/internal/router"

	accessmw "goadmin/internal/modules/access/middleware"
	accesssvc "goadmin/internal/modules/access/service"
	webctl "goadmin/internal/modules/setting/controller/web"
	"goadmin/internal/modules/setting/service"
)

// registerWebRoutes memasang halaman admin setting (HTML, sesi). Hanya dipanggil
// mode full (ctx.Web != nil). Urutan: EnsureAuthenticatedWeb → AuthorizeWeb.
func registerWebRoutes(ctx *router.RegistrationContext) {
	c := ctx.Container

	authSvc, _ := c.Resolve("access.IAuthService").(accesssvc.IAuthService)
	settingSvc, _ := c.Resolve("setting.ISettingService").(service.ISettingService)

	jwtless := accessmw.NewGuardWebOnly(authSvc)
	ctl := webctl.NewSettingController(settingSvc, c.Storage)

	admin := ctx.Web.Group("/admin/v1")
	admin.Use(jwtless.EnsureAuthenticatedWeb("/auth/login"))

	admin.GET("/setting", accessmw.AuthorizeWeb("setting.view"), ctl.Index)
	router.Register("admin.v1.setting.index", "/admin/v1/setting")

	admin.POST("/setting", accessmw.AuthorizeWeb("setting.update"), ctl.Update)
	router.Register("admin.v1.setting.update", "/admin/v1/setting")
}
