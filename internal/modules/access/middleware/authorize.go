package middleware

import (
	"github.com/gin-gonic/gin"

	apperr "goadmin/internal/errors"
)

// Authorize memastikan user terotentikasi memiliki permission tertentu (RBAC).
// HARUS dipasang SETELAH AuthenticatedJWT/sesi. Administrator bypass (HasAccess
// mengembalikan true untuk admin).
func Authorize(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFrom(c)
		if user == nil {
			c.Error(apperr.Unauthorized("Belum terautentikasi"))
			c.Abort()
			return
		}
		if !user.HasAccess(permission) {
			c.Error(apperr.Forbidden("Anda tidak memiliki izin: " + permission))
			c.Abort()
			return
		}
		c.Next()
	}
}
