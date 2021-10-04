package tweets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ChimeraCoder/anaconda"
	cognitosrp "github.com/alexrudd/cognito-srp/v4"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cip "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type OrionClient struct {
	httpClient  *http.Client
	svc         *cip.Client
	baseURL     string
	clientID    string
	userPoolID  string
	serviceName string
	token       string
}

func NewOrionClient(baseURL, clientID, userPoolID, serviceName string) *OrionClient {
	cfg, _ := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
	)
	svc := cip.NewFromConfig(cfg)

	o := OrionClient{
		httpClient:  &http.Client{},
		baseURL:     baseURL,
		clientID:    clientID,
		userPoolID:  userPoolID,
		serviceName: serviceName,
		svc:         svc,
	}

	return &o
}

func (o *OrionClient) Login(username string, password string) error {
	csrp, _ := cognitosrp.NewCognitoSRP(username, password, o.userPoolID, o.clientID, nil)

	params := cip.InitiateAuthInput{
		AuthFlow:       types.AuthFlowTypeUserSrpAuth,
		AuthParameters: csrp.GetAuthParams(),
		ClientId:       aws.String(o.clientID),
		// UserPoolId:     aws.String(o.userPoolID),
	}

	resp, err := o.svc.InitiateAuth(context.Background(), &params)
	if err != nil {
		return fmt.Errorf("svc.InitiateAuth got error: %w", err)
	}

	if resp.ChallengeName == types.ChallengeNameTypePasswordVerifier {
		challengeResponses, _ := csrp.PasswordVerifierChallenge(resp.ChallengeParameters, time.Now())

		resp, err := o.svc.RespondToAuthChallenge(context.Background(), &cip.RespondToAuthChallengeInput{
			ChallengeName:      types.ChallengeNameTypePasswordVerifier,
			ChallengeResponses: challengeResponses,
			ClientId:           aws.String(csrp.GetClientId()),
		})
		if err != nil {
			panic(err)
		}

		// print the tokens
		fmt.Printf("Access Token: %s\n", *resp.AuthenticationResult.AccessToken)
		fmt.Printf("ID Token: %s\n", *resp.AuthenticationResult.IdToken)
		fmt.Printf("Refresh Token: %s\n", *resp.AuthenticationResult.RefreshToken)
		o.token = *resp.AuthenticationResult.IdToken
	} else {
		return fmt.Errorf(
			"challenge name is not %s but %s. Please sign up and change password your user",
			types.ChallengeNameTypePasswordVerifier,
			resp.ChallengeName,
		)
	}

	return nil
}

func (o *OrionClient) Send(e OrionEntity) (*http.Response, error) {
	url := fmt.Sprintf("%s/v2/entities", o.baseURL)
	req, err := o.getRequest(url, o.token, o.serviceName, e)
	if err != nil {
		return nil, fmt.Errorf("o.getRequest got error: %w", err)
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("o.httpClient.Do got error: %w", err)
	}

	return resp, nil
}

func (o *OrionClient) IsExistsEntity(t anaconda.Tweet) (bool, error) {
	id := GenerateID(t)
	url := fmt.Sprintf("%s/v2/entities/%s", o.baseURL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("http.NewRequest got error: %w", err)
	}

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("o.httpClient.Do got error: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

func (o *OrionClient) getRequest(url string, token string, fiwareServiceName string, data interface{}) (*http.Request, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal got error: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest got error: %w", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Fiware-Service", fiwareServiceName)

	return req, nil
}
