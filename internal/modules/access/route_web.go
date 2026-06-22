package access

import (
	"time"

	"goadmin/internal/middleware"
	"goadmin/internal/router"

	webctl "goadmin/internal/modules/access/controller/web"
	accessmw "goadmin/internal/modules/access/middleware"
	"goadmin/internal/modules/access/service"
)

// loginLimiter meredam brute-force login: maks 5 percobaan / menit / IP.
var loginLimiter = middleware.NewRateLimiter(5, time.Minute)

// otpLimiter membatasi permintaan OTP reset: maks 3 / menit / IP.
var otpLimiter = middleware.NewRateLimiter(3, time.Minute)

// registerWebRoutes memasang route web modul access:
//   - Auth sesi PUBLIK: /auth/login (GET/POST), /auth/logout.
//   - Admin (butuh sesi): /admin/v1/users. Urutan middleware WAJIB:
//     EnsureAuthenticatedWeb → AuthorizeWeb.
func registerWebRoutes(ctx *router.RegistrationContext) {
	c := ctx.Container

	userSvc, _ := c.Resolve("access.IUserService").(service.IUserService)
	authSvc, _ := c.Resolve("access.IAuthService").(service.IAuthService)
	roleSvc, _ := c.Resolve("access.IRoleService").(service.IRoleService)
	permSvc, _ := c.Resolve("access.IPermissionService").(service.IPermissionService)
	resetSvc, _ := c.Resolve("access.IPasswordResetService").(service.IPasswordResetService)

	// --- Auth sesi (publik, tanpa guard) ---
	authCtl := webctl.NewAuthController(authSvc, resetSvc)
	ctx.Web.GET("/auth/login", authCtl.ShowLogin)
	router.Register("web.auth.login", "/auth/login")
	ctx.Web.POST("/auth/login", loginLimiter.Middleware(), authCtl.Login)
	ctx.Web.GET("/auth/logout", authCtl.Logout)
	router.Register("web.auth.logout", "/auth/logout")

	// Reset password via OTP email (publik; permintaan OTP di-rate-limit).
	ctx.Web.GET("/auth/forgot", authCtl.ShowForgot)
	router.Register("web.auth.forgot", "/auth/forgot")
	ctx.Web.POST("/auth/forgot", otpLimiter.Middleware(), authCtl.Forgot)
	ctx.Web.GET("/auth/reset", authCtl.ShowReset)
	router.Register("web.auth.reset", "/auth/reset")
	ctx.Web.POST("/auth/reset", authCtl.Reset)

	// --- Admin (butuh sesi) ---
	jwtless := accessmw.NewGuardWebOnly(authSvc)
	userCtl := webctl.NewUserController(userSvc, roleSvc, c.Storage)

	admin := ctx.Web.Group("/admin/v1")
	admin.Use(jwtless.EnsureAuthenticatedWeb("/auth/login"))

	// --- Users (CRUD web) ---
	admin.GET("/users", accessmw.AuthorizeWeb("user.view"), userCtl.Index)
	router.Register("admin.v1.users.index", "/admin/v1/users")
	admin.GET("/users/create", accessmw.AuthorizeWeb("user.create"), userCtl.Create)
	admin.POST("/users", accessmw.AuthorizeWeb("user.create"), userCtl.Store)
	admin.GET("/users/:id/edit", accessmw.AuthorizeWeb("user.update"), userCtl.Edit)
	admin.POST("/users/:id", accessmw.AuthorizeWeb("user.update"), userCtl.Update)
	admin.POST("/users/:id/delete", accessmw.AuthorizeWeb("user.delete"), userCtl.Destroy)

	// --- Roles (CRUD web) ---
	roleCtl := webctl.NewRoleController(roleSvc, permSvc)
	admin.GET("/roles", accessmw.AuthorizeWeb("role.view"), roleCtl.Index)
	router.Register("admin.v1.roles.index", "/admin/v1/roles")
	admin.GET("/roles/create", accessmw.AuthorizeWeb("role.create"), roleCtl.Create)
	admin.POST("/roles", accessmw.AuthorizeWeb("role.create"), roleCtl.Store)
	admin.GET("/roles/:id/edit", accessmw.AuthorizeWeb("role.update"), roleCtl.Edit)
	admin.POST("/roles/:id", accessmw.AuthorizeWeb("role.update"), roleCtl.Update)
	admin.POST("/roles/:id/delete", accessmw.AuthorizeWeb("role.delete"), roleCtl.Destroy)

	// --- Permissions (CRUD web) ---
	permCtl := webctl.NewPermissionController(permSvc)
	admin.GET("/permissions", accessmw.AuthorizeWeb("permission.view"), permCtl.Index)
	router.Register("admin.v1.permissions.index", "/admin/v1/permissions")
	admin.GET("/permissions/create", accessmw.AuthorizeWeb("permission.create"), permCtl.Create)
	admin.POST("/permissions", accessmw.AuthorizeWeb("permission.create"), permCtl.Store)
	admin.GET("/permissions/:id/edit", accessmw.AuthorizeWeb("permission.update"), permCtl.Edit)
	admin.POST("/permissions/:id", accessmw.AuthorizeWeb("permission.update"), permCtl.Update)
	admin.POST("/permissions/:id/delete", accessmw.AuthorizeWeb("permission.delete"), permCtl.Destroy)
}
