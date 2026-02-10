package handlers

import (
	"mk-pos-billing/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	utils.RespondJSON(c, http.StatusOK, "OK", gin.H{
		"status": "healthy",
	}, nil)
}
