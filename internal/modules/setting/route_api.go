package setting

import (
	"goadmin/internal/router"

	accessmw "goadmin/internal/modules/access/middleware"
	apictl "goadmin/internal/modules/setting/controller/api"
)

// registerAPIRoutes memasang endpoint REST setting (singleton: show + update),
// terproteksi JWT + RBAC (setting.view / setting.update; Administrator bypass).
func registerAPIRoutes(ctx *router.RegistrationContext, guard *accessmw.Guard, ctl *apictl.SettingController) {
	g := ctx.API.Group("/v1/setting", guard.AuthenticatedJWT())

	g.GET("", accessmw.Authorize("setting.view"), ctl.Show)
	router.Register("api.v1.setting.show", "/api/v1/setting")

	g.PUT("", accessmw.Authorize("setting.update"), ctl.Update)
	router.Register("api.v1.setting.update", "/api/v1/setting")
}
