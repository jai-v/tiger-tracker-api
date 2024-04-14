package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AppController interface {
	HealthCheck(ctx *gin.Context)
}

type appController struct {
}

func NewAppController() AppController {
	return appController{}
}

// Get API health
// @Tags Health
// @Summary Checks API Service health
// @Description Confirms if the API Service is up and running
// @Produce plain
// @Success 200
// @Router /api/tiger-tracker/v1/health [get]
func (h appController) HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Service is up and running")
}
