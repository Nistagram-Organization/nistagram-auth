package user_grpc_client

import (
	"github.com/Nistagram-Organization/nistagram-auth/src/dto/registration_request"
	"github.com/stretchr/testify/mock"
)

type UserGrpcClientMock struct {
	mock.Mock
}

func (u *UserGrpcClientMock) CreateUser(request registration_request.RegistrationRequest) (*uint, error) {
	args := u.Called(request)
	if args.Get(1) == nil {
		return args.Get(0).(*uint), nil
	}
	return nil, args.Get(1).(error)
}

func (u *UserGrpcClientMock) DeleteUser(u2 *uint, s string) error {
	args := u.Called(u2, s)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}




