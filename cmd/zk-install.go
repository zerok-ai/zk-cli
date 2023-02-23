package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/k8s"
	sentry_utils "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/sentry"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
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

	CONTEXT 							= "ctx"	
	SHOULD_PRINT_SUCCESS_MESSAGE 		= "should_print_success_msg"	
	NODES_NOT_FOUND 					= "no nodes found in the current cluster"	
	SHOULD_VALIDATE 					= "should_validate"	
	YES 								= "yes"	
	NO 									= "no"	
)

type ContextKey struct {
	key string
}

var Ckey ContextKey = ContextKey {key: CONTEXT}

type Values struct {
	m map[string]string
}

func (v Values) Get(key string) string {
	return v.m[key]
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install ZeroK",
	RunE:  runInstallCmd,
}

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.AddCommand(ZKBackendCmd)
	installCmd.AddCommand(ZkOperatorCmd)
}

func runInstallCmd(cmd *cobra.Command, args []string) error {

	err := ValidateAndConfirm(cmd.Context())
	if err != nil {
		return err
	}
	
	var values Values = Values{map[string]string{
		SHOULD_VALIDATE: NO,
		SHOULD_PRINT_SUCCESS_MESSAGE: NO,
	}}
	
	var ctx context.Context = context.WithValue(context.Background(),  Ckey, values)
	
	// ZkOperatorCmd.ExecuteContext(ctx)
	ZKBackendCmd.ExecuteContext(ctx)

	return nil
}

func ValidateAndConfirm(ctx context.Context) error {
	
	var err error

	namespace := viper.GetString(NAMESPACE_FLAG)
	kubeconfig := viper.GetString(KUBECONFIG_FLAG)
	kubecontext := viper.GetString(KUBECONTEXT_FLAG)
	ctxValue := ctx.Value(Ckey)
	
	sentryKubeContext := sentry_utils.NewKubeContext(kubeconfig, kubecontext)
	sentryKubeContext.SetOnCurrentScope()

	var kubeClient *k8s.Client
	if kubeClient, err = k8s.NewKubeClient(kubeconfig, kubecontext); err != nil {
		return err
	}

	var clusterSummary *k8s.ClusterSummary
	if clusterSummary, err = validateCluster(ctx, kubeClient, namespace, sentryKubeContext); err != nil {
		return err
	}

	var nodesReport *k8s.NodesReport
	if nodesReport, err = validateNodes(ctx, kubeClient, sentryKubeContext); err != nil {
		return err
	}

	deployableNodes, tolerations, err := getDeployableNodesAndTolerations(nodesReport, sentryKubeContext)
	if err != nil {
		return err
	}

	dn, _ := json.Marshal(deployableNodes)
	ui.GlobalWriter.Printf("deployableNodes name = %s", string(dn))

	t, _ := json.Marshal(tolerations)
	ui.GlobalWriter.PrintNoticeMessage(fmt.Sprintf("tolerations = %s", string(t)))

	storageProvision := k8s.GenerateStorageProvision(context.Background(), kubeClient, clusterSummary)
	if viper.GetBool(NO_PVC_FLAG) {
		storageProvision.PersistentStorage = false
		storageProvision.Reason = "user used --no-pvc flag"
	}

	if storageProvision.PersistentStorage {
		sentry_utils.SetTagOnCurrentScope(sentry_utils.PERSISTENT_STORAGE_TAG, "true")
	} else {
		sentry_utils.SetTagOnCurrentScope(sentry_utils.PERSISTENT_STORAGE_TAG, "false")
	}

	var shouldInstall bool
	if (ctxValue != nil){
		shouldInstall = ctxValue.(Values).Get(SHOULD_VALIDATE) == YES
		if !shouldInstall {
			return ErrExecutionAborted
		}
	}
	
	if shouldInstall, err = promptInstallSummary(
		clusterSummary.ClusterName, clusterSummary.Namespace, len(deployableNodes), nodesReport.NodesCount(), storageProvision); err != nil {
		return err
	}

	if !shouldInstall {
		return ErrExecutionAborted
	}

	return nil
}

func validateCluster(ctx context.Context, kubeClient *k8s.Client, namespace string, sentryKubeContext *sentry_utils.KubeContext) (*k8s.ClusterSummary, error) {
	var err error

	var clusterSummary *k8s.ClusterSummary
	if clusterSummary, err = kubeClient.GetClusterSummary(namespace); err != nil {
		sentryKubeContext.ClusterReport = &k8s.ClusterReport{
			ClusterSummary: clusterSummary,
		}
		sentryKubeContext.SetOnCurrentScope()
		return nil, err
	}
	ui.GlobalWriter.Printf("Validating cluster compatibility for cluster: %s\n", clusterSummary.ClusterName)

	clusterReport := k8s.DefaultClusterRequirements.Validate(ctx, kubeClient, clusterSummary)

	sentryKubeContext.ClusterReport = clusterReport
	sentryKubeContext.SetOnCurrentScope()

	clusterReport.PrintStatus()

	if clusterReport.IsLocalCluster() {
		viper.Set(LOW_RESOURCES_FLAG, true)
	}

	if !clusterReport.IsCompatible {
		return clusterSummary, errors.New("can't continue with installation, cluster is not compatible for installation. Check solutions suggested by the CLI")
	}

	return clusterSummary, nil
}

