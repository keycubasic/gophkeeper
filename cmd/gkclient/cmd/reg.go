/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
)

// regCmd represents the reg command
var regCmd = &cobra.Command{
	Use:   `reg --u "username" --p "password"`,
	Short: "Register new user",
	Long:  `Registers new user with username and password`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		username, err := cmd.Flags().GetString("u")
		if err != nil {
			log.Error().Err(err)
			return
		}

		password, err := cmd.Flags().GetString("p")
		if err != nil {
			log.Error().Err(err)
			return
		}

		if username == "" || password == "" {
			log.Error().Msg(ParamErrStr)
			return
		}

		authClient := client.NewAuthClient(ctx, ServerAddress)
		user := &apb.User{UserName: username, Password: password}
		response, err := authClient.Client.RegisterUser(ctx, &apb.RegisterUserRequest{User: user})
		if err != nil {
			log.Error().Err(err)
			return
		}

		log.Info().Msg(response.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(regCmd)

	regCmd.PersistentFlags().String("u", "", "Username")
	regCmd.PersistentFlags().String("p", "", "User password")
}
