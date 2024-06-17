/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"context"
	"encoding/json"

	"github.com/keycubasic/gophkeeper/cmd/gkclient/internal/client"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   `get --s "https://ya.ru"`,
	Short: "Get user site credentials",
	Long:  `Get user site credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		// TODO: add auth token to context metadata
		site, err := cmd.Flags().GetString("s")
		if err != nil {
			log.Error().Err(err)
			return
		}

		if site == "" {
			log.Error().Msg(ParamErrStr)
			return
		}

		gkClient := client.NewGophKeeperClient(ctx, ServerAddress)
		request := &gpb.GetDataRequest{Site: site}
		response, err := gkClient.Client.GetData(ctx, request)
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
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// TODO: add auth token
	postCmd.PersistentFlags().String("s", "", "Site URL")
}
