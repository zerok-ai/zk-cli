package cmd

import (
	
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/auth"
	
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/shell"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	SAME_CLUSTER_SETUP_FLAG string = "same_cluster_setup"
	PX_DEV_MODE_FLAG        string = "dev_mode"
	PX_DOMAIN               string = "px_domain"

	pxInstallCLI              string = "/scripts/install-px-cli.sh"
	pxInstallOperator         string = "/zpx/scripts/setup-operator.sh"
	pxInstallBackendPresetup  string = "/zpx/scripts/presetup-backend.sh"
	pxInstallBackend          string = "/zpx/scripts/setup-backend.sh"
	pxInstallBackendPostsetup string = "/zpx/scripts/postsetup-backend.sh"

	px_repo string = "us-west1-docker.pkg.dev/zerok-dev/pixie-dev"

	// Add the following to viper
	// SETUP_NGINX_INGRESS_WAIT_TIME=30
	// SETUP_SECRETS_WAIT_TIME=35
	// SETUP_CERT_MANAGER_WAIT_TIME=15
	// ASK_USER=0
	// SAME_CLUSTER_SETUP=1
	// PIXIE_DEV_MODE=0
)

var (
	sameClusterSetup bool
	pxDevMode        int = 0
	pxDomain         string
	pxClusterName    string

	ZKBackendCmd = &cobra.Command{
		Use:     "backend",
		Short:   "Install ZeroK components",
		PreRunE: preRunE,
		RunE:    RunZKBackendCmd,
	}

	ZkBackendCLICmd = &cobra.Command{
		Use:     "cli",
		Short:   "Install PX CLI",
		PreRunE: preRunE,
		RunE:    RunBackendCLICmd,
	}

	ZkBackendOperatorCmd = &cobra.Command{
		Use:     "operator",
		Short:   "Install PX operator",
		PreRunE: preRunE,
		RunE:    RunPxOperatorCmd,
	}

	ZkBackendPostsetupCmd = &cobra.Command{
		Use:     "postsetup",
		Short:   "Run backend post setup tasks",
		PreRunE: preRunE,
		RunE:    RunZkBackendPostsetupCmd,
	}

	ZKDatastore = &cobra.Command{
		Use:     "datastore",
		Short:   "Install zerok datastore",
		PreRunE: preRunE,
		RunE:    RunZkDatastoreCmd,
	}
)

func init() {

	ZKBackendCmd.AddCommand(ZkBackendCLICmd)
	ZKBackendCmd.AddCommand(ZkBackendOperatorCmd)
	ZKBackendCmd.AddCommand(ZkBackendPostsetupCmd)
	ZKBackendCmd.AddCommand(ZKDatastore)

	// same cluster setup
	ZKBackendCmd.PersistentFlags().BoolVarP(&sameClusterSetup, SAME_CLUSTER_SETUP_FLAG, "", true, "In case you would like to setup backend in the same cluster")
	viper.BindPFlag(SAME_CLUSTER_SETUP_FLAG, ZKBackendCmd.PersistentFlags().Lookup(SAME_CLUSTER_SETUP_FLAG))

	ZKBackendCmd.PersistentFlags().IntVarP(&pxDevMode, PX_DEV_MODE_FLAG, "", 0, "dev mode")
	viper.BindPFlag(PX_DEV_MODE_FLAG, ZKBackendCmd.PersistentFlags().Lookup(PX_DEV_MODE_FLAG))

	ZKBackendCmd.PersistentFlags().StringVarP(&pxDomain, PX_DOMAIN, "", "", "domain for px")
	viper.BindPFlag(PX_DOMAIN, ZKBackendCmd.PersistentFlags().Lookup(PX_DOMAIN))

}