func validateNodes(ctx context.Context, kubeClient *k8s.Client, sentryKubeContext *sentry_utils.KubeContext) (*k8s.NodesReport, error) {
	var err error

	ui.GlobalWriter.PrintlnWithPrefixln("Validating cluster nodes compatibility:")

	var nodesSummeries []*k8s.NodeSummary
	if nodesSummeries, err = kubeClient.GetNodesSummeries(ctx); err != nil {
		return nil, err
	}

	if len(nodesSummeries) == 0 {
		return nil, fmt.Errorf(NODES_NOT_FOUND)
	}

	sentryKubeContext.NodesCount = len(nodesSummeries)
	sentryKubeContext.SetOnCurrentScope()

	nodesReport := k8s.DefaultNodeRequirements.Validate(nodesSummeries)

	sentryKubeContext.SetNodesSamples(nodesReport)
	sentryKubeContext.SetOnCurrentScope()

	nodesReport.PrintStatus()

	if len(nodesReport.CompatibleNodes) == 0 || nodesReport.Schedulable.IsNonCompatible {
		return nil, errors.New("can't continue with installation, no compatible nodes for installation")
	}

	return nodesReport, nil
}

func getDeployableNodesAndTolerations(nodesReport *k8s.NodesReport, sentryKubeContext *sentry_utils.KubeContext) ([]*k8s.NodeSummary, []map[string]interface{}, error) {
	var err error

	tolerations := make([]map[string]interface{}, 0)
	deployableNodes := nodesReport.CompatibleNodes

	if len(nodesReport.TaintedNodes) > 0 {
		tolerationManager := &k8s.TolerationManager{
			TaintedNodes: nodesReport.TaintedNodes,
		}

		var allowedTaints []string
		if allowedTaints, err = promptTaints(tolerationManager, sentryKubeContext); err != nil {
			return nil, nil, err
		}

		if tolerations, err = tolerationManager.GetTolerationsMap(allowedTaints); err != nil {
			return nil, nil, err
		}

		var tolerableNodes []*k8s.NodeSummary
		if tolerableNodes, err = tolerationManager.GetTolerableNodes(allowedTaints); err != nil {
			return nil, nil, err
		}

		deployableNodes = append(deployableNodes, tolerableNodes...)
	}

	return deployableNodes, tolerations, nil
}

func promptTaints(tolerationManager *k8s.TolerationManager, sentryKubeContext *sentry_utils.KubeContext) ([]string, error) {
	var err error

	var taints []string
	if taints, err = tolerationManager.GetTaints(); err != nil {
		return nil, err
	}

	allowedTaints := ui.GlobalWriter.MultiSelectPrompt("Do you want set tolerations to allow scheduling groundcover on following taints:", taints, taints)

	sentryKubeContext.TolerationsAndTaintsRatio = fmt.Sprintf("%d/%d", len(allowedTaints), len(taints))
	sentryKubeContext.SetOnCurrentScope()
	sentry_utils.SetTagOnCurrentScope(sentry_utils.TAINTED_TAG, "true")

	return allowedTaints, nil
}

func promptInstallSummary(clusterName string, namespace string, deployableNodesCount, nodesCount int, storageProvision k8s.StorageProvision) (bool, error) {
	ui.GlobalWriter.PrintlnWithPrefixln("Installing groundcover:")

	if !storageProvision.PersistentStorage {
		ui.GlobalWriter.Printf("Using emptyDir storage, reason: %s\n", storageProvision.Reason)
	}

	var promptMessage string = fmt.Sprintf(
		"Deploy Zerok (cluster: %s, namespace: %s, compatible nodes: %d/%d)",
		clusterName, namespace, deployableNodesCount, nodesCount,
	)

	return ui.GlobalWriter.YesNoPrompt(promptMessage, false), nil
}

// func promptInstallSummary(isUpgrade bool, releaseName string, clusterName string, namespace string, release *helm.Release, chart *helm.Chart, deployableNodesCount, nodesCount int, storageProvision k8s.StorageProvision, sentryHelmContext *sentry_utils.HelmContext) (bool, error) {
// 	ui.GlobalWriter.PrintlnWithPrefixln("Installing groundcover:")

// 	if !storageProvision.PersistentStorage {
// 		ui.GlobalWriter.Printf("Using emptyDir storage, reason: %s\n", storageProvision.Reason)
// 	}

// 	var promptMessage string
// 	if isUpgrade {
// 		sentryHelmContext.Upgrade = isUpgrade
// 		sentryHelmContext.PreviousChartVersion = release.Version().String()
// 		sentry_utils.SetTagOnCurrentScope(sentry_utils.UPGRADE_TAG, strconv.FormatBool(isUpgrade))
// 		sentryHelmContext.SetOnCurrentScope()

// 		if chart.Version().GT(release.Version()) {
// 			promptMessage = fmt.Sprintf(
// 				"Your groundcover version is out of date (cluster: %s, namespace: %s, version: %s), The latest version is %s.\nDo you want to upgrade?",
// 				clusterName, namespace, release.Version(), chart.Version(),
// 			)
// 		} else {
// 			promptMessage = fmt.Sprintf(
// 				"Latest version of groundcover is already installed in your cluster! (cluster: %s, namespace: %s, version: %s).\nDo you want to redeploy?",
// 				clusterName, namespace, chart.Version(),
// 			)
// 		}
// 	} else {
// 		promptMessage = fmt.Sprintf(
// 			"Deploy groundcover (cluster: %s, namespace: %s, compatible nodes: %d/%d, version: %s)",
// 			clusterName, namespace, deployableNodesCount, nodesCount, chart.Version(),
// 		)
// 	}

// 	return ui.GlobalWriter.YesNoPrompt(promptMessage, !isUpgrade), nil
// }

