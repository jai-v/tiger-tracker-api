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

func (h healthCheckController) Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Service is up and running")
}
