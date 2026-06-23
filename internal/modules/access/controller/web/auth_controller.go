package web

import (
	"net/http"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"goadmin/internal/modules/access/dto"
	accessmw "goadmin/internal/modules/access/middleware"
	"goadmin/internal/modules/access/service"
	"goadmin/internal/view"
)

// AuthController menangani login/logout + register + reset password (OTP) jalur
// SESI WEB. Sukses login menyimpan user id di sesi → EnsureAuthenticatedWeb.
type AuthController struct {
	auth  service.IAuthService
	reset service.IPasswordResetService
	users service.IUserService
}

// NewAuthController merakit controller.
func NewAuthController(auth service.IAuthService, reset service.IPasswordResetService, users service.IUserService) *AuthController {
	return &AuthController{auth: auth, reset: reset, users: users}
}

// ShowRegister → GET /auth/register (publik).
func (ctl *AuthController) ShowRegister(c *gin.Context) {
	view.RenderView(c, "auth/register", gin.H{"title": "Daftar"})
}

// Register → POST /auth/register. Buat akun (Active, tanpa role) lalu ke login.
func (ctl *AuthController) Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirm := c.PostForm("password_confirmation")

	if password != confirm {
		setFlashError(sessions.Default(c), "Konfirmasi password tidak cocok.")
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}
	_, err := ctl.users.Store(c.Request.Context(), dto.CreateUserInput{
		Name: name, Email: email, Password: password, Status: "Active",
	}, "")
	if err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/auth/register")
		return
	}
	setFlashSuccess(sessions.Default(c), "Akun berhasil dibuat. Silakan masuk.")
	c.Redirect(http.StatusFound, "/auth/login")
}

// ShowLogin → GET /auth/login (publik). Bila sudah login, lempar ke dashboard.
func (ctl *AuthController) ShowLogin(c *gin.Context) {
	sess := sessions.Default(c)
	if uid, _ := sess.Get(accessmw.SessionUserKey).(string); uid != "" {
		c.Redirect(http.StatusFound, "/admin/v1/dashboard")
		return
	}
	view.RenderView(c, "auth/login", gin.H{
		"title": "Masuk",
	})
}

// Login → POST /auth/login. Verifikasi kredensial → set sesi → redirect dashboard.
func (ctl *AuthController) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := ctl.auth.Authenticate(c.Request.Context(), email, password)
	if err != nil {
		setFlashError(sessions.Default(c), "Email atau password salah.")
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}

	sess := sessions.Default(c)
	sess.Set(accessmw.SessionUserKey, user.ID)
	_ = sess.Save()
	c.Redirect(http.StatusFound, "/admin/v1/dashboard")
}

// Logout → GET /auth/logout. Hapus sesi → kembali ke login.
func (ctl *AuthController) Logout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	_ = sess.Save()
	c.Redirect(http.StatusFound, "/auth/login")
}

// ShowForgot → GET /admin/v1/auth/reset/req (form minta OTP).
func (ctl *AuthController) ShowForgot(c *gin.Context) {
	view.RenderView(c, "auth/forgot", gin.H{"title": "Lupa Password"})
}

// Forgot → POST /admin/v1/auth/reset/request. Kirim OTP (bila email terdaftar) lalu ke form reset.
func (ctl *AuthController) Forgot(c *gin.Context) {
	email := c.PostForm("email")
	if err := ctl.reset.RequestReset(c.Request.Context(), email); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, "/admin/v1/auth/reset/req")
		return
	}
	setFlashSuccess(sessions.Default(c), "Jika email terdaftar, kode OTP telah dikirim.")
	c.Redirect(http.StatusFound, "/admin/v1/auth/reset/proc?email="+url.QueryEscape(email))
}

// ShowReset → GET /admin/v1/auth/reset/proc (form OTP + password baru).
func (ctl *AuthController) ShowReset(c *gin.Context) {
	view.RenderView(c, "auth/reset", gin.H{"title": "Reset Password", "email": c.Query("email")})
}

// Reset → POST /admin/v1/auth/reset/process. Verifikasi OTP → set password → ke login.
func (ctl *AuthController) Reset(c *gin.Context) {
	email := c.PostForm("email")
	otp := c.PostForm("otp")
	password := c.PostForm("password")
	confirm := c.PostForm("password_confirmation")

	back := "/admin/v1/auth/reset/proc?email=" + url.QueryEscape(email)
	if password != confirm {
		setFlashError(sessions.Default(c), "Konfirmasi password tidak cocok.")
		c.Redirect(http.StatusFound, back)
		return
	}
	if err := ctl.reset.Reset(c.Request.Context(), email, otp, password); err != nil {
		setFlashError(sessions.Default(c), errMessage(err))
		c.Redirect(http.StatusFound, back)
		return
	}
	setFlashSuccess(sessions.Default(c), "Password berhasil direset. Silakan masuk.")
	c.Redirect(http.StatusFound, "/auth/login")
}
