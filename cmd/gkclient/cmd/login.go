/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   `login --u <username> --p <password>`,
	Short: "Login user",
	Long:  `Login user with username and password`,
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
		response, err := authClient.Client.LoginUser(ctx, &apb.LoginUserRequest{User: user})
		if err != nil {
			return err
		}

		fmt.Println(response.AccessToken)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().String("u", "", "Username")
	loginCmd.Flags().String("p", "", "User password")
	loginCmd.MarkFlagsRequiredTogether("u", "p")
}
