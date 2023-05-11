package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}
