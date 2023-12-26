package install

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	"zkctl/cmd/internal"
	"zkctl/cmd/internal/shell"
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
	ClusterID  string `json:"cluster_id"`
}

const (
	EbpfMemoryFlag    = "ebpf-mem"
	EbpfMemoryEnvFlag = "ZK_EBPF_MEM"

	ApiKeyFlag    = "apikey"
	ApiKeyEnvFlag = "ZK_API_KEY"

	PostgresVersionKeyFlag    = "postgres-version"
	PostgresVersionKeyEnvFlag = "ZK_POSTGRES_VERSION_KEY"

	RedisVersionKeyFlag    = "redis-version"
	RedisVersionKeyEnvFlag = "ZK_REDIS_VERSION_KEY"

	ApiKeyWarningMessage = "api key is not set. To continue, please get the apikey from zerok dashboard and paste below."
	ApiKeyQuestion       = "enter api key"
	ApiKeyErrorMessage   = "apikey is not set. Run the help command for more details"

	ZkCloudAddressFlag = "ZK_CLOUD_ADDRESS"

	installDataStore            = "installing data stores"
	installDataStoreSuccessText = "installation of data stores completed successfully"
	installDataStoreFailureText = "installation of data stores failed"

	installOlm            = "installing olm"
	installOlmSuccessText = "installation of olm completed successfully"
	installOlmFailureText = "installation of olm failed"

	preInstCreatingCopyOfSecret = "preparing to install zerok daemon and associated CRDs."
	preInstSuccessText          = "pre-installation step completed successfully"
	preInstFailureText          = "pre-installation step failed"

	diSpinnerText = "installing zerok daemonset and associated CRDs"
	diSuccessText = "zerok daemonset installed successfully"
	diFailureText = "failed to install zerok daemonset"

	zksSpinnerText = "installing zerok operator"
	zksSuccessText = "zerok operator installed successfully"
	zksFailureText = "failed to install zerok operator"

	zkEbpfSpinnerText = "installing zerok ebpf layer"
	zkEbpfSuccessText = "zerok ebpf layer installed successfully"
	zkEbpfFailureText = "failed to install zerok ebpf layer"

	waitTimeForService = 160

	VersionKeyFlag    = "zkVersion"
	VersionKeyEnvFlag = "ZK_VERSION"

	olmInstall      string = "scripts/install-olm.sh"
	zkInstallClient string = "scripts/install.sh"
	//zkInstallEbpf         string = "scripts/install-ebpf.sh"
	zkInstallEbpf         string = "zk-ebpf/install-local.sh"
	zkInstallDevClient    string = "helm-charts/install-dev.sh"
	zkInstallSpreadClient string = "scripts/install-spread.sh"
	cliVizierYaml         string = "vizier/vizier.yaml"

	zkInstallStores string = "scripts/install-db.sh"
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

func GetEbpfMemory() string {
	var ebpfMemory string = ""
	if viper.Get(EbpfMemoryFlag) != nil {
		ebpfMemory = viper.Get(EbpfMemoryFlag).(string)
	}
	return ebpfMemory
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

type ClusterMetadata struct {
	ClientVersions map[string]string `json:"clientVersion"`
	NumOfNodes     int64             `json:"numNodes"`
}

func LoginToPX(authAddress, apiKey, clusterName string, clusterMetadata ClusterMetadata) (string, string, error) {

	ui.GlobalWriter.PrintflnWithPrefixArrow("trying to authenticate")
	path := "v1/p/auth/login"
	//authAddress = "http://localhost:80"
	urlToBeCalled := fmt.Sprintf("%s/%s?apikey=%s&clusterName=%s", authAddress, path, apiKey, clusterName)

	// download from server
	jsonResponse := new(AuthPayload)
	err := GetHTTPPOSTResponse(urlToBeCalled, clusterMetadata, jsonResponse)
	// ui.GlobalWriter.PrintflnWithPrefixlnAndArrow("urlToBeCalled = %s\njsonResponse=%v", urlToBeCalled, jsonResponse)
	if err != nil {
		return "", "", fmt.Errorf("error in auth %v", err)
	}

	if jsonResponse.Payload.CliAuthPayload.Token == "" {
		//TODO:AVIN Put a proper error message
		return "", "", errors.New("error in auth")
	}
	authTokenPayloadBytes, err := json.Marshal(jsonResponse.Payload.CliAuthPayload)
	if err != nil {
		return "", "", err
	}
	clusterKeyLocal := jsonResponse.Payload.OperatorAuthPayload.ClusterKey
	if clusterKeyLocal == "" {
		//TODO:AVIN Put proper error message
		return "", "", errors.New("error in auth")
	}

	clusterIdLocal := jsonResponse.Payload.OperatorAuthPayload.ClusterID
	if clusterIdLocal == "" {
		//TODO:AVIN Put proper error message
		return "", "", errors.New("error in auth")
	}

	//write to file
	authPath := utils.GetBackendAuthPath()
	err = os.Remove(authPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return clusterKeyLocal, clusterIdLocal, err
		}
	}
	return clusterKeyLocal, clusterIdLocal, utils.WriteTextToFile(string(authTokenPayloadBytes), authPath)
}

func copySecrets() error {
	var out string
	cmd := "kubectl create namespace pl"

	out, err := shell.Shellout(cmd)
	if err != nil {
		//	 do nothing
		internal.DumpErrorAndPrintLocation(out)
	}

	internal.DumpErrorAndPrintLocation("before reading secret")

	// check if the secret exists
	cmd = "kubectl get secret redis -n zk-client -o jsonpath='{.data.redis-password}' | base64 -d"
	secret, err := shell.Shellout(cmd)
	if err != nil {
		errMessage := fmt.Sprintf("Error in fetching secret %v", err)
		internal.DumpErrorAndPrintLocation(errMessage)
		return err
	}

	//Check if secret already exists
	cmd = "kubectl get secret redis -n pl -o jsonpath='{.data.redis-password}'"
	out, secretExistsErr := shell.Shellout(cmd)
	if secretExistsErr == nil {
		//TODO: Remove
		internal.DumpErrorAndPrintLocation("Secret already exists " + out)
		return nil
	} else {
		internal.DumpErrorAndPrintLocation("Error in checking secret: " + out)
	}

	internal.DumpErrorAndPrintLocation("-----------" + secret)
	cmd = fmt.Sprintf("kubectl create secret generic redis -n pl --from-literal=redis-password=%s", secret)
	out, err = shell.Shellout(cmd)
	if err != nil {
		internal.DumpErrorAndPrintLocation("Error in creating secret: " + out)
		return err
	}

	return nil
}

func copyConfigmaps() error {
	var out string
	cmd := "kubectl create namespace pl"

	out, err := shell.Shellout(cmd)
	if err != nil {
		//	 do nothing
		internal.DumpErrorAndPrintLocation(out)
	}

	internal.DumpErrorAndPrintLocation("before reading configmap")

	// check if the secret exists
	cmd = "kubectl get configmap zk-redis-config -n zk-client -o jsonpath='{.data.redisHost}'"
	config, err := shell.Shellout(cmd)
	if err != nil {
		internal.DumpErrorAndPrintLocation(config)
		return err
	}

	//Check if configmap already exists
	cmd = "kubectl get configmap zk-redis-config -n pl -o jsonpath='{.data.redisHost}'"
	_, configExistsErr := shell.Shellout(cmd)
	if configExistsErr == nil {
		//config already exists, so returning without doing anything
		return nil
	}

	internal.DumpErrorAndPrintLocation("-----------" + config)

	cmd = fmt.Sprintf("kubectl create configmap zk-redis-config -n pl --from-literal=redisHost=%s", config)
	out, err = shell.Shellout(cmd)
	if err != nil {
		internal.DumpErrorAndPrintLocation(out)
		return err
	}

	return nil
}

func InstallDataStores(postgresVersion string, redisVersion string) error {
	// Install zk-client data stores
	inputToShellFile := ""
	inputToShellFile += " POSTGRES_VERSION=" + postgresVersion
	inputToShellFile += " REDIS_VERSION=" + redisVersion

	//Scenario where the project is being run locally.
	if viper.Get(internal.EmbedKeyFlag) == false {
		return shell.ExecuteShellFileWithSpinner(shell.GetPWD()+"/"+zkInstallStores, inputToShellFile, installDataStore, installDataStoreSuccessText, installDataStoreFailureText)
	} else {
		//Production scenario.
		return shell.ExecuteEmbeddedFileWithSpinner(internal.EmbeddedContent, zkInstallStores, inputToShellFile, installDataStore, installDataStoreSuccessText, installDataStoreFailureText)
	}
}

func InstallOlm() (err error) {
	if viper.Get(internal.EmbedKeyFlag) == false {
		return shell.ExecuteShellFileWithSpinner(shell.GetPWD()+"/"+olmInstall, " v0.25.0", installOlm, installOlmSuccessText, installOlmFailureText)
	} else {
		//Production scenario.
		return shell.ExecuteEmbeddedFileWithSpinner(internal.EmbeddedContent, olmInstall, " v0.25.0", installOlm, installOlmSuccessText, installOlmFailureText)
	}
}

func InstallPXOperator(ebpfMemory string) (err error) {

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
		ebpfArgs := " --deploy_olm=false"
		if ebpfMemory != "" {
			ebpfArgs = fmt.Sprintf(" --pem_memory_limit=%s", ebpfMemory)
		}
		cmd := utils.GetBackendCLIPath() + " deploy " + ebpfArgs
		out, err = shell.ShelloutWithSpinner(cmd, diSpinnerText, diSuccessText, diFailureText)

		if err != nil {
			//replace 'pl' with 'zk' in out
			out = strings.ReplaceAll(out, "pl", "zk")
			out = strings.ReplaceAll(out, "pixie", "ebpf")
			out = strings.ReplaceAll(out, "px", "zk")
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
	vizierYamlPath := shell.GetPWD() + "/" + cliVizierYaml
	if viper.Get(internal.EmbedKeyFlag) == true {
		yamlString := utils.GetEmbeddedFileContents(cliVizierYaml, internal.EmbeddedContent)
		vizierYamlPath = shell.GetPWD() + "/tmp_yaml_file.yaml"

		defer shell.DeleteFile(vizierYamlPath)

		err := utils.WriteTextToFile(yamlString, vizierYamlPath)
		if err != nil {
			return err
		}
	}
	patch := func() error {
		out, err := shell.Shellout("kubectl apply -f " + vizierYamlPath)
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

func ExtractZkHelmVersion() (*string, error) {
	passedVersionInterface := viper.Get(VersionKeyFlag)
	if passedVersionInterface != "" {
		passedVersion := passedVersionInterface.(string)
		if passedVersion != "" {
			return &passedVersion, nil
		}
	}
	url := "https://dl.zerok.ai/cli/helmversion.txt"

	// Make an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer response.Body.Close()
	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Convert the response body to a string
	bodyStr := string(body)

	// Find the index of "prod" in the string
	index := strings.Index(bodyStr, "prod/")
	if index == -1 {
		fmt.Println("String 'prod' not found in response")
		return nil, err
	}

	// Extract the version substring after "prod"
	version := bodyStr[index+len("prod/"):]

	// Trim any leading/trailing whitespaces
	version = strings.TrimSpace(version)

	return &version, nil
}

func ExtractZkEbpfHelmVersion() (*string, error) {
	passedVersionInterface := viper.Get(VersionKeyFlag)
	if passedVersionInterface != "" {
		passedVersion := passedVersionInterface.(string)
		if passedVersion != "" {
			return &passedVersion, nil
		}
	}
	url := "https://dl.zerok.ai/cli/ebpfversion.txt"

	// Make an HTTP GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return nil, err
	}
	defer response.Body.Close()
	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Convert the response body to a string
	bodyStr := string(body)

	// Find the index of "prod" in the string
	index := strings.Index(bodyStr, "prod/")
	if index == -1 {
		fmt.Println("String 'prod' not found in response")
		return nil, err
	}

	// Extract the version substring after "prod"
	version := bodyStr[index+len("prod/"):]

	// Trim any leading/trailing whitespaces
	version = strings.TrimSpace(version)

	return &version, nil
}

func InstallZKEbpf(clusterId, clusterName string, zkHelmVersion string) error {
	var inputToShellFile, shellFile string
	var zkEbpfJwtKey string = "9e21b08a50004a740e3ffc972c410e935da595ac570415f15a20870a0aebe68a8a0ca7791fdcaea76c8694a76b7451ae142b5f4f54673b4cb4f4407ef6dcf105"

	shellFile = zkInstallEbpf
	inputToShellFile = " ZK_CLUSTER_ID=" + clusterId +
		" ZK_CLUSTER_NAME=" + clusterName +
		" ZK_HELM_VERSION=" + zkHelmVersion +
		" ZK_EBPF_JWT_KEY=" + zkEbpfJwtKey

	_, err := GenerateCloudCertYAMLs("zk-client")
	if err != nil {
		fmt.Println("Error ", err)
		return err
	}

	//if viper.Get(internal.EmbedKeyFlag) == false {
	//	fmt.Println("@Debug04")
	return shell.ExecuteShellFileWithSpinner(shell.GetPWD()+"/"+shellFile, inputToShellFile, zkEbpfSpinnerText, zkEbpfSuccessText, zkEbpfFailureText)
	//} else {
	//	fmt.Println("@Debug05")
	//	return shell.ExecuteEmbeddedFileWithSpinner(internal.EmbeddedContent, shellFile, inputToShellFile, zkEbpfSpinnerText, zkEbpfSuccessText, zkEbpfFailureText)
	//}
}

func InstallZKServices(apiKey, clusterKey, clusterId, clusterName string, zkHelmVersion string) error {
	var inputToShellFile, shellFile string
	var helmVersion string = ""
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
			" ZK_OTLP_RECIEVER_VERSION=" + keyValueMap["zk-otlp-reciever"] +
			" ZK_PROMTAIL_VERSION=" + keyValueMap["zk-promtail"] +
			" PX_API_KEY=" + apiKey +
			" PX_CLUSTER_KEY=" + clusterKey +
			" ZK_CLUSTER_ID=" + clusterId +
			" PX_CLUSTER_ID=" + clusterName
	} else if viper.Get(internal.SpreadKeyFlag) == true {
		ui.GlobalWriter.Println("zerok spread mode is enabled")

		//get the versions
		versions := viper.Get(VersionKeyFlag).(string)
		keyValueMap := internal.GetKVPairsFromCSV(versions)

		shellFile = zkInstallSpreadClient
		inputToShellFile = " ZK_CLOUD_ADDR=" + zkCloudAddr +
			" ZK_SCENARIO_MANAGER_VERSION=" + keyValueMap["zk-scenario-manager"] +
			" ZK_AXON_VERSION=" + keyValueMap["zk-axon"] +
			" ZK_DAEMONSET_VERSION=" + keyValueMap["zk-daemonset"] +
			" ZK_GPT_VERSION=" + keyValueMap["zk-gpt"] +
			" ZK_WSP_CLIENT_VERSION=" + keyValueMap["zk-wsp-client"] +
			" ZK_OPERATOR_VERSION=" + keyValueMap["zk-operator"] +
			" ZK_OTLP_RECIEVER_VERSION=" + keyValueMap["zk-otlp-reciever"] +
			" ZK_PROMTAIL_VERSION=" + keyValueMap["zk-promtail"] +
			" PX_API_KEY=" + apiKey +
			" PX_CLUSTER_KEY=" + clusterKey +
			" ZK_CLUSTER_ID=" + clusterId +
			" PX_CLUSTER_ID=" + clusterName
	} else {
		ui.GlobalWriter.Println("zerok client mode is enabled")
		shellFile = zkInstallClient
		inputToShellFile = " ZK_CLOUD_ADDR=" + zkCloudAddr +
			" PX_API_KEY=" + apiKey +
			" PX_CLUSTER_KEY=" + clusterKey +
			" ZK_CLUSTER_ID=" + clusterId +
			" PX_CLUSTER_ID=" + clusterName +
			" ZK_HELM_VERSION=" + zkHelmVersion +
			" APP_NAME=zk-client"
		helmVersion = zkHelmVersion
	}
	if viper.Get(internal.GptKeyFlag) == true {
		inputToShellFile += " GPT_ENABLED=true"
	} else {
		inputToShellFile += " GPT_ENABLED=false"
	}

	ui.GlobalWriter.Println("zerok client configs:\n" +
		"    ZK_API_KEY=" + apiKey +
		"    ZK_HELM_VERSION=" + helmVersion)
	if inputToShellFile == "" {
		ui.GlobalWriter.Println("zerok client args are empty")
	} else {
		ui.GlobalWriter.Println("zerok client args are available")
	}

	if viper.Get(internal.EmbedKeyFlag) == false {
		return shell.ExecuteShellFileWithSpinner(shell.GetPWD()+"/"+shellFile, inputToShellFile, zksSpinnerText, zksSuccessText, zksFailureText)
	} else {
		return shell.ExecuteEmbeddedFileWithSpinner(internal.EmbeddedContent, shellFile, inputToShellFile, zksSpinnerText, zksSuccessText, zksFailureText)
	}
}
