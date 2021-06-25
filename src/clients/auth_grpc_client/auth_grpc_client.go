package auth_grpc_client

import (
	"context"
	"github.com/Nistagram-Organization/nistagram-auth/src/dto/registration_request"
	"github.com/Nistagram-Organization/nistagram-shared/src/proto"
	"google.golang.org/grpc"
)

type AuthGrpcClient interface {
	Register(registration_request.RegistrationRequest) error
}

type authGrpcClient struct {
}

func NewAuthGrpcClient() AuthGrpcClient {
	return &authGrpcClient{}
}

func (c *authGrpcClient) Register(user registration_request.RegistrationRequest) error {
	conn, err := grpc.Dial("127.0.0.1:8084", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()


	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := proto.NewAuthServiceClient(conn)

	r, err := client.Register(ctx,
		&proto.RegistrationRequest{
			Registration: user.ToUserMessage(),
		},
	)

	if err != nil || !r.Success {
		return err
	}
	return nil
}
