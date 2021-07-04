package auth_grpc_service

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-auth/src/clients/auth0"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
)

type authGrpcService struct {
	proto.AuthServiceServer
	auth0Client auth0.Auth0Client
}

func NewAuthGrpcService(auth0Client auth0.Auth0Client) proto.AuthServiceServer {
	return &authGrpcService{
		proto.UnimplementedAuthServiceServer{},
		auth0Client,
	}
}

func (s *authGrpcService) TerminateProfile(ctx context.Context, request *proto.TerminateProfileRequest) (*proto.TerminateProfileResponse, error) {
	email := request.Email
	if err := s.auth0Client.BlockUserOnAuth0(email); err != nil {
		return nil, err
	}
	response := proto.TerminateProfileResponse{Success: true}

	return &response, nil
}
