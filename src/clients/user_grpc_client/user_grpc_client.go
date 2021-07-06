package user_grpc_client

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-auth/src/dto/registration_request"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"google.golang.org/grpc"
)

const (
	AGENT = "agent"
	USER  = "user"
)

type UserGrpcClient interface {
	CreateUser(registration_request.RegistrationRequest) (*uint, error)
	DeleteUser(*uint, string) error
}

type userGrpcClient struct {
	address string
}

func NewUserGrpcClient(docker bool) UserGrpcClient {
	var address string
	if docker {
		address = "nistagram-users:8084"
	} else {
		address = "127.0.0.1:8084"
	}
	return &userGrpcClient{
		address: address,
	}
}

func (c *userGrpcClient) CreateUser(user registration_request.RegistrationRequest) (*uint, error) {
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewUserServiceClient(conn)

	r, err := client.CreateUser(ctx,
		&proto.RegistrationRequest{
			Registration: user.ToUserMessage(),
		},
	)

	if err != nil {
		return nil, err
	}

	var id *uint
	id = new(uint)
	*id = uint(r.Id)

	return id, nil
}

func (c *userGrpcClient) DeleteUser(id *uint, role string) error {
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewUserServiceClient(conn)

	_, err = client.DeleteUser(ctx,
		&proto.DeleteUserRequest{
			Id:   uint64(*id),
			Role: getRole(role),
		},
	)
	return err
}

func getRole(role string) proto.Role {
	if role == AGENT {
		return proto.Role_AGENT
	} else if role == USER {
		return proto.Role_USER
	} else {
		return proto.Role_UNKNOWN
	}
}
