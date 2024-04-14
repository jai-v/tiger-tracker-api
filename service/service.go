package service

import (
	"github.com/gin-gonic/gin"
	"tiger-tracker-api/apiError"
	"tiger-tracker-api/logging"
	"tiger-tracker-api/repository"
	"tiger-tracker-api/repository/models"
)

type AppService interface {
	GetAllTigersWithRecentSightingsFirst(ctx *gin.Context, pageNumber, pageSize int) ([]models.TigerDetailWithSightings, error)
}

type appService struct {
	appRepository repository.AppRepository
}

func NewAppService(appRepository repository.AppRepository) AppService {
	return appService{appRepository: appRepository}
}

func (service appService) GetAllTigersWithRecentSightingsFirst(ctx *gin.Context, pageNumber, pageSize int) ([]models.TigerDetailWithSightings, error) {
	logger := logging.GetLogger().WithField("Package", "Service").WithField("Method", "GetAllTigersInDescendingOrderOfLastSeen")
	tigerSightings, err := service.appRepository.GetRecentTigerSightings(ctx, pageNumber, pageSize)
	if err != nil {
		logger.Errorf("failed to fetch recent sightings of tiger from db, apiError:%v", err)
		return nil, apiError.InternalServerError
	}
	return tigerSightings, nil
}
