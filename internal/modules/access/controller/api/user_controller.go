package api

import (
	"github.com/gin-gonic/gin"

	apperr "goadmin/internal/errors"
	"goadmin/internal/helpers"
	accessmw "goadmin/internal/modules/access/middleware"
	"goadmin/internal/modules/access/dto"
	"goadmin/internal/modules/access/service"
)

// UserController = REST CRUD user.
type UserController struct {
	users service.IUserService
}

// NewUserController merakit controller.
func NewUserController(users service.IUserService) *UserController {
	return &UserController{users: users}
}

// Index → GET /api/v1/users.
func (ctl *UserController) Index(c *gin.Context) {
	var q dto.ListQuery
	_ = c.ShouldBindQuery(&q)
	res, err := ctl.users.Index(c.Request.Context(), q)
	if err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "OK", res)
}

// Show → GET /api/v1/users/:id.
func (ctl *UserController) Show(c *gin.Context) {
	user, err := ctl.users.Show(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "OK", user)
}

// Store → POST /api/v1/users.
func (ctl *UserController) Store(c *gin.Context) {
	var in dto.CreateUserInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(apperr.Validation("Input tidak valid", nil))
		return
	}
	actor := actorID(c)
	user, err := ctl.users.Store(c.Request.Context(), in, actor)
	if err != nil {
		c.Error(err)
		return
	}
	helpers.Created(c, "User dibuat", user)
}

// Update → PUT /api/v1/users/:id.
func (ctl *UserController) Update(c *gin.Context) {
	var in dto.UpdateUserInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.Error(apperr.Validation("Input tidak valid", nil))
		return
	}
	user, err := ctl.users.Update(c.Request.Context(), c.Param("id"), in, actorID(c))
	if err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "User diperbarui", user)
}

// Destroy → DELETE /api/v1/users/:id.
func (ctl *UserController) Destroy(c *gin.Context) {
	if err := ctl.users.Destroy(c.Request.Context(), c.Param("id")); err != nil {
		c.Error(err)
		return
	}
	helpers.OK(c, "User dihapus", nil)
}

func actorID(c *gin.Context) string {
	if u := accessmw.UserFrom(c); u != nil {
		return u.ID
	}
	return ""
}
