package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"strings"

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
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Arguments
		filepath := args[0]
		password := args[1]

		// Read in file
		log.Info("Reading input file")
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("There was an issue reading the file")
		}

		// Decrypting the data
		log.Info("Decrypting File")
		plaintext := createPlaintext(&data, password)

		basePath := strings.Split(filepath, ".encrypted")[0]

		// Write the data out to a file
		log.Info(fmt.Sprintf("Outputting file to your directory as %s", basePath))
		err = ioutil.WriteFile(basePath, plaintext, 0777)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("There was an issue writing the decrypted file")
		}

		log.Info("Success!")
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
