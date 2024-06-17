/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   `login --u "username" --p "password"`,
	Short: "Login user",
	Long:  `Login user with username and password`,
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
		response, err := authClient.Client.LoginUser(ctx, &apb.LoginUserRequest{User: user})
		if err != nil {
			log.Error().Err(err)
			return
		}

		log.Info().Msg(response.AccessToken)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.PersistentFlags().String("u", "", "Username")
	loginCmd.PersistentFlags().String("p", "", "User password")
}
