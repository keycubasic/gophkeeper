package server

import (
	"context"
	"fmt"

	"github.com/keycubasic/gophkeeper/internal/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ExampleRegisterUser demonstrates calling of RegisterUser.
func (server *AuthServer) ExampleRegisterUser() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, ":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}

	client := proto.NewAuthServiceClient(conn)

	user := &proto.User{UserName: "username", Password: "password"}
	res, err := client.RegisterUser(ctx, &proto.RegisterUserRequest{User: user})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res.AccessToken)
}

// ExampleLoginUser demonstrates calling of LoginUser.
func (server *AuthServer) ExampleLoginUser() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, ":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}

	client := proto.NewAuthServiceClient(conn)

	user := &proto.User{UserName: "username", Password: "password"}
	res, err := client.LoginUser(ctx, &proto.LoginUserRequest{User: user})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res.AccessToken)
}
