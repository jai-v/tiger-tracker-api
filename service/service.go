package service

import (
	"github.com/gin-gonic/gin"
	"tiger-tracker-api/apiError"
	"tiger-tracker-api/logging"
	"tiger-tracker-api/repository"
	serviceModels "tiger-tracker-api/service/models"
)

//go:generate mockgen -source=./service.go -destination=./mocks/mock_service.go -package=mocks

type AppService interface {
	GetAllTigersWithRecentSightingsFirst(ctx *gin.Context, pageNumber, pageSize int) (serviceModels.ListTigersResponse, error)
}

type appService struct {
	appRepository repository.AppRepository
}

func NewAppService(appRepository repository.AppRepository) AppService {
	return appService{appRepository: appRepository}
}

func (service appService) GetAllTigersWithRecentSightingsFirst(ctx *gin.Context, pageNumber, pageSize int) (serviceModels.ListTigersResponse, error) {
	logger := logging.GetLogger().WithField("Package", "Service").WithField("Method", "GetAllTigersWithRecentSightingsFirst")
	tigerSightings, err := service.appRepository.GetRecentTigerSightings(ctx, pageNumber, pageSize)
	if err != nil {
		logger.Errorf("failed to fetch recent sightings of tiger from db, apiError:%v", err)
		return serviceModels.ListTigersResponse{}, apiError.InternalServerError
	}
	response := serviceModels.ListTigersResponse{Tigers: tigerSightings}
	return response, nil
}
