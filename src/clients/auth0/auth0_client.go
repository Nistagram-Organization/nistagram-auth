package auth0

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	USER_ROLE_ID  = "rol_pi1paAm945Clo77K"
	AGENT_ROLE_ID = "rol_CkudJNyTZ7IERarY"
)

type Auth0Client interface {
	obtainApiToken() (string, error)
	setUserRole(string, string, string) error
	getUserIdOnAuth0(email string, apiToken string) (string, error)
	RegisterUserOnAuth0(email string, password string, role string) (string, error)
	BlockUserOnAuth0(email string) error
}

type ApiTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserRegistrationRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Connection string `json:"connection"`
}

type UserRegistrationResponse struct {
	UserId string `json:"user_id"`
}

type SetUserRolesRequest struct {
	Roles []string `json:"roles"`
}

type GetUserIdResponse struct {
	UserId string `json:"user_id"`
}

type BlockUserRequest struct {
	Blocked bool `json:"blocked"`
}

type auth0Client struct {
	domain       string
	clientId     string
	clientSecret string
	audience     string
}

func NewAuth0Client(domain string, clientId string, clientSecret string, audience string) Auth0Client {
	return &auth0Client{
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

	var roleId string

	if role == "user" {
		roleId = USER_ROLE_ID
	} else {
		roleId = AGENT_ROLE_ID
	}

	b, _ := json.Marshal(&SetUserRolesRequest{[]string{roleId}})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		return errors.New("failed to assign role to user")
	}

	return nil
}

func (s *auth0Client) RegisterUserOnAuth0(email string, password string, role string) (string, error) {
	apiToken, err := s.obtainApiToken()

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%sapi/v2/users", s.domain)

	b, _ := json.Marshal(&UserRegistrationRequest{Email: email, Password: password, Connection: "nistagram-database"})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return "", errors.New("failed to create user on auth0")
	}

	userRegistrationResponse := &UserRegistrationResponse{}
	json.NewDecoder(res.Body).Decode(userRegistrationResponse)

	userId := userRegistrationResponse.UserId

	if err := s.setUserRole(userId, role, apiToken); err != nil {
		return "", err
	}

	return userId, nil
}

func (s *auth0Client) getUserIdOnAuth0(email string, apiToken string) (string, error) {
	url := fmt.Sprintf("%sapi/v2/users", s.domain)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	q := req.URL.Query()
	q.Add("fields", "user_id")
	q.Add("q", fmt.Sprintf("email:\"%s\" AND identities.connection:\"nistagram-database\"", email))
	q.Add("search_engine", "v3")

	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var userIdResponses []GetUserIdResponse
	if err := json.NewDecoder(res.Body).Decode(&userIdResponses); err != nil {
		return "", err
	}

	userId := userIdResponses[0].UserId

	return userId, nil
}

func (s *auth0Client) BlockUserOnAuth0(email string) error {
	apiToken, err := s.obtainApiToken()

	if err != nil {
		return err
	}

	userId, err := s.getUserIdOnAuth0(email, apiToken)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%sapi/v2/users/%s", s.domain, userId)

	b, _ := json.Marshal(&BlockUserRequest{Blocked: true})

	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("failed to block user")
	}

	return nil
}
