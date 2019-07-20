package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var author string

var rootCmd = &cobra.Command{
	Use:   "go-video-cli",
	Short: "A simple cli for encrypting files, decrypting files, uploading files to S3, and deleting files from S3",
	Long:  "A simple cli for encrypting files, decrypting files, uploading files to S3, and deleting files from S3",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&author, "author", "Nicholas McHenry", "Author name for copyright attribution")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}

// Execute : The main function to execute the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
