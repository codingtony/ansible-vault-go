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
	File     string
	Password string
}

type randomTextEncryptFlagsStruct struct {
	Password string
	Length   int
}

var (

	// Runtime State
	fileEncryptFlags       = &fileEncryptFlagsStruct{}
	randomTextEncryptFlags = &randomTextEncryptFlagsStruct{}

	encryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Perfom ansible-vault encryption functions.",
	}

	// Command
	fileEncryptCmd = &cobra.Command{
		Use:                   "file [flags] [vault_password] [file]",
		Short:                 "Encrypt a file.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Errors at this point are no longer related to flags
			rootCmd.SilenceUsage = true
			fileEncryptFlags.Password = args[0]
			fileEncryptFlags.File = args[1]

			return doEncryptFile(fileEncryptFlags)
		},
	}
	randomTextEncryptCmd = &cobra.Command{
		Use:   "randomText [flags] [vault_password]",
		Short: "Generate random text and encrypt it.",
		Long: ` Generate random text and encrypt it.
Random text is alphanumerical and of lenght controlled 
by the lenght parameter 
		`,

		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Errors at this point are no longer related to flags
			rootCmd.SilenceUsage = true
			randomTextEncryptFlags.Password = args[0]

			return doRandomTextEncrypt(randomTextEncryptFlags)
		},
	}
)

// //goland:noinspection GoUnhandledErrorResult
func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.AddCommand(fileEncryptCmd)
	encryptCmd.AddCommand(randomTextEncryptCmd)
	randomTextEncryptCmd.Flags().
		IntVarP(&randomTextEncryptFlags.Length, "length", "l", 32, "length of generated random context.")
}

func doEncryptFile(flags *fileEncryptFlagsStruct) error {
	data, err := ioutil.ReadFile(flags.File)
	if err != nil {
		return err
	}
	f, err := os.Create(flags.File)
	defer f.Close()
	if err != nil {
		return err
	}
	cipher, err := vault.EncryptByteArray(data, flags.Password)
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
	data := randomString(flags.Length)
	out.Debugf("Generated string : %s", data)
	cipher, err := vault.Encrypt(data, flags.Password)
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
