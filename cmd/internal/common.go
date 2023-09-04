package internal

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"
)

const (
	YesFlag = "yes"
	NoFlag  = "no"

	DevKeyFlag    = "dev"
	DevKeyEnvFlag = "ZK_DEV"

	VerboseKeyFlag    = "verbose"
	VerboseKeyEnvFlag = "VERBOSE"
)

var ErrExecutionAborted = errors.New("execution aborted")

func AddBoolFlag(cmd *cobra.Command, name, envFlag, shorthand string, defaultValue bool, usage string, isHidden bool) {
	err := addBoolFlagE(cmd, name, envFlag, shorthand, defaultValue, usage, isHidden)
	if err == nil {
		ui.GlobalWriter.Errorf("Unable to add flag %s", name)
	}
}

func addBoolFlagE(cmd *cobra.Command, name, envFlag, shorthand string, defaultValue bool, usage string, isHidden bool) error {
	cmd.PersistentFlags().BoolP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		return err
	}
	err = viper.BindEnv(name, envFlag)
	if err != nil {
		return err
	}

	if isHidden {
		err = cmd.PersistentFlags().MarkHidden(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddStringFlag(cmd *cobra.Command, name, envFlag, shorthand, defaultValue, usage string, isHidden bool) {
	err := addStringFlagE(cmd, name, envFlag, shorthand, defaultValue, usage, isHidden)
	if err != nil {
		ui.GlobalWriter.Errorf("Unable to add flag %s", name)
	}
}

func addStringFlagE(cmd *cobra.Command, name, envFlag, shorthand, defaultValue, usage string, isHidden bool) error {
	cmd.PersistentFlags().StringP(name, shorthand, defaultValue, usage)
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		return err
	}
	err = viper.BindEnv(name, envFlag)
	if err != nil {
		return err
	}

	if isHidden {
		err = cmd.PersistentFlags().MarkHidden(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetKVPairsFromCSV(csv string) map[string]string {
	pairs := strings.Split(csv, ",") // Split by comma
	keyValueMap := make(map[string]string)
	for _, pair := range pairs {
		parts := strings.Split(pair, "=") // Split by equal sign
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			keyValueMap[key] = value
		}
	}
	return keyValueMap
}

func ExecuteShellFile(shellFile, inputParameters, spinnerText, successMessage, errorMessage string) error {

	printError := viper.Get(VerboseKeyFlag) == true

	// make the file executable
	out, err := shell.ShelloutWithSpinner("chmod +x "+shellFile, "checking for appropriate permissions", "permissions added", "failed to get permissions")
	if err != nil {
		utils.DumpErrorAndPrintLocation(out, printError)
		return err
	}

	out, err = shell.ShelloutWithSpinner(shellFile+inputParameters, spinnerText, successMessage, errorMessage)
	if err != nil {
		utils.DumpErrorAndPrintLocation(out, printError)
		return err
	}
	return nil
}
