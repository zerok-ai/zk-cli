package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
	"k8s.io/client-go/util/homedir"
	// operator "github.com/zerok-ai/operator/cmd"
	// operatior "github.com/zerok-ai/zk-cli/zkctl/cmd/operator/cmd"
)

const (
	NAMESPACE_FLAG    = "namespace"
	KUBECONFIG_FLAG   = "kubeconfig"
	KUBECONTEXT_FLAG  = "kube-context"
	CLUSTER_NAME_FLAG = "cluster-name"
)

var (
	rootCmd = &cobra.Command{
		Use:           "zkctl",
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "zkctl - a CLI to interact with ZeroK",
		Version:       "0.0.1",
		Long: `
  
	❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄
	__  ___  __   __           ___                     
	 / |__  |__) /  \    |__/ |__  |    \  / | |\ |    
	/_ |___ |  \ \__/    |  \ |___ |___  \/  | | \| 

	                                ❄ As Kool as it G8s!
	❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄
	
	 	`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(os.Stderr, "Nothing to run here. Use --help for details \n")
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// this code is inherited by all the commands and run

			// #### Add session logic here ##################

			return nil
		},
	}

	// private variables
	cfgFile string
)

func ExecuteContext(ctx context.Context) error {
	ui.GlobalWriter.Println("")
	err := rootCmd.ExecuteContext(ctx)

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

	// path of the config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ${PWD}/zkctl.yaml)")

	// cloud provider
	rootCmd.PersistentFlags().StringP("cloudprovider", "P", "", "Name of a cloud provider. Allowed values gke|eks|minikube")
	viper.BindPFlag("cloudprovider", rootCmd.PersistentFlags().Lookup("cloudprovider"))

	rootCmd.PersistentFlags().String(CLUSTER_NAME_FLAG, "", "cluster name")
	viper.BindPFlag(CLUSTER_NAME_FLAG, rootCmd.PersistentFlags().Lookup(CLUSTER_NAME_FLAG))

	rootCmd.PersistentFlags().String(KUBECONTEXT_FLAG, "", "name of the kubeconfig context to use")
	viper.BindPFlag(KUBECONTEXT_FLAG, rootCmd.PersistentFlags().Lookup(KUBECONTEXT_FLAG))

	home := homedir.HomeDir()
	rootCmd.PersistentFlags().String(KUBECONFIG_FLAG, filepath.Join(home, ".kube", "config"), "path to the kubeconfig file")
	viper.BindPFlag(KUBECONFIG_FLAG, rootCmd.PersistentFlags().Lookup(KUBECONFIG_FLAG))
	viper.BindEnv(KUBECONFIG_FLAG)

}

func initConfigFromFile() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// // Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)
		// viper.AddConfigPath(home)

		// Find in current directory.
		workingDir, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(workingDir)

		// Search config in home directory with name "zkctl" (without extension).
		viper.SetConfigType("env")
		viper.SetConfigName("zkctl")
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
