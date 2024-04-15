package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"testing"
	"tiger-tracker-api/apiError"
	"tiger-tracker-api/repository/mocks"
	"tiger-tracker-api/repository/models"
	models2 "tiger-tracker-api/service/models"
	"time"
)

type AppServiceTestSuite struct {
	suite.Suite
	context           *gin.Context
	mockController    *gomock.Controller
	mockAppRepository *mocks.MockAppRepository
	appService        AppService
}

func TestAppServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AppServiceTestSuite))
}

func (suite *AppServiceTestSuite) SetupSuite() {

	suite.mockController = gomock.NewController(suite.T())
	suite.mockAppRepository = mocks.NewMockAppRepository(suite.mockController)
	suite.appService = NewAppService(suite.mockAppRepository)
	suite.context, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.context.Request = httptest.NewRequest("GET", "/api", nil)
}

func (suite *AppServiceTestSuite) TestListALlTigers_WhenRepositoryReturnsResponse_ShouldReturnResponse() {
	mockRepositoryResponse := []models.TigerDetailWithSightings{{
		TigerId:      1,
		TigerName:    "raja",
		DateOfBirth:  time.Time{},
		LastSeen:     time.Time{},
		LastPosition: models.Location{Latitude: 21.789, Longitude: 91.238},
		ReportedBy:   1,
		Image:        nil,
	}}
	expectedResponse := models2.ListTigersResponse{Tigers: mockRepositoryResponse}
	suite.mockAppRepository.EXPECT().GetRecentTigerSightings(suite.context, 1, 5).Return(mockRepositoryResponse, nil)
	actualResponse, actualErr := suite.appService.GetAllTigersWithRecentSightingsFirst(suite.context, 1, 5)

	suite.Nil(actualErr)
	suite.Equal(actualResponse, expectedResponse)
}

func (suite *AppServiceTestSuite) TestListALlTigers_WhenRepositoryReturnsError_ShouldReturnAPIError() {
	suite.mockAppRepository.EXPECT().GetRecentTigerSightings(suite.context, 1, 5).Return(nil, errors.New("some-error"))
	actualResponse, actualErr := suite.appService.GetAllTigersWithRecentSightingsFirst(suite.context, 1, 5)

	suite.Equal(apiError.InternalServerError, actualErr)
	suite.Empty(actualResponse)
}
