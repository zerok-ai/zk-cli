package install

import (
	"context"
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

	diSpinnerText = "installing zerok daemon and associated CRDs"
	diSuccessText = "zerok daemon installed successfully"
	diFailureText = "failed to install zerok daemon"

	waitTimeForService = 160

	DevKeyFlag        = "dev"
	DevKeyEnvFlag     = "ZK_DEV"
	VersionKeyFlag    = "zkVersion"
	VersionKeyEnvFlag = "ZK_VERSION"

	zkInstallClient    string = "/base-charts/install.sh"
	zkInstallDevClient string = "/scripts/install-dev.sh"

	pxVizierDevModeSetup string = "/zpx/scripts/setup-vizier-export.sh"
	pxVizierYaml         string = "/zpx/scripts/modified/vizier/exported-vizier.yaml"
	cliVizierYaml        string = "/vizier/vizier.yaml"
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

func GetCloudAddress() string {
	zkCloudAddr := viper.Get(ZkCloudAddressFlag)
	if zkCloudAddr == nil {
		//TODO replace this with the production address
		zkCloudAddr = "zkcloud02.getanton.com"
		viper.Set(ZkCloudAddressFlag, zkCloudAddr)
	}
	return zkCloudAddr.(string)
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

func InstallPXOperator(ctx context.Context, apiKey string) (err error) {

	// start deployment in background
	go func() {

		ui.GlobalWriter.PrintflnWithPrefixArrow("installing operator for managing data store")

		var out string
		cmd := utils.GetBackendCLIPath() + " deploy"
		out, err = shell.ShelloutWithSpinner(cmd, diSpinnerText, diSuccessText, diFailureText)

		if err != nil {
			filePath, _ := utils.DumpError(out)
			ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
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

	//// install the kustomization for vizier over the default code
	//out, err := shell.Shellout("VIZIER_TAG="+vizierTag+" ./"+pxVizierDevModeSetup, false)
	//if err != nil {
	//	filePath, _ := utils.DumpError(out)
	//	ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("vizier installation failed, Check %s for details\n", filePath))
	//}

	return nil
}

func InstallVizier(vizierTag string) error {
	out, err := shell.Shellout("VIZIER_TAG="+vizierTag+" ./"+pxVizierDevModeSetup, false)
	if err != nil {
		filePath, _ := utils.DumpError(out)
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("vizier export failed, Check %s for details\n", filePath))
		return err
	}

	out, err = shell.Shellout("cp ./"+pxVizierYaml+" ./"+cliVizierYaml, false)
	if err != nil {
		filePath, _ := utils.DumpError(out)
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("vizier copy failed, Check %s for details\n", filePath))
		return err
	}

	out, err = shell.Shellout("kubectl apply -f ./"+cliVizierYaml, true)
	if err != nil {
		filePath, _ := utils.DumpError(out)
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("vizier install failed, Check %s for details\n", filePath))
	}
	return err
}

func InstallZKServices(apiKey, clusterKey string) error {
	var inputToShellFile, shellFile string
	zkCloudAddr := viper.Get(ZkCloudAddressFlag).(string)
	if viper.Get(DevKeyFlag) == true {
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
			" PX_API_KEY=" + apiKey +
			" PX_CLUSTER_KEY=" + clusterKey
	} else {
		shellFile = zkInstallClient
		inputToShellFile = " ZK_CLOUD_ADDR=" + zkCloudAddr +
			" PX_API_KEY=" + apiKey +
			" PX_CLUSTER_KEY=" + clusterKey +
			" APP_NAME=zk-client"
	}
	return internal.ExecuteShellFile(shell.GetPWD()+shellFile, inputToShellFile, "failed to install zk_operator", "zk_operator installed successfully")
}