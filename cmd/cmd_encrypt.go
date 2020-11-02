package cmd

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/codingtony/ansible-vault-go/vault"
	"github.com/spf13/cobra"
)

type fileEncryptFlagsStruct struct {
	file     string
	password string
}

type randomTextEncryptFlagsStruct struct {
	password string
	length   int
}

var (

	// Runtime State
	fileEncryptFlags       = &fileEncryptFlagsStruct{}
	randomTextEncryptFlags = &randomTextEncryptFlagsStruct{}

	// Command
	fileEncryptCmd = &cobra.Command{
		Use:                   "encrypt [flags] [file]",
		Short:                 "Encrypt a file.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Errors at this point are no longer related to flags
			rootCmd.SilenceUsage = true
			fileEncryptFlags.password = RootPFlags.Password
			fileEncryptFlags.file = args[0]

			return doEncryptFile(fileEncryptFlags)
		},
	}
	randomTextEncryptCmd = &cobra.Command{
		Use:   "random_text_encrypt [flags] [vault_password]",
		Short: "Generate random text and encrypt it.",
		Long: ` Generate random text and encrypt it.
Random text is alphanumerical and of lenght controlled 
by the lenght parameter 
		`,

		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Errors at this point are no longer related to flags
			rootCmd.SilenceUsage = true
			randomTextEncryptFlags.password = RootPFlags.Password

			return doRandomTextEncrypt(randomTextEncryptFlags)
		},
	}
)

// //goland:noinspection GoUnhandledErrorResult
func init() {
	rootCmd.AddCommand(fileEncryptCmd)
	rootCmd.AddCommand(randomTextEncryptCmd)
	randomTextEncryptCmd.Flags().
		IntVarP(&randomTextEncryptFlags.length, "length", "l", 32, "length of generated random content")
}

func doEncryptFile(flags *fileEncryptFlagsStruct) error {
	data, err := ioutil.ReadFile(flags.file)
	if err != nil {
		return err
	}
	f, err := os.Create(flags.file)
	defer f.Close()
	if err != nil {
		return err
	}
	cipher, err := vault.EncryptByteArray(data, flags.password)
	if err != nil {
		return err
	}
	_, err = f.WriteString(cipher)
	if err != nil {
		return err
	}
	out.Infof("Encryption successful")
	return nil
}

func doRandomTextEncrypt(flags *randomTextEncryptFlagsStruct) error {
	data := randomString(flags.length)
	out.Debugf("Generated string : %s", data)
	cipher, err := vault.Encrypt(data, flags.password)
	if err != nil {
		return err
	}
	fmt.Println(cipher)
	return nil
}

// from https://www.calhoun.io/creating-random-strings-in-go/
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, charset)
}
