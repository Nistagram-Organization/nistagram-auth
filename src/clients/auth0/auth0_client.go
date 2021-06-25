package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
	"strings"
)

type Auth0Client interface {
	obtainApiToken() (string, error)
	setUserRole(string, string, string) error
	RegisterUserOnAuth0(email string, password string, role string) (string, error)
}

type ApiTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	UserId string `json:"user_id"`
}

type SetUserRolesRequest struct {
	Roles []string `json:"roles"`
}

type auth0Client struct {
	restyClient  *resty.Client
	domain       string
	clientId     string
	clientSecret string
	audience     string
}

func NewAuth0Client(domain string, clientId string, clientSecret string, audience string) Auth0Client {
	return &auth0Client{
		resty.New(),
		domain,
		clientId,
		clientSecret,
		audience,
	}
}

func (s *auth0Client) obtainApiToken() (string, error) {
	endpoint := fmt.Sprintf("%soauth/token", s.domain)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.clientId)
	data.Set("client_secret", s.clientSecret)
	data.Set("audience", s.audience)

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	apiTokenResponse := &ApiTokenResponse{}
	json.NewDecoder(res.Body).Decode(apiTokenResponse)

	return apiTokenResponse.AccessToken, nil
}

func (s *auth0Client) setUserRole(userId string, role string, apiToken string) error {
	url := fmt.Sprintf("%sapi/v2/users/%s/roles", s.domain, userId)

	b, _ := json.Marshal(&SetUserRolesRequest{[]string{role}})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (s *auth0Client) RegisterUserOnAuth0(email string, password string, role string) (string, error) {
	apiToken, err := s.obtainApiToken()

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%sapi/v2/users", s.domain)

	b, _ := json.Marshal(&UserRegistrationRequest{Email: email, Password: password})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	userRegistrationResponse := &UserRegistrationResponse{}
	json.NewDecoder(res.Body).Decode(userRegistrationResponse)

	userId := userRegistrationResponse.UserId

	if err := s.setUserRole(userId, role, apiToken); err != nil {
		return "", nil
	}

	return userId, nil
}
