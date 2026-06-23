package middleware

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Kunci flash di sesi (one-shot feedback pasca-redirect, pola PRG).
const (
	FlashSuccessKey = "flash_success"
	FlashErrorKey   = "flash_error"
	FieldErrorsKey  = "field_errors" // JSON map field→pesan (validasi inline)
	FieldOldKey     = "field_old"    // JSON map field→nilai lama (repopulasi)
)

// SetFieldErrors menyimpan error per-field + nilai lama (old input) ke sesi untuk
// SATU redirect (pola PRG). Padanan `req.session.errors`/`req.session.old` +
// helper `getError`/`old` NodeAdmin — view menampilkan error inline & mengisi
// ulang form. Disimpan sebagai JSON (cookie-session aman tanpa gob-register).
func SetFieldErrors(sess sessions.Session, errs, old map[string]string) {
	if b, err := json.Marshal(errs); err == nil {
		sess.Set(FieldErrorsKey, string(b))
	}
	if b, err := json.Marshal(old); err == nil {
		sess.Set(FieldOldKey, string(b))
	}
	_ = sess.Save()
}

// Flash memindahkan pesan flash dari sesi ke context (sekali pakai → dihapus
// dari sesi). RenderView lalu menaruhnya ke locals (`flash_success`/`flash_error`
// + `errors`/`old`) agar chrome/halaman menampilkannya. Hanya jalur web (sesi).
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
		// Error per-field + old input (JSON) → parse ke map → context.
		if v, ok := sess.Get(FieldErrorsKey).(string); ok && v != "" {
			var m map[string]string
			if json.Unmarshal([]byte(v), &m) == nil {
				c.Set(FieldErrorsKey, m)
			}
			sess.Delete(FieldErrorsKey)
			changed = true
		}
		if v, ok := sess.Get(FieldOldKey).(string); ok && v != "" {
			var m map[string]string
			if json.Unmarshal([]byte(v), &m) == nil {
				c.Set(FieldOldKey, m)
			}
			sess.Delete(FieldOldKey)
			changed = true
		}
		if changed {
			_ = sess.Save()
		}
		c.Next()
	}
}
