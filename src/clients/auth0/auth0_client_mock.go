package auth0

import (
	"github.com/stretchr/testify/mock"
)

type Auth0ClientMock struct {
	mock.Mock
}

func (a *Auth0ClientMock) RegisterUserOnAuth0(email string, password string, role string) (string, error) {
	args := a.Called(email, password, role)
	if args.Get(1) == nil {
		return args.Get(0).(string), nil
	}
	return "", args.Get(1).(error)
}

func (a *Auth0ClientMock) obtainApiToken() (string, error) {
	panic("implement me")
}

func (a *Auth0ClientMock) setUserRole(s string, s2 string, s3 string) error {
	panic("implement me")
}

func (a *Auth0ClientMock) getUserIdOnAuth0(email string, apiToken string) (string, error) {
	panic("implement me")
}

func (a *Auth0ClientMock) BlockUserOnAuth0(email string) error {
	panic("implement me")
}

