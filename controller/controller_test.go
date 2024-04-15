package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"tiger-tracker-api/repository/models"
	serviceMocks "tiger-tracker-api/service/mocks"
	serviceModels "tiger-tracker-api/service/models"
	"time"
)

type AppControllerTestSuite struct {
	suite.Suite
	recorder       *httptest.ResponseRecorder
	context        *gin.Context
	mockController *gomock.Controller
	mockAppService *serviceMocks.MockAppService
	appController  AppController
}

func TestAppControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AppControllerTestSuite))
}

func (suite *AppControllerTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockAppService = serviceMocks.NewMockAppService(suite.mockController)
	suite.appController = NewAppController(suite.mockAppService)
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *AppControllerTestSuite) TestHealthCheck_ShouldReturnResponse() {
	suite.appController.HealthCheck(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal("Service is up and running", suite.recorder.Body.String())
}

func (suite *AppControllerTestSuite) TestListALlTigers_WhenRequestIsValidAndServiceReturnResponse_ShouldReturnSuccessResponse() {
	gin.SetMode(gin.TestMode)
	suite.context.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	suite.context.Request.Method = "GET"
	suite.context.Request.Header.Set("Content-Type", "application/json")
	u := url.Values{}
	u.Add("pageNo", "1")
	u.Add("pageSize", "5")
	suite.context.Request.URL.RawQuery = u.Encode()

	mockServiceResponse := serviceModels.ListTigersResponse{Tigers: []models.TigerDetailWithSightings{{
		TigerId:     1,
		TigerName:   "raja",
		DateOfBirth: time.Time{},
		LastSeen:    time.Time{},
		LastPosition: models.Location{
			Latitude:  21.789,
			Longitude: 91.239,
		},
		ReportedBy: 1,
		Image:      nil,
	}}}

	suite.mockAppService.EXPECT().GetAllTigersWithRecentSightingsFirst(suite.context, 1, 5).Return(mockServiceResponse, nil)
	suite.appController.ListAllTigers(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	expectedResponse := `{"tigers":[{"tigerName":"raja","dob":"0001-01-01T00:00:00Z","lastSeen":"0001-01-01T00:00:00Z","lastPosition":{"Latitude":21.789,"Longitude":91.239}}]}`
	suite.Equal(expectedResponse, suite.recorder.Body.String())
}

func (suite *AppControllerTestSuite) TestListALlTigers_WhenRequestIsValidAndServiceReturnsError_ShouldInternalServerErrorResponse() {
	gin.SetMode(gin.TestMode)
	suite.context.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	suite.context.Request.Method = "GET"
	suite.context.Request.Header.Set("Content-Type", "application/json")
	u := url.Values{}
	u.Add("pageNo", "1")
	u.Add("pageSize", "5")
	suite.context.Request.URL.RawQuery = u.Encode()

	suite.mockAppService.EXPECT().GetAllTigersWithRecentSightingsFirst(suite.context, 1, 5).Return(serviceModels.ListTigersResponse{}, errors.New("some-error"))
	suite.appController.ListAllTigers(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	expectedResponse := `{"errorCode":"ERR_INTERNAL_SERVER_ERROR","errorMessage":"Internal Server Error"}`
	suite.Equal(expectedResponse, suite.recorder.Body.String())
}

func (suite *AppControllerTestSuite) TestListALlTigers_WhenRequestIsMissingPageNoQueryParam_ShouldReturnBadRequestErrorResponse() {
	gin.SetMode(gin.TestMode)
	suite.context.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	suite.context.Request.Method = "GET"
	suite.context.Request.Header.Set("Content-Type", "application/json")
	u := url.Values{}
	u.Add("pageSize", "5")
	suite.context.Request.URL.RawQuery = u.Encode()

	suite.appController.ListAllTigers(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	expectedResponse := `{"errorCode":"ERR_BAD_REQUEST","errorMessage":"Invalid Request"}`
	suite.Equal(expectedResponse, suite.recorder.Body.String())
}

func (suite *AppControllerTestSuite) TestListALlTigers_WhenRequestIsHavingInvalidPageNoQueryParam_ShouldReturnBadRequestErrorResponse() {
	gin.SetMode(gin.TestMode)
	suite.context.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	suite.context.Request.Method = "GET"
	suite.context.Request.Header.Set("Content-Type", "application/json")
	u := url.Values{}
	u.Add("pageNo", "one")
	u.Add("pageSize", "5")
	suite.context.Request.URL.RawQuery = u.Encode()

	suite.appController.ListAllTigers(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	expectedResponse := `{"errorCode":"ERR_BAD_REQUEST","errorMessage":"Invalid Request"}`
	suite.Equal(expectedResponse, suite.recorder.Body.String())
}

func (suite *AppControllerTestSuite) TestListALlTigers_WhenRequestIsMissingPageSizeQueryParam_ShouldReturnBadRequestErrorResponse() {
	gin.SetMode(gin.TestMode)
	suite.context.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	suite.context.Request.Method = "GET"
	suite.context.Request.Header.Set("Content-Type", "application/json")
	u := url.Values{}
	u.Add("pageNo", "1")
	suite.context.Request.URL.RawQuery = u.Encode()

	suite.appController.ListAllTigers(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	expectedResponse := `{"errorCode":"ERR_BAD_REQUEST","errorMessage":"Invalid Request"}`
	suite.Equal(expectedResponse, suite.recorder.Body.String())
}

func (suite *AppControllerTestSuite) TestListALlTigers_WhenRequestIsHavingInvalidPageSizeQueryParam_ShouldReturnBadRequestErrorResponse() {
	gin.SetMode(gin.TestMode)
	suite.context.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}
	suite.context.Request.Method = "GET"
	suite.context.Request.Header.Set("Content-Type", "application/json")
	u := url.Values{}
	u.Add("pageNo", "1")
	u.Add("pageSize", "five")
	suite.context.Request.URL.RawQuery = u.Encode()

	suite.appController.ListAllTigers(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	expectedResponse := `{"errorCode":"ERR_BAD_REQUEST","errorMessage":"Invalid Request"}`
	suite.Equal(expectedResponse, suite.recorder.Body.String())
}
