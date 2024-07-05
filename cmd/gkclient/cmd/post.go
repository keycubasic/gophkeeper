/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
)

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   `post --s "https://ya.ru" --l "login" --p "password" --m "metadata"`,
	Short: "Post user site credentials",
	Long:  `Post user site credentials.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		token, err := cmd.Flags().GetString("t")
		if err != nil {
			return err
		}

		newCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

		site, err := cmd.Flags().GetString("s")
		if err != nil {
			return err
		}

		login, err := cmd.Flags().GetString("l")
		if err != nil {
			return err
		}

		password, err := cmd.Flags().GetString("p")
		if err != nil {
			return err
		}

		if site == "" || login == "" || password == "" {
			return errors.New(ParamErrStr)
		}

		metadata, err := cmd.Flags().GetString("m")
		if err != nil {
			return err
		}

		gkClient := client.NewGophKeeperClient(newCtx, ServerAddress, EnableTLS)
		defer gkClient.Conn.Close()

		creds := &gpb.Credentials{Login: login, Password: password}
		request := &gpb.StoreDataRequest{
			Site:        site,
			Credentials: creds,
			Metadata:    metadata,
		}

		response, err := gkClient.Client.StoreData(newCtx, request)
		if err != nil {
			return err
		}

		fmt.Println(response.Info)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(postCmd)

	postCmd.Flags().String("t", "", "Access token")
	postCmd.Flags().String("s", "", "Site URL")
	postCmd.Flags().String("l", "", "Site login")
	postCmd.Flags().String("p", "", "Site password")
	postCmd.MarkFlagsRequiredTogether("t", "s", "l", "p")
	postCmd.Flags().String("m", "", "Additional metadata")
}
