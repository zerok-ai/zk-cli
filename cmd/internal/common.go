package internal

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"
)

const (
	YesFlag = "yes"
	NoFlag  = "no"

	DevKeyFlag    = "dev"
	DevKeyEnvFlag = "ZK_DEV"

	GptKeyFlag    = "gpt"
	GptKeyEnvFlag = "ZK_GPT"

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

func DumpErrorAndPrintLocation(errorMessage string) {
	printOnConsole := viper.Get(VerboseKeyFlag) == true
	utils.DumpErrorAndPrintLocation(errorMessage, printOnConsole)
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
