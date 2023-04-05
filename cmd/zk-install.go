package cmd

import (
	"os"
	"fmt"
	"time"
	"errors"
	"runtime"
	"strings"
	"context"
	"net/http"
	"encoding/json"

	"zkctl/cmd/pkg/k8s"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sentry_utils "zkctl/cmd/pkg/sentry"
)


type AuthPayload struct {
	Payload AuthRefreshToken `json:"payload"`
}

type AuthRefreshToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	OrgName   string `json:"orgName,omitempty"`
	OrgID     string `json:"orgID,omitempty"`
}

const (
	ZK_CLOUD_ADDRESS_FLAG   = "ZK_CLOUD_ADDRESS"
	API_KEY_FLAG            = "apikey"
	API_KEY_ENV_FLAG        = "API_KEY"
	API_KEY_WARNING_MESSAGE = "api key is not set. To continue, please get the apikey from zerok dashboard and paste below."
	API_KEY_QUESTION        = "enter api key"
	API_KEY_ERROR_MESSAGE   = "apikey is not set. Run the help command for more details"

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

	CONTEXT                      = "ctx"
	SHOULD_PRINT_SUCCESS_MESSAGE = "should_print_success_msg"
	NODES_NOT_FOUND              = "no nodes found in the current cluster"
	SHOULD_VALIDATE              = "should_validate"
	YES                          = "yes"
	NO                           = "no"

	APPLY_POLLING_RETRIES  = 1
	APPLY_POLLING_TIMEOUT  = time.Minute * 3
	APPLY_POLLING_INTERVAL = time.Second / 10

	waitTimeForPodsInSeconds         = 100
	waitTimeForService               = 60
	waitTimeForInstallationInSeconds = 240

	zkInstallOperator  string = "/operator/buildAndInstall.sh"
	pxDownloadAuthJson string = "/zpx/scripts/setup-px-auth-json.sh"
	pxSetupDomain      string = "/zpx/scripts/setup-domain.sh"
	pxSetupIngress     string = "/zpx/scripts/setup-ingress.sh"

	diSpinnerText = "installing zerok daemon and associated CRDs"
	diSuccessText = "zerok daemon installed successfully"
	diFailureText = "failed to install zerok daemon"
)

var (
	authAddress string
	apiKey      string
)

type ContextKey struct {
	key string
}

var Ckey ContextKey = ContextKey{key: CONTEXT}

type Values struct {
	m map[string]string
}

func (v Values) Get(key string) string {
	return v.m[key]
}

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install ZeroK",
	PreRunE: RunInstallPreCmd,
	RunE:    RunInstallCmd,
}

func init() {
	RootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().String(API_KEY_FLAG, "", "api key. This can also be set through environment variable "+API_KEY_ENV_FLAG+" instead of passing the parameter")
	viper.BindPFlag(API_KEY_FLAG, installCmd.PersistentFlags().Lookup(API_KEY_FLAG))
	viper.BindEnv(API_KEY_FLAG, API_KEY_ENV_FLAG)

	viper.BindEnv("PL_CLOUD_ADDR", ZK_CLOUD_ADDRESS_FLAG)
}

func CheckFlags() error {

	//Check if API Key is present
	apiKey = viper.Get(API_KEY_FLAG).(string)
	if apiKey == "" {
		ui.GlobalWriter.PrintlnWarningMessageln(API_KEY_WARNING_MESSAGE)
		apiKey = ui.GlobalWriter.QuestionPrompt(API_KEY_QUESTION)
		if apiKey == "" {
			return errors.New(API_KEY_ERROR_MESSAGE)
		}
		viper.Set(API_KEY_FLAG, apiKey)
	}

	// set cloud address as an env variable for shell scripts
	var pl_cloud_addr string = viper.Get(ZK_CLOUD_ADDRESS_FLAG).(string)
	if pl_cloud_addr == "" {
		pl_cloud_addr = "zkcloud02.getanton.com"
		viper.Set("ZK_CLOUD_ADDRESS_FLAG", pl_cloud_addr)
		os.Setenv("PL_CLOUD_ADDR", pl_cloud_addr)
	}

	//TODO: Change it to https
	authAddress = fmt.Sprintf("http://api.%s", pl_cloud_addr)
	return nil
}

func RunInstallPreCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	err := CheckFlags()
	if err != nil {
		return err
	}

	utils.InitializeFolders()

	err = validateAndConfirm(ctx)
	ui.GlobalWriter.PrintflnWithPrefixArrow("running pre-installation checks")
	if err != nil {
		ui.GlobalWriter.Printf("In PreInstallChecksAndTasks err %s \n", err)
		return err
	}

	// install px cli
	return installBackendCLI(ctx)
}

func RunInstallCmd(cmd *cobra.Command, args []string) error {

	_, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+zkInstallOperator, "zeroK operator installed successfully")
	if err != nil {
		ui.LogAndPrintError(fmt.Errorf("failed to install zkoperator: %v", err))
		return err
	}

	err = installPXOperator(cmd.Context(), apiKey)

	if err == nil {
		ui.GlobalWriter.PrintlnSuccessMessageln("installation done")
	}
	return err
}

