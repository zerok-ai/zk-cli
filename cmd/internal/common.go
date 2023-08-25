package internal

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
)

const (
	YesFlag = "yes"
	NoFlag  = "no"
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

func ExecuteShellFile(shellFile, inputParameters, successMessage, errorMessage string) error {

	// make the file executable
	_, chmodErr := shell.ExecWithDurationAndSuccessM("chmod +x "+shellFile, "")
	if chmodErr != nil {
		ui.LogAndPrintError(fmt.Errorf(errorMessage+": %v", chmodErr))
		return chmodErr
	}

	// execute the file
	_, err := shell.ExecWithDurationAndSuccessM(shellFile+inputParameters, successMessage)
	if err != nil {
		ui.LogAndPrintError(fmt.Errorf(errorMessage+": %v", err))
	}
	return err
}
