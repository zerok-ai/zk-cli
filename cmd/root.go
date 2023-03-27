package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"k8s.io/client-go/util/homedir"
	// operator "github.com/zerok-ai/operator/cmd"
	// operatior "zkctl/cmd/operator/cmd"
)

const (
	zerok_dir_name = ".zerok"

	NAMESPACE_FLAG    = "namespace"
	KUBECONFIG_FLAG   = "kubeconfig"
	KUBECONTEXT_FLAG  = "kube-context"
	CLUSTER_NAME_FLAG = "cluster-name"
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
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Fprintf(os.Stderr, "Nothing to run here. Use --help for details \n")
		// },
		// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 	var err error
	
		// 	segment.SetScope(cmd.Name())
		// 	sentry_utils.SetTransactionOnCurrentScope(cmd.Name())
	
		// 	event := segment.NewEvent(cmd.Name())
		// 	defer event.Start()
	
		// 	if err = validateAuthentication(cmd, args); err != nil {
		// 		return err
		// 	}
	
		// 	if !viper.GetBool(SKIP_CLI_UPDATE_FLAG) {
		// 		return checkAndUpgradeVersion(cmd.Context())
		// 	}
	
		// 	return nil
		// },
	}

	ErrExecutionAborted = errors.New("execution aborted")

	ErrForceAborted = errors.New("force abort")
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

	// cloud provider
	RootCmd.PersistentFlags().String("cloudprovider", "P", "Name of a cloud provider. Allowed values gke|eks|minikube")
	viper.BindPFlag("cloudprovider", RootCmd.PersistentFlags().Lookup("cloudprovider"))

	RootCmd.PersistentFlags().String(utils.ZEROK_DIR_PATH_FLAG, getZerokDirDefaultPath(), "default location of zerok's directory on dev machine")
	viper.BindPFlag(utils.ZEROK_DIR_PATH_FLAG, RootCmd.PersistentFlags().Lookup(utils.ZEROK_DIR_PATH_FLAG))

	RootCmd.PersistentFlags().String(CLUSTER_NAME_FLAG, "", "cluster name")
	viper.BindPFlag(CLUSTER_NAME_FLAG, RootCmd.PersistentFlags().Lookup(CLUSTER_NAME_FLAG))

	RootCmd.PersistentFlags().String(KUBECONTEXT_FLAG, "", "name of the kubeconfig context to use")
	viper.BindPFlag(KUBECONTEXT_FLAG, RootCmd.PersistentFlags().Lookup(KUBECONTEXT_FLAG))

	home := homedir.HomeDir()
	RootCmd.PersistentFlags().String(KUBECONFIG_FLAG, filepath.Join(home, ".kube", "config"), "path to the kubeconfig file")
	viper.BindPFlag(KUBECONFIG_FLAG, RootCmd.PersistentFlags().Lookup(KUBECONFIG_FLAG))
	viper.BindEnv(KUBECONFIG_FLAG)
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

	plc_cluster := viper.Get("CLUSTER_NAME")

	os.Setenv("PL_CLOUD_ADDR", fmt.Sprintf("%s.getanton.com", plc_cluster))

	pl_cloud_addr := viper.Get("PL_CLOUD_ADDR") //, fmt.Sprintf("%s.getanton.com", plc_cluster))

	if pl_cloud_addr == nil {
		err := fmt.Errorf("PL_CLOUD_ADDR not set, set it to \n\nexport PL_CLOUD_ADDR=%s.getanton.com", plc_cluster)
		ui.LogAndPrintError(err)
		panic(err)
	}

	viper.Set("PLC_CLUSTER", fmt.Sprintf("gke_zerok-dev_us-west1-b_%s", plc_cluster))
}

func getZerokDirDefaultPath() string {
	var err error

	var baseDir string
	if baseDir, err = os.UserHomeDir(); err != nil {
		baseDir = os.TempDir()
	}

	return fmt.Sprintf("%s%c%s", baseDir, os.PathSeparator, zerok_dir_name)
}
