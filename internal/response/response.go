package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Message: message, Data: data})
}

func OKWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{Success: true, Message: message, Data: data, Meta: meta})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Message: message, Data: data})
}

func NoContent(c *gin.Context) {
	c.JSON(http.StatusOK, Response{Success: true, Message: "deleted successfully"})
}

// ValidationError formats go-playground/validator errors into human-readable messages
func ValidationError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		BadRequest(c, err.Error())
		return
	}

	msgs := make([]string, 0, len(ve))
	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		switch fe.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", field))
		case "email":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid email address", field))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s characters", field, fe.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s characters", field, fe.Param()))
		case "len":
			msgs = append(msgs, fmt.Sprintf("%s must be exactly %s characters", field, fe.Param()))
		case "eqfield":
			msgs = append(msgs, fmt.Sprintf("%s must match %s", field, strings.ToLower(fe.Param())))
		case "min=1":
			msgs = append(msgs, fmt.Sprintf("%s must be greater than 0", field))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", field))
		}
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "validation failed",
		Data:    msgs,
	})
}

func BadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, Response{Success: false, Message: message})
}

func Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{Success: false, Message: message})
}

func NotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, Response{Success: false, Message: message})
}

func InternalError(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{Success: false, Message: message})
}
