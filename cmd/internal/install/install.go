package install

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
	"time"
	"zkctl/cmd/internal"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"
)

type AuthPayload struct {
	Payload AuthCompletePayload `json:"payload"`
}

type AuthCompletePayload struct {
	CliAuthPayload      AuthRefreshToken `json:"cliAuth"`
	OperatorAuthPayload OperatorToken    `json:"operatorAuth"`
}

type AuthRefreshToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	OrgName   string `json:"orgName,omitempty"`
	OrgID     string `json:"orgID,omitempty"`
}

type OperatorToken struct {
	ClusterKey string `json:"cluster_key"`
}

const (
	ApiKeyFlag           = "apikey"
	ApiKeyEnvFlag        = "ZK_API_KEY"
	ApiKeyWarningMessage = "api key is not set. To continue, please get the apikey from zerok dashboard and paste below."
	ApiKeyQuestion       = "enter api key"
	ApiKeyErrorMessage   = "apikey is not set. Run the help command for more details"

	ZkCloudAddressFlag = "ZK_CLOUD_ADDRESS"

	installDataStore            = "installing data stores"
	installDataStoreSuccessText = "installation of data stores completed successfully"
	installDataStoreFailureText = "installation of data stores failed"

	preInstCreatingCopyOfSecret = "preparing to install zerok daemon and associated CRDs."
	preInstSuccessText          = "pre-installation step completed successfully"
	preInstFailureText          = "pre-installation step failed"

	diSpinnerText = "installing zerok daemonset and associated CRDs"
	diSuccessText = "zerok daemonset installed successfully"
	diFailureText = "failed to install zerok daemonset"

	waitTimeForService = 160

	VersionKeyFlag    = "zkVersion"
	VersionKeyEnvFlag = "ZK_VERSION"

	zkInstallClient    string = "/scripts/install.sh"
	zkInstallDevClient string = "/helm-charts/install-dev.sh"
	cliVizierYaml      string = "/vizier/vizier.yaml"

	zkInstallStores string = "/db-helm-charts/install-db.sh"
)

func GetAPIKey() string {
	var apiKey string
	if viper.Get(ApiKeyFlag) != nil {
		apiKey = viper.Get(ApiKeyFlag).(string)
	}
	if apiKey == "" {
		ui.GlobalWriter.PrintlnWarningMessageln(ApiKeyWarningMessage)
		apiKey = ui.GlobalWriter.QuestionPrompt(ApiKeyQuestion)
		if apiKey != "" {
			viper.Set(ApiKeyFlag, apiKey)
		}
	}
	return apiKey
}

func DownloadAndInstallPXCLI(ctx context.Context) error {

	useVersion := "latest"
	artifactBucket := "pixie-dev-public"
	if strings.Contains(useVersion, "-") {
		artifactBucket = "pixie-prod-artifacts"
	}
	artifactBasePath := fmt.Sprintf("https://storage.googleapis.com/%s/cli", artifactBucket)

	err := utils.BackendCLIExists()
	if err != nil {
		ui.GlobalWriter.PrintflnWithPrefixArrow("installing full CLI support")

		artifactName := "cli_darwin_universal"
		if runtime.GOOS == "linux" {
			artifactName = "cli_linux_amd64"
		}
		url := fmt.Sprintf("%s/%s/%s", artifactBasePath, useVersion, artifactName)
		err = utils.DownloadBackendCLI(url)
	} else {
		ui.GlobalWriter.PrintSuccessMessage("full CLI support already installed\n")
	}
	return err
}

func LoginToPX(authAddress, apiKey, clusterName string) (string, error) {

	ui.GlobalWriter.PrintflnWithPrefixArrow("trying to authenticate")
	path := "v1/p/auth/login"
	urlToBeCalled := fmt.Sprintf("%s/%s?apikey=%s&clusterName=%s", authAddress, path, apiKey, clusterName)

	// download from server
	jsonResponse := new(AuthPayload)
	err := GetHTTPGETResponse(urlToBeCalled, jsonResponse)
	// ui.GlobalWriter.PrintflnWithPrefixlnAndArrow("urlToBeCalled = %s\njsonResponse=%v", urlToBeCalled, jsonResponse)
	if err != nil {
		return "", fmt.Errorf("error in auth %v", err)
	}

	if jsonResponse.Payload.CliAuthPayload.Token == "" {
		//TODO:AVIN Put a proper error message
		return "", errors.New("error in auth")
	}
	authTokenPayloadBytes, err := json.Marshal(jsonResponse.Payload.CliAuthPayload)
	if err != nil {
		return "", err
	}
	clusterKeyLocal := jsonResponse.Payload.OperatorAuthPayload.ClusterKey
	if clusterKeyLocal == "" {
		//TODO:AVIN Put proper error message
		return "", errors.New("error in auth")
	}

	//write to file
	authPath := utils.GetBackendAuthPath()
	err = os.Remove(authPath)
	if err != nil {
		return clusterKeyLocal, err
	}
	return clusterKeyLocal, utils.WriteTextToFile(string(authTokenPayloadBytes), authPath)
}

