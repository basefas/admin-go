package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	InvalidArguments = errors.New("Invalid arguments")
	AuthError        = errors.New("auth error.")
)

func Health(c *gin.Context) {
	Re(c, 0, "success", nil)
}
