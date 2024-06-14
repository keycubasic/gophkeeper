package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
)

// TODO: add cobra CLI - https://dev.to/mauriciolinhares/building-and-distributing-a-command-line-tool-in-golang-go0

const (
	Reg         = "reg"
	Login       = "login"
	Post        = "post"
	Get         = "get"
	paramErrStr = "Insufficient input parameters. For details use --help"
)

type parameters struct {
	ServerAddress string
	Username      string
	Password      string
	Site          string
	Login         string
	Metadata      string
}

func main() {
	ctx := context.Background()
	method := os.Args[1]
	p := new(parameters)
	parseParameters(p)

	conn, err := grpc.DialContext(ctx, p.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()

	log.Info().Msgf("Server %s connected", p.ServerAddress)
	authClient := apb.NewAuthServiceClient(conn)
	dataClient := gpb.NewGophKeeperServiceClient(conn)
	Processing(ctx, authClient, dataClient, method, *p)
}

func parseParameters(p *parameters) {
	flag.StringVar(&p.ServerAddress, "a", ":8080", "Server address")
	flag.StringVar(&p.Username, "u", "", "Username")
	flag.StringVar(&p.Password, "p", "", "User password/Site password")
	flag.StringVar(&p.Site, "s", "", "Site URL")
	flag.StringVar(&p.Login, "l", "", "Site login")
	flag.StringVar(&p.Metadata, "m", "", "Additional metadata")
	flag.Parse()
	params, err := json.Marshal(p)
	if err != nil {
		log.Error().Err(err)
		return
	}
	log.Info().Msgf("Input params:\n%s", params)
}

func Processing(ctx context.Context, authClient apb.AuthServiceClient, dataClient gpb.GophKeeperServiceClient, method string, params parameters) {
	switch method {
	case Reg: // $ cmd/gkclient/gkclient reg -u "username" -p "password"
		{
			if params.Username == "" || params.Password == "" {
				log.Error().Msg(paramErrStr)
				return
			}

			user := &apb.User{UserName: params.Username, Password: params.Password}
			response, err := authClient.RegisterUser(ctx, &apb.RegisterUserRequest{User: user})
			if err != nil {
				log.Error().Err(err)
				return
			}

			log.Info().Msg(response.AccessToken)
		}
	case Login: // $ cmd/gkclient/gkclient login -u "username" -p "password"
		{
			if params.Username == "" || params.Password == "" {
				log.Error().Msg(paramErrStr)
				return
			}

			user := &apb.User{UserName: params.Username, Password: params.Password}
			response, err := authClient.LoginUser(ctx, &apb.LoginUserRequest{User: user})
			if err != nil {
				log.Error().Err(err)
				return
			}

			log.Info().Msg(response.AccessToken)
		}
	case Post: // $ cmd/gkclient/gkclient post -s "https://ya.ru" -l "login" -p "password" -m "metadata"
		{
			if params.Site == "" || params.Login == "" || params.Password == "" {
				log.Error().Msg(paramErrStr)
				return
			}

			creds := &gpb.Credentials{Login: params.Login, Password: params.Password}
			request := &gpb.StoreDataRequest{
				Site:        params.Site,
				Credentials: creds,
				Metadata:    params.Metadata,
			}

			response, err := dataClient.StoreData(ctx, request)
			if err != nil {
				log.Error().Err(err)
				return
			}

			log.Info().Msg(response.Info)
		}
	case Get: // $ cmd/gkclient/gkclient get -s "https://ya.ru"
		{
			if params.Site == "" {
				log.Error().Msg(paramErrStr)
				return
			}

			request := &gpb.GetDataRequest{Site: params.Site}
			response, err := dataClient.GetData(ctx, request)
			if err != nil {
				log.Error().Err(err)
				return
			}

			res, err := json.Marshal(response)
			if err != nil {
				log.Error().Err(err)
				return
			}

			log.Info().Msg(string(res))
		}
	}
}
