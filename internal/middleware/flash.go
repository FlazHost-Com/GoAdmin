package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Kunci flash di sesi (one-shot feedback pasca-redirect, pola PRG).
const (
	FlashSuccessKey = "flash_success"
	FlashErrorKey   = "flash_error"
)

// Flash memindahkan pesan flash dari sesi ke context (sekali pakai → dihapus
// dari sesi). RenderView lalu menaruhnya ke locals (`flash_success`/`flash_error`)
// agar chrome/halaman menampilkannya. Hanya jalur web (butuh sesi).
func Flash() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		changed := false
		if v, ok := sess.Get(FlashSuccessKey).(string); ok && v != "" {
			c.Set(FlashSuccessKey, v)
			sess.Delete(FlashSuccessKey)
			changed = true
		}
		if v, ok := sess.Get(FlashErrorKey).(string); ok && v != "" {
			c.Set(FlashErrorKey, v)
			sess.Delete(FlashErrorKey)
			changed = true
		}
		if changed {
			_ = sess.Save()
		}
		c.Next()
	}
}
