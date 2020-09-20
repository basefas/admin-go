package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	InvalidArguments = errors.New("Invalid arguments")
	AuthError        = errors.New("auth error.")
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Re(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
