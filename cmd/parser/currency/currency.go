package currency

import (
	"fmt"

	"github.com/dealense7/go-rate-app/cmd/parser/root"
	"github.com/spf13/cobra"
)

func init() {
	root.ParseCmd.AddCommand(currencyCmd)
}

var currencyCmd = &cobra.Command{
	Use:   "currency",
	Short: "Parse a currency amount",
	Long:  `Parses the given currency amount (e.g. "USD 123.45").`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("parsing currency data")
		return nil
	},
}
