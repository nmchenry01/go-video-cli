package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

/*
	TODO Better documentation around the usage
	TODO Look in to more flags/validation options
	TODO Write tests
	TODO Benchmarking
*/

func init() {
	rootCmd.AddCommand(encrypt)

	encrypt.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli encrypt [FILEPATH] [PASSWORD]\n")
}

var encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file provided on the command line",
	Long:  "Encrypt a file provided on the command line by providing a path to the file and a password",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// Arguments
		filepath := args[0]
		password := args[1]

		// Read in file
		log.Info("Reading input file")
		data, err := ioutil.ReadFile(filepath)
		check(err, "There was an issue reading the file")

		// Encrypting the data
		log.Info("Encrypting File")
		cipherText := createCipherText(&data, password)

		outputPath := fmt.Sprintf("%s.encrypted", filepath)

		// Write the data out to a file
		log.Info(fmt.Sprintf("Outputting file to your directory as %s", outputPath))
		err = ioutil.WriteFile(outputPath, cipherText, 0777)
		check(err, "There was an issue writing the encrypted file")

		log.Info("Success!")
	},
}

func createCipherText(data *[]byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))

	block, err := aes.NewCipher(key)
	check(err, "There was an issue creating the cipher block")

	gcm, err := cipher.NewGCM(block)
	check(err, "There was an issue creating the GCM cipher")

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	check(err, "There was an issue reading the nonce")

	ciphertext := gcm.Seal(nonce, nonce, *data, nil)
	return ciphertext
}
