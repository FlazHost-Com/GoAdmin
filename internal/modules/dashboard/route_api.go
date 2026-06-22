package dashboard

import (
	"goadmin/internal/router"

	accessmw "goadmin/internal/modules/access/middleware"
	apictl "goadmin/internal/modules/dashboard/controller/api"
)

// registerAPIRoutes memasang endpoint statistik dashboard. Cukup terautentikasi
// (JWT) — statistik ringkas untuk admin yang sudah login.
func registerAPIRoutes(ctx *router.RegistrationContext, guard *accessmw.Guard, ctl *apictl.DashboardController) {
	g := ctx.API.Group("/v1/dashboard", guard.AuthenticatedJWT())

	g.GET("/stats", ctl.Stats)
	router.Register("api.v1.dashboard.stats", "/api/v1/dashboard/stats")
}
