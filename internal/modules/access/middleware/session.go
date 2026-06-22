package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	apperr "goadmin/internal/errors"
)

// SessionUserKey = kunci penyimpanan user id di sesi web.
const SessionUserKey = "user_id"

// EnsureAuthenticatedWeb memastikan ada sesi login (web). Bila tidak, redirect
// ke halaman login. Padanan ensureAuthenticated di NodeAdmin (jalur sesi).
func (g *Guard) EnsureAuthenticatedWeb(loginPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		uid, _ := sess.Get(SessionUserKey).(string)
		if uid == "" {
			c.Redirect(302, loginPath)
			c.Abort()
			return
		}
		user, err := g.auths.FindByID(c.Request.Context(), uid)
		if err != nil {
			sess.Clear()
			_ = sess.Save()
			c.Redirect(302, loginPath)
			c.Abort()
			return
		}
		c.Set(ctxUserKey, user)
		c.Next()
	}
}

// AuthorizeWeb = RBAC untuk jalur web (render 403 lewat error terpusat).
func AuthorizeWeb(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := UserFrom(c)
		if user == nil || !user.HasAccess(permission) {
			c.Error(apperr.Forbidden("Anda tidak memiliki izin: " + permission))
			c.Abort()
			return
		}
		c.Next()
	}
}
