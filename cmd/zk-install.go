package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
	"zkctl/cmd/internal"
	"zkctl/cmd/internal/install"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	authAddress       string
	apiKey            string
	ebpfMemory        string
	postgresVersion   string
	redisVersion      string
	clusterName       string
	clusterKey        string
	clusterId         string
	zkHelmVersion     string
	zkEbpfHelmVersion string
)

type ContextKey struct {
	key string
}

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

	// add flags
	internal.AddStringFlag(RootCmd, install.EbpfMemoryFlag, install.EbpfMemoryEnvFlag, "", "", "ebpf memory. This can also be set through environment variable "+install.EbpfMemoryEnvFlag+" instead of passing the parameter", true)
	internal.AddStringFlag(RootCmd, install.ApiKeyFlag, install.ApiKeyEnvFlag, "", "", "api key. This can also be set through environment variable "+install.ApiKeyEnvFlag+" instead of passing the parameter", false)
	internal.AddStringFlag(RootCmd, install.VersionKeyFlag, install.VersionKeyEnvFlag, "", "", "version of the installation", false)
	internal.AddStringFlag(RootCmd, install.PostgresVersionKeyFlag, install.PostgresVersionKeyEnvFlag, "", "0.1.1-alpha", "postgres version to be installed", true)
	internal.AddStringFlag(RootCmd, install.RedisVersionKeyFlag, install.RedisVersionKeyEnvFlag, "", "0.1.1-alpha", "redis version to be installed", true)
}

func isHelmInstalled() bool {
	cmd := exec.Command("helm", "version", "--short")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false
	}

	// Convert output to string and check if it contains the version information
	outputStr := string(output)
	return strings.Contains(outputStr, "v")
}

func validateRequiredTools() error {
	ui.GlobalWriter.Println("Validating tools:")
	//Validate helm
	isHelmInstalled := isHelmInstalled()
	if !isHelmInstalled {
		return errors.New("helm is not installed. Please install helm and try again")
	} else {
		ui.GlobalWriter.PrintSuccessMessageln("helm is installed")
		ui.GlobalWriter.Println("")
	}
	return nil
}

func printVersionsAndOptions() {
	ui.GlobalWriter.Println("Versions:")
	ui.GlobalWriter.PrintflnWithPrefixBullet("zkctl version: %s", BinaryVersion)
	ui.GlobalWriter.PrintflnWithPrefixBullet("zerok services version: %s", zkHelmVersion)
	if viper.Get(internal.PxKeyFlag) == false {
		ui.GlobalWriter.PrintflnWithPrefixBullet("ebpf: %s", "false")
	}
	if ebpfMemory != "" {
		ui.GlobalWriter.PrintflnWithPrefixBullet("ebpf memory: %s", ebpfMemory)
	}
	if viper.Get(internal.ZksKeyFlag) == false {
		ui.GlobalWriter.PrintflnWithPrefixBullet("zerok services: %s", "false")
	}
	if viper.Get(internal.ZkStoresKeyFlag) == false {
		ui.GlobalWriter.PrintflnWithPrefixBullet("zerok stores: %s", "false")
	}
	if viper.Get(internal.OlmKeyFlag) == false {
		ui.GlobalWriter.PrintflnWithPrefixBullet("olm: %s", "false")
	}
	if viper.Get(internal.EbpfKeyFlag) == false {
		ui.GlobalWriter.PrintflnWithPrefixBullet("ebpf probes: %s", "false")
	}
	ui.GlobalWriter.Println("")
}

func RunInstallPreCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	err := LoadAndValidateFlags()
	if err != nil {
		return err
	}

	err = utils.InitializeFolders()
	if err != nil {
		return err
	}

	printVersionsAndOptions()

	err = validateRequiredTools()
	if err != nil {
		return err
	}

	clusterName, err = install.ValidateClusterAndTakeConsent(ctx)
	ui.GlobalWriter.PrintflnWithPrefixArrow("running pre-installation checks")
	if err != nil {
		ui.GlobalWriter.Printf("In PreInstallChecksAndTasks err %s \n", err)
		return err
	}

	// download and install px cli
	//return install.DownloadAndInstallPXCLI(ctx)
	return nil
}

