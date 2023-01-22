package zkctl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "zkctl",
		Short:   "zkctl - a CLI to interact with ZeroK",
		Version: "0.0.1",
		Long: `zkctl is the interface between a ZeroK user ZeroK
	
	It helps a user install and manage all the components of ZeroK`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(os.Stderr, "Nothing to run here. Use --help for details \n")
		},
	}

	cfgFile string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error in zkctl '%s'", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfigFromFile)

	// path of the config file
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ${PWD}/zkctl.yaml)")

	// cloud provider
	rootCmd.PersistentFlags().StringP("cloudprovider", "P", "", "Name of a cloud provider. Allowed values gke|eks|minikube")
	viper.BindPFlag("cloudprovider", rootCmd.PersistentFlags().Lookup("cloudprovider"))
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
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
