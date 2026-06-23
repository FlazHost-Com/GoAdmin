package middleware

import (
	"github.com/gin-gonic/gin"

	apperr "goadmin/internal/errors"
	"goadmin/internal/router"
)

// Authorize (RBAC route-driven, a la NodeAdmin AccessMiddleware): menurunkan
// NAMA route + METHOD dari request berjalan, lalu memeriksa user memiliki
// permission (name+method) tsb. TANPA argumen subjek — granularitas per-route.
// HARUS dipasang SETELAH AuthenticatedJWT. Administrator bypass.
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFrom(c)
		if user == nil {
			c.Error(apperr.Unauthorized("Belum terautentikasi"))
			c.Abort()
			return
		}
		// Nama route diturunkan dari (method, pola-path); "" bila route tak
		// bernama → HasAccess false utk non-admin (admin tetap bypass).
		name := router.NameByMethodPath(c.Request.Method, c.FullPath())
		if !user.HasAccess(name, c.Request.Method) {
			c.Error(apperr.Forbidden("Anda tidak memiliki izin untuk aksi ini"))
			c.Abort()
			return
		}
		c.Next()
	}
}