func RunInstallCmd(cmd *cobra.Command, args []string) error {

	var err error

	// 1. login to px
	clusterMetadata := install.ClusterMetadata{
		ClientVersions: map[string]string{
			"cli":                 "",
			"zk-client":           "",
			"zk-axon":             "",
			"zk-scenario-manager": "",
			"zk-gpt":              "",
			"zk-operator":         "",
			"zk-daemonset":        "",
			"zk-wsp-client":       "",
		},
	}
	if clusterKey, clusterId, err = install.LoginToPX(authAddress, apiKey, clusterName, clusterMetadata); err != nil {
		return err
	}

	// 2. Install zk-client data stores
	if viper.Get(internal.ZkStoresKeyFlag) == true {
		err = install.InstallDataStores(postgresVersion, redisVersion)
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping stores installation")
	}

	// 5. Install zeroK services
	if viper.Get(internal.ZksKeyFlag) == true {
		err = install.InstallZKServices(apiKey, clusterKey, clusterId, clusterName, zkHelmVersion)
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping zk client services installation")
	}

	// 6. Install zeroK ebpf
	if viper.Get(internal.EbpfKeyFlag) == true {
		//clusterId = "e38ee038-32af-43a0-8e1b-510d31b0e8c9"
		//clusterName = "gke_zerok-dev_us-west1-b_devclient03"
		//zkHelmVersion = "0.0.1"
		err = install.InstallZKEbpf(clusterId, clusterName, zkEbpfHelmVersion)
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping zk ebpf installation")
	}

	// 7. print success message
	if err == nil {
		ui.GlobalWriter.PrintlnSuccessMessageln("installation done")
	}
	return err
}

var (
	// prodCloudAddress should be set through ldflags option in `go build` command
	prodCloudAddress string
)

func getCloudAddress() string {
	zkCloudAddr := viper.Get(install.ZkCloudAddressFlag)
	if zkCloudAddr == nil {
		zkCloudAddr = prodCloudAddress
		viper.Set(install.ZkCloudAddressFlag, zkCloudAddr)
	}
	return zkCloudAddr.(string)
}

func getPostgresVersion() string {
	postgresVersion := viper.Get(install.PostgresVersionKeyFlag)
	if postgresVersion == nil {
		postgresVersion = "0.1.1-alpha"
	}
	return postgresVersion.(string)
}

func getRedisVersion() string {
	redisVersion := viper.Get(install.RedisVersionKeyFlag)
	if redisVersion == nil {
		redisVersion = "0.1.1-alpha"
	}
	return redisVersion.(string)
}

func LoadAndValidateFlags() error {

	//Get API Key
	apiKey = install.GetAPIKey()
	if apiKey == "" {
		return errors.New(install.ApiKeyErrorMessage)
	}

	//Ebpf Memory
	ebpfMemory = install.GetEbpfMemory()

	//Stores' versions
	postgresVersion = getPostgresVersion()
	redisVersion = getRedisVersion()

	// set cloud address as an env variable for shell scripts
	zkCloudAddr := getCloudAddress()

	// set auth address
	plCloudAddr := "px." + zkCloudAddr
	os.Setenv("PL_CLOUD_ADDR", plCloudAddr)
	authAddress = fmt.Sprintf("https://api.%s", zkCloudAddr)

	zkHelmVersionLocal, err := install.ExtractZkHelmVersion()
	if err != nil {
		return err
	}
	zkHelmVersion = *zkHelmVersionLocal

	zkEbpfHelmVersionLocal, err := install.ExtractZkEbpfHelmVersion()
	if err != nil {
		return err
	}
	zkEbpfHelmVersion = *zkEbpfHelmVersionLocal

	return nil
}
