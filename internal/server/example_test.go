// Package server shows examples.
package server

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
)

// ExampleRegisterUser demonstrates calling of RegisterUser.
func ExampleRegisterUser() {
	ctx := context.Background()
	authClient, err := client.NewAuthClient(ctx, ":8080")
	if err != nil {
		log.Fatal(err)
	}

	user := &proto.User{UserName: "username", Password: "password"}
	res, err = authClient.Client.RegisterUser(ctx, &proto.RegisterUserRequest{User: user})
	if err != nil {
		log.Error(err)
	}

	log.Info(res.AccessToken)
}

// ExampleLoginUser demonstrates calling of LoginUser.
func ExampleLoginUser() {
	ctx := context.Background()
	authClient, err := client.NewAuthClient(ctx, suite.cfg.ServerAddress)
	if err != nil {
		log.Error(err)
	}

	user := &proto.User{UserName: "username", Password: "password"}
	res, err = authClient.Client.LoginUser(ctx, &proto.LoginUserRequest{User: user})
	if err != nil {
		log.Error(err)
	}

	log.Info(res.AccessToken)
}

// ExampleStoreData demonstrates calling of StoreData.
func ExampleStoreData() {
	ctx := context.Background()
	gkClient, err := client.NewGophKeeperClient(ctx, ":8080")
	if err != nil {
		log.Error(err)
	}

	authCtx := metadata.AppendToOutgoingContext(suite.ctx, "authorization", "Bearer <token>")

	creds := &proto.Credentials{
		Login:    "login",
		Password: "123456",
	}

	res, err := gkClient.Client.StoreData(authCtx, &proto.StoreDataRequest{Site: "somesite.com", Credentials: creds, Metadata: "somemeta"})
	if err != nil {
		log.Error(err)
	}

	log.Info(res.Info)
}

// ExampleGetData demonstrates calling of GetData.
func ExampleGetData() {
	ctx := context.Background()
	gkClient, err := client.NewGophKeeperClient(ctx, ":8080")
	if err != nil {
		log.Error(err)
	}

	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer <token>")
	res, err := gkClient.Client.GetData(authCtx, &proto.GetDataRequest{Site: "somesite.com"})
	if err != nil {
		log.Error(err)
	}

	log.Info(res.Password)
}