func copySecrets() error {
	var out string
	cmd := "kubectl create namespace pl"

	out, err := shell.Shellout(cmd, true)
	if err != nil {
		//	 do nothing
		internal.DumpErrorAndPrintLocation(out)
	}

	internal.DumpErrorAndPrintLocation("before reading secret")

	// check if the secret exists
	cmd = "kubectl get secret redis -n zk-client -o jsonpath='{.data.redis-password}' | base64 -d"
	secret, err := shell.Shellout(cmd, false)
	if err != nil {
		internal.DumpErrorAndPrintLocation(secret)
		return err
	}

	internal.DumpErrorAndPrintLocation("-----------" + secret)

	cmd = fmt.Sprintf("kubectl create secret generic redis -n pl --from-literal=redis-password=%s", secret)
	out, err = shell.Shellout(cmd, false)
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
		return err
	}

	return nil
}

func copyConfigmaps() error {
	var out string
	cmd := "kubectl create namespace pl"

	out, err := shell.Shellout(cmd, true)
	if err != nil {
		//	 do nothing
		internal.DumpErrorAndPrintLocation(out)
	}

	internal.DumpErrorAndPrintLocation("before reading configmap")

	// check if the secret exists
	cmd = "kubectl get configmap zk-redis-config -n zk-client -o jsonpath='{.data.redisHost}'"
	config, err := shell.Shellout(cmd, false)
	if err != nil {
		internal.DumpErrorAndPrintLocation(config)
		return err
	}

	internal.DumpErrorAndPrintLocation("-----------" + config)

	cmd = fmt.Sprintf("kubectl create configmap zk-redis-config -n pl --from-literal=redisHost=%s", config)
	out, err = shell.Shellout(cmd, false)
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
		return err
	}

	return nil
}

func InstallDataStores() error {
	// Install zk-client data stores
	return internal.ExecuteShellFileWithSpinner(shell.GetPWD()+zkInstallStores, " ", installDataStore, installDataStoreSuccessText, installDataStoreFailureText)
}

func InstallPXOperator() (err error) {

	// start deployment in background
	go func() {

		ui.GlobalWriter.PrintflnWithPrefixArrow("installing operator for managing data store")

		_, err = shell.RunWithSpinner(copySecrets, preInstCreatingCopyOfSecret, preInstSuccessText, preInstFailureText)
		if err != nil {
			return
		}

		_, err = shell.RunWithSpinner(copyConfigmaps, preInstCreatingCopyOfSecret, preInstSuccessText, preInstFailureText)
		if err != nil {
			return
		}

		var out string
		cmd := utils.GetBackendCLIPath() + " deploy"
		out, err = shell.ShelloutWithSpinner(cmd, diSpinnerText, diSuccessText, diFailureText)

		if err != nil {
			internal.DumpErrorAndPrintLocation(out)
		}
	}()

	// wait for `waitTimeForService` seconds for the service
	// TODO: replace this code with proper health check
	for i := 0; i <= waitTimeForService; i++ {
		time.Sleep(1 * time.Second)
		if err != nil {
			return err
		}
	}

	return nil
}

func InstallVizier() error {

	patch := func() error {
		out, err := shell.Shellout("kubectl apply -f "+shell.GetPWD()+cliVizierYaml, false)
		if err != nil {
			internal.DumpErrorAndPrintLocation(fmt.Sprintf("vizier install failed, %s", out))
		}
		return err
	}

	out, err := shell.RunWithSpinner(patch, "applying final patches", "patches applied successfully", "patches failed")
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
	}
	return err
}

func InstallZKServices(apiKey, clusterKey string) error {
	var inputToShellFile, shellFile string
	zkCloudAddr := viper.Get(ZkCloudAddressFlag).(string)
	if viper.Get(internal.DevKeyFlag) == true {
		ui.GlobalWriter.Println("zerok dev mode is enabled")

		//get the versions
		versions := viper.Get(VersionKeyFlag).(string)
		keyValueMap := internal.GetKVPairsFromCSV(versions)

		shellFile = zkInstallDevClient
		inputToShellFile = " ZK_CLOUD_ADDR=" + zkCloudAddr +
			" ZK_SCENARIO_MANAGER_VERSION=" + keyValueMap["zk-scenario-manager"] +
			" ZK_AXON_VERSION=" + keyValueMap["zk-axon"] +
			" ZK_DAEMONSET_VERSION=" + keyValueMap["zk-daemonset"] +
			" ZK_GPT_VERSION=" + keyValueMap["zk-gpt"] +
			" ZK_WSP_CLIENT_VERSION=" + keyValueMap["zk-wsp-client"] +
			" ZK_OPERATOR_VERSION=" + keyValueMap["zk-operator"] +
			" ZK_APP_INIT_CONTAINERS_VERSION=" + keyValueMap["zk-app-init-containers"] +
			" PX_API_KEY=" + apiKey +
			" PX_CLUSTER_KEY=" + clusterKey
	} else {
		shellFile = zkInstallClient
		inputToShellFile = " ZK_CLOUD_ADDR=" + zkCloudAddr +
			" PX_API_KEY=" + apiKey +
			" PX_CLUSTER_KEY=" + clusterKey +
			" APP_NAME=zk-client"
	}
	return internal.ExecuteShellFileWithSpinner(shell.GetPWD()+shellFile, inputToShellFile, "installing zk operator", "zk_operator installed successfully", "failed to install zk_operator")
}
