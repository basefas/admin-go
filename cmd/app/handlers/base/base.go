package base

import (
	"go-admin/cmd/app/handlers"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	handlers.Re(c, 0, "success", nil)
}
