package auth

import (
	auth02 "github.com/Nistagram-Organization/nistagram-auth/src/clients/auth0"
	"github.com/Nistagram-Organization/nistagram-auth/src/clients/user_grpc_client"
	"github.com/Nistagram-Organization/nistagram-auth/src/dto/registration_request"
	"github.com/Nistagram-Organization/nistagram-shared/src/utils/rest_error"
)

type AuthService interface {
	Register(request registration_request.RegistrationRequest) rest_error.RestErr
}

type authService struct {
	auth0Client    auth02.Auth0Client
	userGrpcClient user_grpc_client.UserGrpcClient
}

func NewAuthService(auth0Client auth02.Auth0Client, userGrpcClient user_grpc_client.UserGrpcClient) AuthService {
	return &authService{
		auth0Client,
		userGrpcClient,
	}
}

func (s *authService) Register(registrationRequest registration_request.RegistrationRequest) rest_error.RestErr {
	if err := registrationRequest.Validate(); err != nil {
		return err
	}

	var id *uint
	var err error

	if id, err = s.userGrpcClient.CreateUser(registrationRequest); err != nil {
		return rest_error.NewInternalServerError("user grpc client error when creating user", err)
	}

	if _, err := s.auth0Client.RegisterUserOnAuth0(registrationRequest.Email, registrationRequest.Password, registrationRequest.Role); err != nil {
		if err = s.userGrpcClient.DeleteUser(id, registrationRequest.Role); err != nil {
			return rest_error.NewInternalServerError("user grpc client error when deleting user", err)
		}
		return rest_error.NewInternalServerError("auth0 client error", err)
	}

	return nil
}
