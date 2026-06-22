package access

import (
	"github.com/gin-gonic/gin"

	"goadmin/internal/router"

	apictl "goadmin/internal/modules/access/controller/api"
	accessmw "goadmin/internal/modules/access/middleware"
)

// apiDeps mengumpulkan controller API modul.
type apiDeps struct {
	auth *apictl.AuthController
	user *apictl.UserController
	role *apictl.RoleController
	perm *apictl.PermissionController
}

// registerAPIRoutes memasang seluruh endpoint REST modul access.
// Pola: route bernama (router.Register) + middleware authenticated→authorize.
func registerAPIRoutes(ctx *router.RegistrationContext, guard *accessmw.Guard, d apiDeps) {
	api := ctx.API.Group("/v1")

	// Auth (publik: login; terproteksi: logout, me).
	authGrp := api.Group("/auth")
	named(authGrp, "POST", "api.v1.auth.login", "/login", d.auth.Login)
	authGrp.POST("/logout", guard.AuthenticatedJWT(), d.auth.Logout)
	router.Register("api.v1.auth.logout", "/api/v1/auth/logout")
	authGrp.GET("/me", guard.AuthenticatedJWT(), d.auth.Me)
	router.Register("api.v1.auth.me", "/api/v1/auth/me")

	// Semua resource di bawah ini butuh auth + permission spesifik (RBAC).
	auth := guard.AuthenticatedJWT()

	users := api.Group("/users", auth)
	resource(users, "/api/v1/users", "api.v1.users", "user", d.user.Index, d.user.Show, d.user.Store, d.user.Update, d.user.Destroy)

	roles := api.Group("/roles", auth)
	resource(roles, "/api/v1/roles", "api.v1.roles", "role", d.role.Index, d.role.Show, d.role.Store, d.role.Update, d.role.Destroy)

	perms := api.Group("/permissions", auth)
	resource(perms, "/api/v1/permissions", "api.v1.permissions", "permission", d.perm.Index, d.perm.Show, d.perm.Store, d.perm.Update, d.perm.Destroy)
}

// named mendaftarkan satu route bernama (registry) + memasangnya ke grup.
func named(g *gin.RouterGroup, method, name, path string, h gin.HandlerFunc, mw ...gin.HandlerFunc) {
	handlers := append(mw, h)
	switch method {
	case "GET":
		g.GET(path, handlers...)
	case "POST":
		g.POST(path, handlers...)
	case "PUT":
		g.PUT(path, handlers...)
	case "DELETE":
		g.DELETE(path, handlers...)
	}
	router.Register(name, fullPath(g, path))
}

// resource memasang 5 endpoint CRUD standar dengan guard permission per-aksi
// (RBAC: <subject>.view / .create / .update / .delete; Administrator bypass).
func resource(g *gin.RouterGroup, basePath, nameBase, subject string,
	index, show, store, update, destroy gin.HandlerFunc) {

	g.GET("", accessmw.Authorize(subject+".view"), index)
	router.Register(nameBase+".index", basePath)

	g.GET("/:id", accessmw.Authorize(subject+".view"), show)
	router.Register(nameBase+".show", basePath+"/:id")

	g.POST("", accessmw.Authorize(subject+".create"), store)
	router.Register(nameBase+".store", basePath)

	g.PUT("/:id", accessmw.Authorize(subject+".update"), update)
	router.Register(nameBase+".update", basePath+"/:id")

	g.DELETE("/:id", accessmw.Authorize(subject+".delete"), destroy)
	router.Register(nameBase+".destroy", basePath+"/:id")
}

// fullPath menggabungkan base path grup dengan path relatif untuk registry.
func fullPath(g *gin.RouterGroup, path string) string {
	return g.BasePath() + path
}
