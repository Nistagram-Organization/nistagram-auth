package auth

import (
	"errors"
	"github.com/Nistagram-Organization/nistagram-auth/src/clients/auth0"
	"github.com/Nistagram-Organization/nistagram-auth/src/clients/user_grpc_client"
	"github.com/Nistagram-Organization/nistagram-auth/src/dto/registration_request"
	"github.com/Nistagram-Organization/nistagram-shared/src/model/user"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthServiceUnitTestsSuite struct {
	suite.Suite
	userGrpcClientMock  *user_grpc_client.UserGrpcClientMock
	auth0GrpcClientMock *auth0.Auth0ClientMock
	service             AuthService
}

func TestAuthServiceUnitTestsSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceUnitTestsSuite))
}

func (suite *AuthServiceUnitTestsSuite) SetupSuite() {
	suite.userGrpcClientMock = new(user_grpc_client.UserGrpcClientMock)
	suite.auth0GrpcClientMock = new(auth0.Auth0ClientMock)
	suite.service = NewAuthService(suite.auth0GrpcClientMock, suite.userGrpcClientMock)
}

func (suite *AuthServiceUnitTestsSuite) TestNewAuthService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *AuthServiceUnitTestsSuite) TestAuthService_Register_InvalidRegistrationRequest() {
	registrationRequest := registration_request.RegistrationRequest{
		User: user.User{
			Name:      "Ime",
			OwnerID:   1,
			OwnerType: "registered_users",
			Username:  "",
			FirstName: "FirstName",
			LastName:  "LastName",
			Phone:     "02123456678",
			BirthDate: 123,
			Website:   "www.site.com",
			Biography: "Biography",
			Gender:    1,
			Email:     "mail@mail.com",
		},
		Password: "123123123",
		Role:     "user",
	}
	err := rest_error.NewBadRequestError("Username cannot be empty")

	_, regErr := suite.service.Register(registrationRequest)

	assert.Equal(suite.T(), err, regErr)
}

func (suite *AuthServiceUnitTestsSuite) TestAuthService_Register_UserGrpcErrorCreatingUser() {
	registrationRequest := registration_request.RegistrationRequest{
		User: user.User{
			Name:      "Ime",
			OwnerID:   1,
			OwnerType: "registered_users",
			Username:  "Username",
			FirstName: "FirstName",
			LastName:  "LastName",
			Phone:     "02123456678",
			BirthDate: 123,
			Website:   "www.site.com",
			Biography: "Biography",
			Gender:    1,
			Email:     "mail@mail.com",
		},
		Password: "123123123",
		Role:     "user",
	}
	errGrpc := errors.New("")
	err := rest_error.NewInternalServerError("user grpc client error when creating user", errGrpc)

	suite.userGrpcClientMock.On("CreateUser", registrationRequest).Return(new(uint), errGrpc).Once()

	_, regErr := suite.service.Register(registrationRequest)

	assert.Equal(suite.T(), err, regErr)
}

func (suite *AuthServiceUnitTestsSuite) TestAuthService_Register_UserGrpcErrorDeletingUser() {
	registrationRequest := registration_request.RegistrationRequest{
		User: user.User{
			Name:      "Ime",
			OwnerID:   1,
			OwnerType: "registered_users",
			Username:  "Username",
			FirstName: "FirstName",
			LastName:  "LastName",
			Phone:     "02123456678",
			BirthDate: 123,
			Website:   "www.site.com",
			Biography: "Biography",
			Gender:    1,
			Email:     "mail@mail.com",
		},
		Password: "123123123",
		Role:     "user",
	}
	errGrpc := errors.New("")
	err := rest_error.NewInternalServerError("user grpc client error when deleting user", errGrpc)

	suite.userGrpcClientMock.On("CreateUser", registrationRequest).Return(new(uint), nil).Once()
	suite.auth0GrpcClientMock.On("RegisterUserOnAuth0", registrationRequest.Email, registrationRequest.Password,
		registrationRequest.Role).Return("", errGrpc).Once()
	suite.userGrpcClientMock.On("DeleteUser", new(uint), registrationRequest.Role).Return(errGrpc).Once()

	_, regErr := suite.service.Register(registrationRequest)

	assert.Equal(suite.T(), err, regErr)
}

func (suite *AuthServiceUnitTestsSuite) TestAuthService_Register_AuthGrpcError() {
	registrationRequest := registration_request.RegistrationRequest{
		User: user.User{
			Name:      "Ime",
			OwnerID:   1,
			OwnerType: "registered_users",
			Username:  "Username",
			FirstName: "FirstName",
			LastName:  "LastName",
			Phone:     "02123456678",
			BirthDate: 123,
			Website:   "www.site.com",
			Biography: "Biography",
			Gender:    1,
			Email:     "mail@mail.com",
		},
		Password: "123123123",
		Role:     "user",
	}
	errGrpc := errors.New("")
	err := rest_error.NewInternalServerError("auth0 client error", nil)

	suite.userGrpcClientMock.On("CreateUser", registrationRequest).Return(new(uint), nil).Once()
	suite.auth0GrpcClientMock.On("RegisterUserOnAuth0", registrationRequest.Email, registrationRequest.Password,
		registrationRequest.Role).Return("", errGrpc).Once()
	suite.userGrpcClientMock.On("DeleteUser", new(uint), registrationRequest.Role).Return(nil).Once()

	_, regErr := suite.service.Register(registrationRequest)

	assert.Equal(suite.T(), err, regErr)
}

func (suite *AuthServiceUnitTestsSuite) TestAuthService_Register() {
	registrationRequest := registration_request.RegistrationRequest{
		User: user.User{
			Name:      "Ime",
			OwnerID:   1,
			OwnerType: "registered_users",
			Username:  "Username",
			FirstName: "FirstName",
			LastName:  "LastName",
			Phone:     "02123456678",
			BirthDate: 123,
			Website:   "www.site.com",
			Biography: "Biography",
			Gender:    1,
			Email:     "mail@mail.com",
		},
		Password: "123123123",
		Role:     "user",
	}

	suite.userGrpcClientMock.On("CreateUser", registrationRequest).Return(new(uint), nil).Once()
	suite.auth0GrpcClientMock.On("RegisterUserOnAuth0", registrationRequest.Email, registrationRequest.Password,
		registrationRequest.Role).Return("", nil).Once()

	id, regErr := suite.service.Register(registrationRequest)

	assert.Equal(suite.T(), uint(0), id)
	assert.Equal(suite.T(), nil, regErr)
}