package base

import (
	"admin-go/cmd/app/handlers/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	http.Re(c, 0, "success", nil)
}
