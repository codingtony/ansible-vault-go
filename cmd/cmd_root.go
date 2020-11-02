package cmd

import (
	"fmt"

	"github.com/juju/loggo"
	"github.com/spf13/cobra"
)

type rootPFlagsStruct struct {
	Verbose bool
}

var (
	// This two variables are set by the Makefile
	version   string
	buildTime string

	// Logger
	out = loggo.GetLogger("cmd")

	// Runtime State
	rootPFlags = &rootPFlagsStruct{}

	// Command
	rootCmd = &cobra.Command{
		Use:           "ansible-vault-go",
		Short:         "Golang port of ansible-vault that can perform basic functions",
		Version:       fmt.Sprintf("%s (Built on: %s)", version, buildTime),
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			//goland:noinspection GoUnhandledErrorResult
			rootLogger := loggo.GetLogger("")

			if rootPFlags.Verbose {
				rootLogger.SetLogLevel(loggo.DEBUG)
			} else {
				rootLogger.SetLogLevel(loggo.INFO)
			}

			topLevelCmd := cmd
			for {
				if !topLevelCmd.HasParent() {
					break
				}

				topLevelCmd = topLevelCmd.Parent()
			}

			versionString := version

			out.Debugf("%s", versionString)

			return nil
		},
	}
)

//goland:noinspection GoUnhandledErrorResult
func init() {
	rootCmd.PersistentFlags().
		BoolVarP(&rootPFlags.Verbose, "verbose", "v", false, "enable verbose output")
}

func Execute() {
	//goland:noinspection GoUnhandledErrorResult

	err := rootCmd.Execute()

	if err != nil {
		if rootPFlags.Verbose {
			out.Errorf("%v", err)
		} else {
			out.Errorf("%s", err.Error())
		}
	}
}
