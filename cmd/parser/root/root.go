package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "An example app to parse currencies",
	Long:  `This application demonstrates how to use Cobra to parse a currency amount.`,
	Run: func(cmd *cobra.Command, args []string) {
		// default behavior when no subcommand is provided
		fmt.Println("Use the 'parse' subcommand.")
	},
}

// Execute runs the root command.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var ParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse various types (currency, gas, store,…)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify what to parse. Available subcommands: currency, gas, store")
	},
}

func init() {
	RootCmd.AddCommand(ParseCmd)
}
