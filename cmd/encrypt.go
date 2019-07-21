package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

	encrypt.SetUsageTemplate("Example Usage:\n" + "\tgo-video-cli encrypt [FILEPATH] [KEYNAME] [REGION] \n")
}

var encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file provided on the command line",
	Long:  "Encrypt a file provided on the command line by providing a path to the file, the name of the key stored in AWS Secrets Manager, and the AWS region where the key is stored",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {

		// Arguments
		filepath := args[0]
		keyName := args[1]
		region := args[2]

		// Create AWS Session
		log.Info("Creating AWS Session")
		session := createAWSSession(region)

		// Create Secrets Manager service
		svc := secretsmanager.New(session)

		// Build Get Secrets Value input
		input := buildGetSecretInput(keyName)

		// Get Secrets Value
		log.Info("Retrieving Secret Key")
		key := getSecret(svc, input)

		// Read in file
		log.Info("Reading input file")
		data, err := ioutil.ReadFile(filepath)
		check(err, "There was an issue reading the file")

		// Encrypting the data
		log.Info("Encrypting File")
		cipherText := createCipherText(&data, key)

		outputPath := fmt.Sprintf("%s.encrypted", filepath)

		// Write the data out to a file
		log.Info(fmt.Sprintf("Outputting file to your directory as %s", outputPath))
		err = ioutil.WriteFile(outputPath, cipherText, 0777)
		check(err, "There was an issue writing the encrypted file")

		log.Info("Success!")
	},
}

func createCipherText(data *[]byte, key string) []byte {
	byteKey := []byte(key)

	block, err := aes.NewCipher(byteKey)
	check(err, "There was an issue creating the cipher block")

	gcm, err := cipher.NewGCM(block)
	check(err, "There was an issue creating the GCM cipher")

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	check(err, "There was an issue reading the nonce")

	ciphertext := gcm.Seal(nonce, nonce, *data, nil)
	return ciphertext
}
