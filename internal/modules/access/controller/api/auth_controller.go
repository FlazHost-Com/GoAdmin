// Package api berisi controller REST modul access. Controller TIPIS: parse
// input + panggil service + format respons. Error diteruskan via c.Error()
// ke middleware terpusat (tanpa try/catch manual / mapping status di sini).
package api

import (
	"github.com/gin-gonic/gin"

	"goadmin/internal/auth"
	apperr "goadmin/internal/errors"
	"goadmin/internal/helpers"
	accessmw "goadmin/internal/modules/access/middleware"
	"goadmin/internal/modules/access/service"
)

// AuthController menangani login/logout API (JWT).
type AuthController struct {
	auths service.IAuthService
}

// NewAuthController merakit controller (service di-inject, bukan di-new di sini).
func NewAuthController(auths service.IAuthService) *AuthController {
	return &AuthController{auths: auths}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login → POST /api/v1/auth/login.
func (ctl *AuthController) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperr.Validation("Input tidak valid", nil))
		return
	}
	token, user, err := ctl.auths.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "Login berhasil", gin.H{"token": token, "user": user})
}

// Logout → POST /api/v1/auth/logout (butuh AuthenticatedJWT).
// Mencabut token saat ini (blacklist) → akses berikutnya 401.
func (ctl *AuthController) Logout(c *gin.Context) {
	v, ok := c.Get("jwt_claims")
	if !ok {
		c.Error(apperr.Unauthorized("Token tidak ada"))
		return
	}
	claims, ok := v.(*auth.Claims)
	if !ok || claims.ExpiresAt == nil {
		c.Error(apperr.Unauthorized("Token tidak valid"))
		return
	}
	if err := ctl.auths.Logout(c.Request.Context(), claims.ID, claims.ExpiresAt.Time); err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "Logout berhasil", nil)
}

// Me → GET /api/v1/auth/me (profil user terautentikasi).
func (ctl *AuthController) Me(c *gin.Context) {
	user := accessmw.UserFrom(c)
	if user == nil {
		c.Error(apperr.Unauthorized("Belum terautentikasi"))
		return
	}
	helpers.OK(c, "OK", user)
}
