package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}
