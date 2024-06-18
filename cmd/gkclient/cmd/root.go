/*
Copyright © 2024 keycubasic keycube@yandex.ru
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const ParamErrStr = "insufficient input parameters. For details use --help"

var (
	ServerAddress string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gkclient",
		Short: "CLI for GophKeeper Password Manager",
		Long:  `GkClient is a CLI for for GophKeeper Password Manager.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ServerAddress, "address", "a", ":8080", "Server address")
	rootCmd.MarkPersistentFlagRequired("a")
}
