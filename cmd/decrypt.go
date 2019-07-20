package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(decrypt)
}

var decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file provided on the command line",
	Long:  "Decrypt a file provided on the command line",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}

func createPlaintext(data *[]byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("There was an issue creating the cipher block")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("There was an issue creating the GCM block")
	}

	encryptedData := *data
	nonceSize := gcm.NonceSize()

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("There was an issue creating the decrypting the ciphertext")
	}

	return plaintext
}
