package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"tiger-tracker-api/apiError"
	"tiger-tracker-api/clients/oauth"
	"tiger-tracker-api/constants"
	"tiger-tracker-api/controller/models"
	"tiger-tracker-api/logging"
	"tiger-tracker-api/repository"
	serviceModels "tiger-tracker-api/service/models"
)

//go:generate mockgen -source=./service.go -destination=./mocks/mock_service.go -package=mocks

type AppService interface {
	GetAllTigersWithRecentSightingsFirst(ctx *gin.Context, pageNumber, pageSize int) (serviceModels.ListTigersResponse, error)
	AuthenticateLoginByPassword(ctx *gin.Context, request models.LoginByPasswordRequest) (string, *apiError.APIError)
	CreateNewUser(ctx *gin.Context, request models.SignupRequest) (models.SignupResponse, *apiError.APIError)
}

type appService struct {
	appRepository repository.AppRepository
	oauthClient   oauth.OauthClient
}

func NewAppService(appRepository repository.AppRepository, oauthClient oauth.OauthClient) AppService {
	return appService{appRepository: appRepository, oauthClient: oauthClient}
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

func (service appService) AuthenticateLoginByPassword(ctx *gin.Context, request models.LoginByPasswordRequest) (string, *apiError.APIError) {
	logger := logging.GetLogger().WithField("Package", "Service").WithField("Method", "AuthenticateLoginByPassword")
	loginFailureRedirectTo := "http://localhost:5001/authentication/login?login_challenge=%s&error=%s"
	userDetails, err := service.appRepository.GetUserDetailsByUsername(ctx, request.Username)
	if err != nil {
		logger.Errorf("failed to get password hash for user:%s, error:%v", request.Username, err)
		return fmt.Sprintf(loginFailureRedirectTo, request.LoginChallenge, "invalid username/password"), nil
	}
	if compareErr := bcrypt.CompareHashAndPassword([]byte(userDetails.PasswordHash), []byte(request.Password)); compareErr != nil {
		logger.Errorf("password does not match")
		return fmt.Sprintf(loginFailureRedirectTo, request.LoginChallenge, "invalid username/password"), nil
	}

	resp, err := service.oauthClient.AcceptLoginRequest(ctx, strconv.FormatInt(userDetails.Id, 10), request.RememberMe, request.LoginChallenge)
	if err != nil {
		logger.Errorf("failed to accept login request, error:%v", err)
		return fmt.Sprintf(loginFailureRedirectTo, request.LoginChallenge, "service temporarily unavailable"), nil
	}
	logger.Infof("login accepted for userId: %d", userDetails.Id)
	redirectTo := *resp.GetPayload().RedirectTo
	return redirectTo, nil
}

func (service appService) CreateNewUser(ctx *gin.Context, request models.SignupRequest) (models.SignupResponse, *apiError.APIError) {
	logger := logging.GetLogger().WithField("Package", "Service").WithField("Method", "CreateNewUser")
	count, err := service.appRepository.CheckIfUsernameOrEmailExists(ctx, request.Username, request.Email)
	var signupResponse models.SignupResponse
	if err != nil {
		logger.Errorf("failed to check if username or email already exists, error:%v", err)
		return signupResponse, apiError.InternalServerError
	}
	if count > 0 {
		logger.Errorf("username or email already registered")
		return signupResponse, apiError.UsernameOrEmailAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), constants.PASSWORD_HASH_COST)
	if err != nil {
		logger.Error("failed to generate password hash, error: %v", err)
		return signupResponse, apiError.InternalServerError
	}

	userId, err := service.appRepository.CreateUser(ctx, request.Username, request.Email, string(hashedPassword))
	if err != nil {
		logger.Errorf("failed to create user, error: %v", err)
		return signupResponse, apiError.InternalServerError
	}
	logger.Infof("created user with id: %d", userId)
	signupResponse.Status = "success"
	return signupResponse, nil
}