func validateAndConfirm(ctx context.Context) error {

	var err error

	ctxValue := ctx.Value(Ckey)
	if ctxValue != nil && ctxValue.(Values).Get(SHOULD_VALIDATE) != YES {
		return nil
	}

	namespace := viper.GetString(NAMESPACE_FLAG)
	kubeconfig := viper.GetString(KUBECONFIG_FLAG)
	kubecontext := viper.GetString(KUBECONTEXT_FLAG)

	return nil

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

	if !promptInstallSummary(clusterSummary.ClusterName, clusterSummary.Namespace, len(deployableNodes),
		nodesReport.NodesCount(), storageProvision, tolerations) {
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

	allowedTaints := ui.GlobalWriter.MultiSelectPrompt("Do you want set tolerations to allow scheduling zerok on following taints:", taints, taints)

	sentryKubeContext.TolerationsAndTaintsRatio = fmt.Sprintf("%d/%d", len(allowedTaints), len(taints))
	sentryKubeContext.SetOnCurrentScope()
	sentry_utils.SetTagOnCurrentScope(sentry_utils.TAINTED_TAG, "true")

	return allowedTaints, nil
}

func promptInstallSummary(clusterName string, namespace string, deployableNodesCount, nodesCount int, storageProvision k8s.StorageProvision, tolerations []map[string]interface{}) bool {
	ui.GlobalWriter.PrintlnWithPrefixln("Installing Zerok:")

	if !storageProvision.PersistentStorage {
		ui.GlobalWriter.Printf("Using emptyDir storage, reason: %s\n", storageProvision.Reason)
	}
	t, _ := json.Marshal(tolerations)

	var promptMessage string = fmt.Sprintf(
		"Deploy Zerok (cluster: %s, compatible nodes: %d/%d, tolerations: %s)",
		clusterName, deployableNodesCount, nodesCount, string(t),
	)

	return ui.GlobalWriter.YesNoPrompt(promptMessage, true)
}

func installBackendCLI(ctx context.Context) error {

	use_version := "latest"
	artifact_bucket := "pixie-dev-public"
	if strings.Contains(use_version, "-") {
		artifact_bucket = "pixie-prod-artifacts"
	}
	artifact_base_path := fmt.Sprintf("https://storage.googleapis.com/%s/cli", artifact_bucket)

	err := utils.BackendCLIExists()
	if err != nil {
		ui.GlobalWriter.PrintflnWithPrefixArrow("installing full CLI support")

		artifact_name := "cli_darwin_universal"
		if runtime.GOOS == "linux" {
			artifact_name = "cli_linux_amd64"
		}
		url := fmt.Sprintf("%s/%s/%s", artifact_base_path, use_version, artifact_name)
		// ui.GlobalWriter.Printf("downloading %s to %s \n", url, backend_cli_filepath)

		err = utils.DownloadBackendCLI(url)
	} else {
		ui.GlobalWriter.PrintSuccessMessage("full CLI support already installed\n")
	}
	return err
}

func LoginToPX(ctx context.Context, apiKey string) error {

	ui.GlobalWriter.PrintflnWithPrefixArrow("trying to authenticate")
	path := "v1/p/auth/login"
	urlToBeCalled := fmt.Sprintf("%s/%s?apikey=%s", authAddress, path, apiKey)

	// download from server
	jsonResponse := new(AuthPayload)
	err := GetHTTPGETResponse(urlToBeCalled, jsonResponse)
	if err != nil {
		return fmt.Errorf("error in auth %v", err)
	}
	authTokenPayloadBytes, err := json.Marshal(jsonResponse.Payload)
	if err != nil {
		return err
	}

	//write to file
	authPath := utils.GetBackendAuthPath()
	os.Remove(authPath)
	return utils.WriteTextToFile(string(authTokenPayloadBytes), authPath)
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetHTTPGETResponse(url string, target interface{}) error {
	response, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

func installPXOperator(ctx context.Context, apiKey string) (err error) {

	// login to px
	if err = LoginToPX(ctx, apiKey); err != nil {
		// send to sentry and print
		return ui.GlobalWriter.Errorf("login error %v", err)
	}

	// wg.Add(1)

	go func() {

		// defer wg.Done()
		// go func() error {
		// deploy px operator
		ui.GlobalWriter.PrintflnWithPrefixArrow("installing operator for managing data store")

		var out string
		/*/
		cmd := utils.GetBackendCLIPath() + " deploy"
		out, err := shell.ExecWithLogsDurationAndSuccessM(cmd, "Zerok daemon installed successfully")
		/*/
		cmd := utils.GetBackendCLIPath() + " deploy"
		out, err = shell.ShelloutWithSpinner(cmd, diSpinnerText, diSuccessText, diFailureText)
		/**/

		filePath, _ := utils.DumpError(out)
		if err != nil {
			// send to sentry and print
			ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
			// wg.Done()
		}
		// return nil
		// }()

		// ui.GlobalWriter.PrintflnWithPrefixArrow("Waiting for %d sec for the installation", waitTimeForInstallationInSeconds)
		// time.Sleep(waitTimeForInstallationInSeconds * time.Second)
	}()

	if err != nil {
		return err
	}

	// set ingress to post deploy settings
	// ui.GlobalWriter.PrintflnWithPrefixArrow("Waiting for %d sec for the pods to get deployed", waitTimeForPodsInSeconds)

	for i := 0; i <= waitTimeForPodsInSeconds; i++ {
		time.Sleep(1 * time.Second)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	for i := 0; i <= waitTimeForService; i++ {
		time.Sleep(1 * time.Second)
		if err != nil {
			return err
		}
	}

	return nil
}
