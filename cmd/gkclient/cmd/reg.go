/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
)

// regCmd represents the reg command
var regCmd = &cobra.Command{
	Use:   `reg --u <username> --p <password>`,
	Short: "Register new user",
	Long:  `Registers new user with username and password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		username, err := cmd.Flags().GetString("u")
		if err != nil {
			return err
		}

		password, err := cmd.Flags().GetString("p")
		if err != nil {
			return err
		}

		if username == "" || password == "" {
			return errors.New(ParamErrStr)
		}

		authClient := client.NewAuthClient(ctx, ServerAddress, EnableTLS)
		defer authClient.Conn.Close()

		user := &apb.User{UserName: username, Password: password}
		response, err := authClient.Client.RegisterUser(ctx, &apb.RegisterUserRequest{User: user})
		if err != nil {
			return err
		}

		fmt.Println(response.AccessToken)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(regCmd)

	regCmd.Flags().String("u", "", "Username")
	regCmd.Flags().String("p", "", "User password")
	regCmd.MarkFlagsRequiredTogether("u", "p")
}
