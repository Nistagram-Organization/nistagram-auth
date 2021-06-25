package auth

import (
	"github.com/Nistagram-Organization/agent-shared/src/utils/rest_error"
	auth02 "github.com/Nistagram-Organization/nistagram-auth/src/clients/auth0"
	"github.com/Nistagram-Organization/nistagram-auth/src/clients/auth_grpc_client"
	"github.com/Nistagram-Organization/nistagram-auth/src/dto/registration_request"
)

type AuthService interface {
	Register(request registration_request.RegistrationRequest) rest_error.RestErr
}

type authService struct {
	auth0Client    auth02.Auth0Client
	authGrpcClient auth_grpc_client.AuthGrpcClient
}

func NewAuthService(auth0Client auth02.Auth0Client, authGrpcClient auth_grpc_client.AuthGrpcClient) AuthService {
	return &authService{
		auth0Client,
		authGrpcClient,
	}
}

func (s *authService) Register(registrationRequest registration_request.RegistrationRequest) rest_error.RestErr {
	if err := registrationRequest.Validate(); err != nil {
		return err
	}

	if err := s.authGrpcClient.Register(registrationRequest); err != nil {
		return rest_error.NewInternalServerError("auth grpc client error", err)
	}

	if _, err := s.auth0Client.RegisterUserOnAuth0(registrationRequest.Email, registrationRequest.Password, registrationRequest.Role); err != nil {
		return rest_error.NewInternalServerError("auth0 client error", err)
	}

	return nil
}
