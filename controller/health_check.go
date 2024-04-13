package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthCheckController interface {
	Status(ctx *gin.Context)
}

type healthCheckController struct {
}

func NewHealthCheckController() HealthCheckController {
	return healthCheckController{}
}

// Get API health
// @Tags Health
// @Summary Checks API Service health
// @Description Confirms if the API Service is up and running
// @Produce plain
// @Success 200
// @Router /api/tiger-tracker/v1/health [get]
func (h healthCheckController) Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Service is up and running")
}
