package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(encrypt)
}

var encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file provided on the command line",
	Long:  "Encrypt a file provided on the command line",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}
