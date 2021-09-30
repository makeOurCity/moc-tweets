package tweets

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type OrionClient struct {
	httpClient *http.Client
	baseURL    string
	clientID   string
	userPoolID string
}

func NewOrionClient(baseURL, clientID, userPoolID string) *OrionClient {
	o := OrionClient{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		clientID:   clientID,
		userPoolID: userPoolID,
	}

	return &o
}

func (o *OrionClient) Login(username string, password string) (string, error) {
	svc := cognitoidentityprovider.New(
		session.New(),
		&aws.Config{Region: aws.String("ap-northeast-1")},
	)

	authParams := map[string]*string{
		"USERNAME": aws.String(username),
		"PASSWORD": aws.String(password),
	}

	params := &cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow:       aws.String("ADMIN_NO_SRP_AUTH"),
		AuthParameters: authParams,
		ClientId:       aws.String(o.clientID),
		UserPoolId:     aws.String(o.userPoolID),
	}

	resp, err := svc.AdminInitiateAuth(params)
	if err != nil {
		return "", fmt.Errorf("svc.AdminInitiateAuth got error: %w", err)
	}

	token := *resp.AuthenticationResult.AccessToken

	return token, nil
}
