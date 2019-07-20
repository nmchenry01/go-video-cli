package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(encrypt)
}

var encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file provided on the command line",
	Long:  "Encrypt a file provided on the command line",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		// Arguments
		filepath := args[0]
		password := args[1]

		// Read in file
		data, readErr := ioutil.ReadFile(filepath)
		if readErr != nil {
			panic(readErr.Error())
		}

		// Hash the password to ensure 32 bytes
		hashedPassword := createHash(password)

		// Encrypt the data
		cipherText := createCipherText(data, hashedPassword)

		// Write the data out to a file
		writeErr := ioutil.WriteFile("encrypted", cipherText, 0777)
		if writeErr != nil {
			panic(writeErr.Error())
		}
	},
}

func createHash(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func createCipherText(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}
