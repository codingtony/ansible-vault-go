package cmd

import (
	"os"

	"github.com/codingtony/ansible-vault-go/vault"
	"github.com/spf13/cobra"
)

type fileDecryptFlagsStruct struct {
	file     string
	password string
}

var (

	// Runtime State
	fileDecryptFlags = &fileDecryptFlagsStruct{}

	// Command
	fileDecryptCmd = &cobra.Command{
		Use:                   "decrypt [flags] [file]",
		Short:                 "Decrypt a file.",
		DisableFlagsInUseLine: true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Errors at this point are no longer related to flags
			rootCmd.SilenceUsage = true
			fileDecryptFlags.file = args[0]
			fileDecryptFlags.password = RootPFlags.Password

			return doDecryptFile(fileDecryptFlags)
		},
	}
)

// //goland:noinspection GoUnhandledErrorResult
func init() {
	rootCmd.AddCommand(fileDecryptCmd)
}

func doDecryptFile(flags *fileDecryptFlagsStruct) error {
	plaintext, err := vault.DecryptFile(flags.file, flags.password)
	if err != nil {
		return err
	}

	f, err := os.Create(flags.file)
	defer f.Close()
	if err != nil {
		return err
	}

	_, err = f.WriteString(plaintext)
	if err != nil {
		return err
	}
	out.Infof("Decryption successful")
	return nil
}
