package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"zkctl/cmd/internal"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"k8s.io/client-go/util/homedir"
	// operator "github.com/zerok-ai/operator/cmd"
	// operatior "zkctl/cmd/operator/cmd"
)

const (
	zerok_dir_name = ".zerok"

	KUBECONFIG_FLAG   = "kubeconfig"
	KUBECONTEXT_FLAG  = "kube-context"
	CLUSTER_NAME_FLAG = "cluster-name"
	PRECISE_FLAG      = "precise"
)

var (
	RootCmd = &cobra.Command{
		Use:               "zkctl",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		SilenceUsage:      true,
		SilenceErrors:     true,
		Short:             "zkctl - a CLI to interact with ZeroK",
		// Version:           "0.0.1",
		Long: `
  
	❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄
	__  ___  __   __           ___                     
	 / |__  |__) /  \    |__/ |__  |    \  / | |\ |    
	/_ |___ |  \ \__/    |  \ |___ |___  \/  | | \| 

	                                ❄ As Kool as it G8s!
	❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄
	
	 	`,
	}
)

func ExecuteContext(ctx context.Context) error {
	err := RootCmd.ExecuteContext(ctx)

	if err == nil {
		// ui.GlobalWriter.PrintSuccessMessageln("Executed successfully")
		return nil
	}
	ui.GlobalWriter.PrintErrorMessageln(fmt.Sprintf("Whoops. There was an error in zkctl '%s'", err))
	os.Exit(1)
	return err
}

func init() {

	cobra.OnInitialize(initConfigFromFile)

	RootCmd.PersistentFlags().String(utils.ZEROK_DIR_PATH_FLAG, getZerokDirDefaultPath(), "default location of zerok's directory on dev machine")
	viper.BindPFlag(utils.ZEROK_DIR_PATH_FLAG, RootCmd.PersistentFlags().Lookup(utils.ZEROK_DIR_PATH_FLAG))

	RootCmd.PersistentFlags().String(CLUSTER_NAME_FLAG, "", "cluster name")
	viper.BindPFlag(CLUSTER_NAME_FLAG, RootCmd.PersistentFlags().Lookup(CLUSTER_NAME_FLAG))

	internal.AddBoolFlag(RootCmd, PRECISE_FLAG, "", "", false, "for internal use only", true)

	RootCmd.PersistentFlags().BoolP(internal.YesFlag, "y", false, "Automatic yes to prompts")
	viper.BindPFlag(internal.YesFlag, RootCmd.PersistentFlags().Lookup(internal.YesFlag))

	RootCmd.PersistentFlags().String(KUBECONTEXT_FLAG, "", "name of the kubeconfig context to use")
	viper.BindPFlag(KUBECONTEXT_FLAG, RootCmd.PersistentFlags().Lookup(KUBECONTEXT_FLAG))

	home := homedir.HomeDir()
	RootCmd.PersistentFlags().String(KUBECONFIG_FLAG, filepath.Join(home, ".kube", "config"), "path to the kubeconfig file")
	viper.BindPFlag(KUBECONFIG_FLAG, RootCmd.PersistentFlags().Lookup(KUBECONFIG_FLAG))
	viper.BindEnv(KUBECONFIG_FLAG)

	//internal flags
	internal.AddBoolFlag(RootCmd, internal.VerboseKeyFlag, internal.VerboseKeyEnvFlag, "v", false, "for internal use only", true)
	internal.AddBoolFlag(RootCmd, internal.DevKeyFlag, internal.DevKeyEnvFlag, "d", false, "for internal use only", true)

}

func initConfigFromFile() {

	viper.AutomaticEnv()

	// if cfgFile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cfgFile)
	// } else {
	// Find in current directory.
	// workingDir, err := os.Getwd()
	// cobra.CheckErr(err)
	// viper.AddConfigPath(workingDir)
	// viper.SetConfigType("env")
	// viper.SetConfigName("zkctl")
	// }
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	panic(err)
	// }

	// viper.Set("PLC_CLUSTER", fmt.Sprintf("gke_zerok-dev_us-west1-b_%s", plc_cluster))
}

func getZerokDirDefaultPath() string {
	var err error

	var baseDir string
	if baseDir, err = os.UserHomeDir(); err != nil {
		baseDir = os.TempDir()
	}

	return fmt.Sprintf("%s%c%s", baseDir, os.PathSeparator, zerok_dir_name)
}
