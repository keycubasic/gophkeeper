package server

import (
	"context"
	"fmt"

	"github.com/keycubasic/gophkeeper/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// ExampleStoreData demonstrates calling of StoreData.
func (s *GrpcServer) ExampleStoreData() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, ":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}

	client := proto.NewGophKeeperServiceClient(conn)
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer <token>")

	creds := &proto.Credentials{
		Login:    "login",
		Password: "123456",
	}

	res, err := client.StoreData(authCtx, &proto.StoreDataRequest{Site: "somesite.com", Credentials: creds, Metadata: "somemeta"})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res.Info)
}

// ExampleGetData demonstrates calling of GetData.
func (s *GrpcServer) ExampleGetData() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, ":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err.Error())
	}

	client := proto.NewGophKeeperServiceClient(conn)
	authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer <token>")
	res, err := client.GetData(authCtx, &proto.GetDataRequest{Site: "somesite.com"})
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res.Credentials.Password)
}
