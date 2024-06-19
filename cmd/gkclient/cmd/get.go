/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   `get --s "https://ya.ru"`,
	Short: "Get user site credentials",
	Long:  `Get user site credentials.`,
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

		if site == "" {
			return errors.New(ParamErrStr)
		}

		gkClient := client.NewGophKeeperClient(newCtx, ServerAddress, EnableTLS)
		defer gkClient.Conn.Close()

		request := &gpb.GetDataRequest{Site: site}
		response, err := gkClient.Client.GetData(newCtx, request)
		if err != nil {
			return err
		}

		res, err := json.Marshal(response)
		if err != nil {
			return err
		}

		fmt.Println(string(res))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().String("t", "", "Access token")
	getCmd.Flags().String("s", "", "Site URL")
	getCmd.MarkFlagsRequiredTogether("t", "s")
}
