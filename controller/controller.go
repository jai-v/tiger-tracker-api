package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiger-tracker-api/apiError"
	"tiger-tracker-api/service"
	"tiger-tracker-api/service/models"
)

type AppController interface {
	HealthCheck(ctx *gin.Context)
	ListAllTigers(ctx *gin.Context)
}

type appController struct {
	appService service.AppService
}

func NewAppController(appService service.AppService) AppController {
	return appController{appService: appService}
}

// Get API health
// @Tags Health
// @Summary Checks API Service health
// @Description Confirms if the API Service is up and running
// @Produce plain
// @Success 200
// @Router /v1/health [get]
func (h appController) HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Service is up and running")
}

func (h appController) readQueryParamAsInt(ctx *gin.Context, name string) (int, error) {
	strValue, exists := ctx.GetQuery(name)
	if !exists {
		return 0, errors.New(fmt.Sprintf("%s query param is missing", name))
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return 0, err
	}

	return intValue, nil
}

// Get All Tigers
// @Tags Tigers
// @Summary Returns a list of all the tigers with their details.
// @Description It is a paginated endpoint. The tigers are sorted by the last time they were seen.
// @Produce json
// @Param pageNo query int true "page number"
// @Param pageSize query int true "page size"
// @Success 200 {object} models.ListTigersResponse
// @Failure 400 {object} apiError.APIError
// @Failure 500 {object} apiError.APIError
// @Router /v1/tigers [get]
func (h appController) ListAllTigers(ctx *gin.Context) {
	pageNo, err := h.readQueryParamAsInt(ctx, "pageNo")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError.BadRequestError)
		return
	}
	pageSize, err := h.readQueryParamAsInt(ctx, "pageSize")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError.BadRequestError)
		return
	}

	tigerSightings, err := h.appService.GetAllTigersWithRecentSightingsFirst(ctx, pageNo, pageSize)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, apiError.InternalServerError)
		return
	}
	response := models.ListTigersResponse{Tigers: tigerSightings}
	ctx.JSON(http.StatusOK, response)
	return
}
