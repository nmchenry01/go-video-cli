package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var author string

var rootCmd = &cobra.Command{
	Use:   "go-video-cli",
	Short: "A simple cli for encrypting files, decrypting files, uploading files to S3, and deleting files from S3",
	Long:  "A simple cli for encrypting files, decrypting files, uploading files to S3, and deleting files from S3",
}

// Execute : The main function to execute the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check(err error, message string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal(message)
	}
}

func createHash(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
