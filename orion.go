package tweets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type OrionClient struct {
	httpClient  *http.Client
	baseURL     string
	clientID    string
	userPoolID  string
	serviceName string
	token       string
}

func NewOrionClient(baseURL, clientID, userPoolID, serviceName string) *OrionClient {
	o := OrionClient{
		httpClient:  &http.Client{},
		baseURL:     baseURL,
		clientID:    clientID,
		userPoolID:  userPoolID,
		serviceName: serviceName,
	}

	return &o
}

func (o *OrionClient) Login(username string, password string) error {
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
		return fmt.Errorf("svc.AdminInitiateAuth got error: %w", err)
	}

	token := *resp.AuthenticationResult.AccessToken

	return nil
}

func (o *OrionClient) Send(e OrionEntity) error {
	url := fmt.Sprintf("%s/v2/entities", o.baseURL)
	req, err := o.getRequest(url, o.token, o.serviceName, e)
	if err != nil {
		return fmt.Errorf("o.getRequest got error: %w", err)
	}
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
