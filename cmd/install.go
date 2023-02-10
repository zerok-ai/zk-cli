package cmd

import (
	"time"

	// "errors"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
	// "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/k8s"
	install "github.com/zerok-ai/zk-cli/zkctl/cmd/install"
	// logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
)

const (
	VALUES_FLAG                         = "values"
	EXPERIMENTAL_FLAG                   = "experimental"
	NO_PVC_FLAG                         = "no-pvc"
	LOW_RESOURCES_FLAG                  = "low-resources"
	STORE_ISSUES_LOGS_ONLY_FLAG         = "store-issues-logs-only"
	STORE_ISSUES_LOGS_ONLY_KEY          = "storeIssuesLogsOnly"
	COMMIT_HASH_KEY_NAME_FLAG           = "git-commit-hash-key-name"
	REPOSITORY_URL_KEY_NAME_FLAG        = "git-repository-url-key-name"
	EXPERIMENTAL_PRESET_PATH            = "presets/agent/experimental.yaml"
	LOW_RESOURCES_NOTICE_MESSAGE_FORMAT = "We get it, you like things light 🪁\n   But since you’re deploying on a %s we’ll have to limit some of our features to make sure it’s smooth sailing.\n   For the full experience, try deploying on a different cluster\n"
	WAIT_FOR_GET_LATEST_CHART_FORMAT    = "Waiting for downloading latest chart to complete"
	WAIT_FOR_GET_LATEST_CHART_SUCCESS   = "Downloading latest chart completed successfully"
	WAIT_FOR_GET_LATEST_CHART_FAILURE   = "Latest chart download failed:"
	WAIT_FOR_GET_LATEST_CHART_TIMEOUT   = "Latest chart download timeout"
	GET_LATEST_CHART_POLLING_RETIRES    = 3
	GET_LATEST_CHART_POLLING_INTERVAL   = time.Second * 1
	GET_LATEST_CHART_POLLING_TIMEOUT    = time.Second * 10
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install ZeroK",
	RunE:  runInstallCmd,
}

func init() {
	Initialize(rootCmd)
}

func Initialize(root *cobra.Command) {
	root.AddCommand(installCmd)

	installCmd.AddCommand(install.ZKBackendCmd)
	installCmd.AddCommand(install.ZkOperatorCmd)
}

func runInstallCmd(cmd *cobra.Command, args []string) error {

	// Nothing here, how help
	cmd.Help()

	return nil
}

// func validateCluster(ctx context.Context, namespace string) error {
// 	// var err error

// 	ui.GlobalWriter.PrintlnWithPrefixln("Validating cluster compatibility:")

// 	// var clusterSummary *k8s.ClusterSummary
// 	// if clusterSummary, err = kubeClient.GetClusterSummary(namespace); err != nil {
// 	// 	// print the error
// 	// 	return err
// 	// }

// 	// clusterReport := k8s.DefaultClusterRequirements.Validate(ctx, kubeClient, clusterSummary)

// 	// clusterReport.PrintStatus()

// 	// if clusterReport.IsLocalCluster() {
// 	// 	viper.Set(LOW_RESOURCES_FLAG, true)
// 	// }

// 	// if !clusterReport.IsCompatible {
// 	// 	return errors.New("can't continue with installation, cluster is not compatible for installation. Check solutions suggested by the CLI")
// 	// }

// 	return nil
// }
