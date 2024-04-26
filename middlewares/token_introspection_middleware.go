package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tiger-tracker-api/clients/oauth"
	"tiger-tracker-api/constants"
	"tiger-tracker-api/logging"
)

type TokenIntrospectionMiddleware struct {
	oauthClient oauth.OauthClient
}

func (t TokenIntrospectionMiddleware) ValidateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logging.GetLogger().WithField("Package", "middleware").WithField("Method", "ValidateToken")
		accessToken, err := t.getAccessToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errorCode": "Unauthorized", "errorMessage": err.Error()})
			return
		}
		resp, err := t.oauthClient.IntrospectAccessToken(ctx, accessToken)
		if err != nil {
			logger.Errorf("failed to introspect access token, error:%v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errorCode": "InternalServerError", "errorMessage": "server-error"})
			return
		}
		if !*resp.Active {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errorCode": "Unauthorized", "errorMessage": "expired token"})
			return
		}
		ctx.Next()
		logger.Info("token is active, allowing to proceed further")
	}
}

func NewTokenIntrospectionMiddleware(oauthClient oauth.OauthClient) TokenIntrospectionMiddleware {
	return TokenIntrospectionMiddleware{oauthClient: oauthClient}
}

func (t TokenIntrospectionMiddleware) getTokenFromHeaders(ctx *gin.Context) (string, error) {
	bearerToken := ctx.Request.Header.Get(constants.AUTHORIZATION_HEADER_KEY)
	if bearerToken == "" {
		return "", errors.New("missing or empty bearer token")
	}
	if strings.HasPrefix(bearerToken, "Bearer") && len(strings.Split(bearerToken, " ")) != 2 {
		return "", errors.New("invalid bearer token")
	}
	token := strings.Split(bearerToken, " ")[1]
	return token, nil
}

func (t TokenIntrospectionMiddleware) getTokenFromCookie(ctx *gin.Context) (string, error) {
	return ctx.Cookie(constants.COOKIE_ACCESS_TOKEN)
}

func (t TokenIntrospectionMiddleware) getAccessToken(ctx *gin.Context) (string, error) {
	logger := logging.GetLogger().WithField("Package", "middleware").WithField("Method", "getAccessToken")
	if accessToken, bearerTokenErr := t.getTokenFromHeaders(ctx); bearerTokenErr == nil {
		logger.Info("using access token from header")
		return accessToken, nil
	} else if accessToken, cookieErr := t.getTokenFromCookie(ctx); cookieErr == nil {
		logger.Info("using access token from cookie")
		return accessToken, nil
	} else {
		logger.Warn(cookieErr.Error())
		return "", errors.New("invalid bearer/cookie header")
	}
}
