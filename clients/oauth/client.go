package oauth

import (
	"github.com/gin-gonic/gin"
	hydraClient "github.com/ory/hydra-client-go/client"
	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"net/url"
	"tiger-tracker-api/logging"
)

type OauthClient interface {
	AcceptLoginRequest(ctx *gin.Context, userId string, rememberMe bool, loginChallenge string) (hydraAdmin.AcceptLoginRequestOK, error)
	IntrospectAccessToken(ctx *gin.Context, accessToken string) (models.OAuth2TokenIntrospection, error)
}

type oauthClient struct {
	hydraAdminService hydraAdmin.ClientService
}

func NewOauthClient(adminBaseUrl string) OauthClient {
	adminURL, _ := url.Parse(adminBaseUrl)
	hydraAdminHttpClient := hydraClient.NewHTTPClientWithConfig(nil,
		&hydraClient.TransportConfig{
			Schemes:  []string{adminURL.Scheme},
			Host:     adminURL.Host,
			BasePath: adminURL.Path,
		},
	)
	return oauthClient{
		hydraAdminService: hydraAdminHttpClient.Admin,
	}
}

func (oc oauthClient) AcceptLoginRequest(ctx *gin.Context, userId string, rememberMe bool, loginChallenge string) (hydraAdmin.AcceptLoginRequestOK, error) {
	logger := logging.GetLogger().WithField("Package", "OauthClient").WithField("Method", "AcceptLoginRequest")
	loginAcceptParam := hydraAdmin.NewAcceptLoginRequestParams()
	loginAcceptParam.WithContext(ctx)
	loginAcceptParam.SetLoginChallenge(loginChallenge)

	subject := userId
	loginAcceptParam.SetBody(&models.AcceptLoginRequest{
		Subject:  &subject,
		Remember: rememberMe,
	})

	resp, err := oc.hydraAdminService.AcceptLoginRequest(loginAcceptParam)
	if err != nil {
		logger.Errorf("failed to accept login request, error: %v", err)
		return hydraAdmin.AcceptLoginRequestOK{}, err
	}
	logger.Info("accepted login request for user_id: %s", userId)
	return *resp, err
}

func (oc oauthClient) IntrospectAccessToken(ctx *gin.Context, accessToken string) (models.OAuth2TokenIntrospection, error) {
	logger := logging.GetLogger().WithField("Package", "OauthClient").WithField("Method", "IntrospectAccessToken")
	introspectParam := hydraAdmin.NewIntrospectOAuth2TokenParams()
	introspectParam.WithContext(ctx)
	introspectParam.SetToken(accessToken)

	resp, err := oc.hydraAdminService.IntrospectOAuth2Token(introspectParam)
	if err != nil {
		logger.Errorf("failed to introspect token, error: %v", err)
		return models.OAuth2TokenIntrospection{}, err
	}
	return *resp.GetPayload(), nil
}
