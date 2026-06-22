package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseHandler menstandarkan bentuk respons JSON API (padanan ResponseHandler
// di core NodeAdmin): { success, message, data }.
type apiEnvelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// JSONSuccess mengirim respons sukses standar.
func JSONSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, apiEnvelope{Success: true, Message: message, Data: data})
}

// JSONError mengirim respons error standar (dipakai middleware ErrorHandler).
func JSONError(c *gin.Context, status int, message string, errs interface{}) {
	c.JSON(status, apiEnvelope{Success: false, Message: message, Errors: errs})
}

// OK = 200 sukses.
func OK(c *gin.Context, message string, data interface{}) {
	JSONSuccess(c, http.StatusOK, message, data)
}

// Created = 201 sukses.
func Created(c *gin.Context, message string, data interface{}) {
	JSONSuccess(c, http.StatusCreated, message, data)
}