func preRunE(cmd *cobra.Command, args []string) error {
	pxClusterName = viper.GetString("PX_CLUSTER_NAME")

	//	#Cloud proxy service domain
	//
	// Port forward
	if sameClusterSetup {
		pxDomain = fmt.Sprintf("%s.testdomain.com", pxClusterName)
		// PX_CLUSTER_NAME="$CLUSTER_NAME"
	} else {
		pxDomain = "pxtest.getanton.com"
	}

	// PL_CLOUD_ADDR=$PX_DOMAIN
	// if [ "$PIXIE_DEV_MODE" == '1' ]
	// then
	//
	//	PL_CLOUD_ADDR=$PX_DOMAIN:443
	//	PL_TESTING_ENV=dev
	//
	// fi

	return nil
}

func RunZkBackendPostsetupCmd(cmd *cobra.Command, args []string) error {
	_, err := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallBackendPostsetup, "PX backend postsetup done")
	return err
}

func RunBackendCLICmd(cmd *cobra.Command, args []string) error {
	_, err := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallCLI, "PX cli installed successfully")
	return err
}

func RunPxOperatorShellCmd(cmd *cobra.Command, args []string) error {
	_, err := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallOperator, "PX operator installed successfully")
	return err
}

func RunZkDatastoreCmd(cmd *cobra.Command, args []string) error {
	ui.GlobalWriter.PrintNoticeMessage("Running pre-setup scripts ")
	_, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+pxInstallBackendPresetup, "PX backend presetup done")
	if err != nil {
		return err
	}

	_, err1 := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallBackend, "PX backend installed successfully")
	return err1
}

func RunZKBackendCmd(cmd *cobra.Command, args []string) error {

	var err error

	if err := RunBackendCLICmd(cmd, args); err != nil {
		return err
	}

	if err = RunZkDatastoreCmd(cmd, args); err != nil {
		return err
	}

	for x := 0; x < 5; x++ {
		time.Sleep(30 * time.Second)

		// check whether pods are up
		fmt.Println("Checking readiness of the backend pods")
	}
	ui.GlobalWriter.PrintSuccessMessageln("Pods ready")

	if err = RunPxOperatorShellCmd(cmd, args); err != nil {
		return err
	}

	// if err = RunZkBackendPostsetupCmd(cmd, args); err != nil {
	// 	return err
	// }

	return nil
}

func RunPxOperatorCmd(cmd *cobra.Command, args []string) error {

	err := ValidateAndConfirm(cmd.Context())
	if err != nil {
		return err
	}


	// if chartValues, err = generateChartValues(chartValues, clusterName, storageProvision.PersistentStorage, deployableNodes, tolerations, sentryHelmContext); err != nil {
	// 	return err
	// }

	// if err = installHelmRelease(ctx, helmClient, releaseName, chart, chartValues); err != nil {
	// 	return err
	// }

	// err = validateInstall(ctx, kubeClient, namespace, chart.AppVersion(), &auth0Token, clusterName, len(deployableNodes), storageProvision.PersistentStorage, sentryHelmContext)
	// reportPodsStatus(ctx, kubeClient, namespace, sentryHelmContext)

	// if err != nil {
	// 	ui.GlobalWriter.PrintflnWithPrefixln("Installation takes longer then expected, you can check the status using \"kubectl get pods -n %s\"", namespace)
	// 	ui.GlobalWriter.PrintUrl(fmt.Sprintf("If pods in %q namespce are running, Check out: ", namespace), fmt.Sprintf("%s/?clusterId=%s&viewType=Overview\n", GROUNDCOVER_URL, clusterName))
	// 	return err
	// }

	ui.GlobalWriter.PrintlnWithPrefixln("\nThings are cooling down. Approaching ZeroK!")
	// utils.TryOpenBrowser(ui.GlobalWriter, "Check out: ", fmt.Sprintf("%s/?clusterId=%s&viewType=Overview", GROUNDCOVER_URL, clusterName))
	// ui.GlobalWriter.PrintUrl(fmt.Sprintf("\n%s\n", JOIN_SLACK_MESSAGE), JOIN_SLACK_LINK)

	return nil
}
