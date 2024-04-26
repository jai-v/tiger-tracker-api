package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiger-tracker-api/apiError"
	"tiger-tracker-api/controller/models"
	"tiger-tracker-api/logging"
	"tiger-tracker-api/service"
)

type AppController interface {
	HealthCheck(ctx *gin.Context)
	ListAllTigers(ctx *gin.Context)
	LoginByPassword(ctx *gin.Context)
	Signup(ctx *gin.Context)
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
func (ac appController) HealthCheck(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Service is up and running")
}

func (ac appController) readQueryParamAsInt(ctx *gin.Context, name string) (int, error) {
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
// @Param Authorization header string true "starts with Bearer"
// @Success 200 {object} models.ListTigersResponse
// @Failure 400 {object} apiError.APIError
// @Failure 500 {object} apiError.APIError
// @Router /v1/tigers [get]
func (ac appController) ListAllTigers(ctx *gin.Context) {
	pageNo, err := ac.readQueryParamAsInt(ctx, "pageNo")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError.BadRequestError)
		return
	}
	pageSize, err := ac.readQueryParamAsInt(ctx, "pageSize")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError.BadRequestError)
		return
	}

	response, err := ac.appService.GetAllTigersWithRecentSightingsFirst(ctx, pageNo, pageSize)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, apiError.InternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, response)
	return
}

// Login By Password
// @Tags Login
// @Summary Identity provider api to verify login.
// @Description Accepts the login consent when credentials are correct and redirects to the consent page.
// @Accept application/x-www-form-urlencoded
// @Param login_challenge formData string true "login challenge"
// @Param user_name formData string true "login user name"
// @Param password formData string true "login password"
// @Produce json
// @Success 302
// @Failure 422 {object} apiError.APIError
// @Router /v1/login/password [post]
func (ac appController) LoginByPassword(ctx *gin.Context) {
	logger := logging.GetLogger().WithField("Package", "Controller").WithField("Method", "LoginByPassword")
	rememberMe, _ := strconv.ParseBool(ctx.Request.FormValue("remember_me"))
	req := models.LoginByPasswordRequest{
		LoginChallenge: ctx.Request.FormValue("login_challenge"),
		Username:       ctx.Request.FormValue("user_name"),
		Password:       ctx.Request.FormValue("password"),
		RememberMe:     rememberMe,
	}
	//TODO: validate request
	redirectTo, err := ac.appService.AuthenticateLoginByPassword(ctx, req)
	if err != nil {
		logger.Errorf("failed to authenticate, error: %s", err)
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	ctx.Redirect(http.StatusFound, redirectTo)
	return
}

// Signup
// @Tags Signup
// @Summary Creates a new user.
// @Description Creates a new user.
// @Produce json
// @Accept application/x-www-form-urlencoded
// @Param user_name formData string true "user name"
// @Param email formData string true "email"
// @Param password formData string true "password"
// @Success 201 {object} models.SignupResponse
// @Failure 500 {object} apiError.APIError
// @Router /v1/signup [post]
func (ac appController) Signup(ctx *gin.Context) {
	logger := logging.GetLogger().WithField("Package", "Controller").WithField("Method", "Signup")
	signupRequest := models.SignupRequest{
		Username: ctx.Request.FormValue("user_name"),
		Email:    ctx.Request.FormValue("email"),
		Password: ctx.Request.FormValue("password"),
	}

	//TODO: validate request
	signupResponse, apiErr := ac.appService.CreateNewUser(ctx, signupRequest)
	if apiErr != nil {
		logger.Errorf("failed to create new user, error: %s", apiErr)
		ctx.AbortWithStatusJSON(apiErr.HttpStatusCode, apiErr)
		return
	}
	ctx.JSON(http.StatusCreated, signupResponse)
}
