/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
)

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   `post --url "https://ya.ru" --login "login" --pass "password" --meta "metadata"`,
	Short: "Post user site credentials",
	Long:  `Post user site credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// TODO: add auth token to context metadata
		site, err := cmd.Flags().GetString("url")
		if err != nil {
			log.Error().Err(err)
			return
		}

		login, err := cmd.Flags().GetString("login")
		if err != nil {
			log.Error().Err(err)
			return
		}

		password, err := cmd.Flags().GetString("pass")
		if err != nil {
			log.Error().Err(err)
			return
		}

		if site == "" || login == "" || password == "" {
			log.Error().Msg(ParamErrStr)
			return
		}

		metadata, err := cmd.Flags().GetString("meta")
		if err != nil {
			log.Error().Err(err)
			return
		}

		gkClient := client.NewGophKeeperClient(ctx, ServerAddress)
		creds := &gpb.Credentials{Login: login, Password: password}
		request := &gpb.StoreDataRequest{
			Site:        site,
			Credentials: creds,
			Metadata:    metadata,
		}

		response, err := gkClient.Client.StoreData(ctx, request)
		if err != nil {
			log.Error().Err(err)
			return
		}

		log.Info().Msg(response.Info)
	},
}

func init() {
	rootCmd.AddCommand(postCmd)

	// TODO: add auth token
	postCmd.PersistentFlags().String("url", "", "Site URL")
	postCmd.PersistentFlags().String("login", "", "Site login")
	postCmd.PersistentFlags().String("pass", "", "Site password")
	postCmd.PersistentFlags().String("meta", "", "Additional metadata")
}
