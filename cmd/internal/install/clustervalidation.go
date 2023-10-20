package install

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"zkctl/cmd/internal"
	"zkctl/cmd/pkg/k8s"
	sentryUtils "zkctl/cmd/pkg/sentry"
	"zkctl/cmd/pkg/ui"
)

const (
	NamespaceFlag   = "namespace"
	KubeConfigFlag  = "kubeconfig"
	KubeContextFlag = "kube-context"

	NODES_NOT_FOUND = "no nodes found in the current cluster"
	SHOULD_VALIDATE = "should_validate"
	YES_FLAG        = "yes"
	NO              = "no"

	NO_PVC_FLAG        = "no-pvc"
	LOW_RESOURCES_FLAG = "low-resources"
)

// ValidateClusterAndTakeConsent validates the cluster and takes consent from the user to install. It returns the cluster name
func ValidateClusterAndTakeConsent(ctx context.Context) (string, error) {

	var err error

	namespace := viper.GetString(NamespaceFlag)
	kubeConfig := viper.GetString(KubeConfigFlag)
	kubeContext := viper.GetString(KubeContextFlag)

	sentryKubeContext := sentryUtils.NewKubeContext(kubeConfig, kubeContext)
	sentryKubeContext.SetOnCurrentScope()

	var kubeClient *k8s.Client
	if kubeClient, err = k8s.NewKubeClient(kubeConfig, kubeContext); err != nil {
		return "", err
	}

	var clusterSummary *k8s.ClusterSummary
	if clusterSummary, err = validateCluster(ctx, kubeClient, namespace, sentryKubeContext); err != nil {
		return "", err
	}
	clusterName := clusterSummary.ClusterName

	var nodesReport *k8s.NodesReport
	if nodesReport, err = validateNodes(ctx, kubeClient, sentryKubeContext); err != nil {
		return clusterName, err
	}

	deployableNodes, tolerations, err := getDeployableNodesAndTolerations(nodesReport, sentryKubeContext)
	if err != nil {
		return clusterName, err
	}

	storageProvision := k8s.GenerateStorageProvision(context.Background(), kubeClient, clusterSummary)

	if !promptInstallSummary(clusterSummary.ClusterName, clusterSummary.Namespace, len(deployableNodes),
		nodesReport.NodesCount(), storageProvision, tolerations) {
		return clusterName, internal.ErrExecutionAborted
	}

	return clusterName, nil
}

func validateCluster(ctx context.Context, kubeClient *k8s.Client, namespace string, sentryKubeContext *sentryUtils.KubeContext) (*k8s.ClusterSummary, error) {
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

func validateNodes(ctx context.Context, kubeClient *k8s.Client, sentryKubeContext *sentryUtils.KubeContext) (*k8s.NodesReport, error) {
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

func getDeployableNodesAndTolerations(nodesReport *k8s.NodesReport, sentryKubeContext *sentryUtils.KubeContext) ([]*k8s.NodeSummary, []map[string]interface{}, error) {
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

func promptTaints(tolerationManager *k8s.TolerationManager, sentryKubeContext *sentryUtils.KubeContext) ([]string, error) {
	var err error

	var taints []string
	if taints, err = tolerationManager.GetTaints(); err != nil {
		return nil, err
	}

	// Allowing all taints
	//allowedTaints := ui.GlobalWriter.MultiSelectPrompt("Do you want set tolerations to allow scheduling zerok on following taints:", taints, taints)
	allowedTaints := taints

	sentryKubeContext.TolerationsAndTaintsRatio = fmt.Sprintf("%d/%d", len(allowedTaints), len(taints))
	sentryKubeContext.SetOnCurrentScope()
	sentryUtils.SetTagOnCurrentScope(sentryUtils.TAINTED_TAG, "true")

	return allowedTaints, nil
}

func promptInstallSummary(clusterName string, namespace string, deployableNodesCount, nodesCount int, storageProvision k8s.StorageProvision, tolerations []map[string]interface{}) bool {
	ui.GlobalWriter.PrintlnWithPrefixln("Installing Zerok:")

	if !storageProvision.PersistentStorage {
		ui.GlobalWriter.Printf("Using emptyDir storage, reason: %s\n", storageProvision.Reason)
	}
	t, _ := json.Marshal(tolerations)

	promptMessage := fmt.Sprintf(
		"Deploy Zerok (cluster: %s, compatible nodes: %d/%d, tolerations: %s)",
		clusterName, deployableNodesCount, nodesCount, string(t),
	)

	yesFlag := viper.GetBool(YES_FLAG)
	return ui.GlobalWriter.YesNoPrompt(promptMessage, yesFlag)
}
