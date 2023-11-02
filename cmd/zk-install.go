package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"zkctl/cmd/internal"
	"zkctl/cmd/internal/install"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	authAddress   string
	apiKey        string
	ebpfMemory    string
	clusterName   string
	clusterKey    string
	zkHelmVersion string
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
}

func PrintVersionsAndOptions() {
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

	PrintVersionsAndOptions()

	clusterName, err = install.ValidateClusterAndTakeConsent(ctx)
	ui.GlobalWriter.PrintflnWithPrefixArrow("running pre-installation checks")
	if err != nil {
		ui.GlobalWriter.Printf("In PreInstallChecksAndTasks err %s \n", err)
		return err
	}

	// download and install px cli
	return install.DownloadAndInstallPXCLI(ctx)
}

func RunInstallCmd(cmd *cobra.Command, args []string) error {

	var err error

	// 1. login to px
	if clusterKey, err = install.LoginToPX(authAddress, apiKey, clusterName); err != nil {
		return err
	}

	// 2. Install zk-client data stores
	if viper.Get(internal.ZkStoresKeyFlag) == true {
		err = install.InstallDataStores()
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping stores installation")
	}

	// 3 Install olm
	if viper.Get(internal.OlmKeyFlag) == true {
		err = install.InstallOlm()
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping olm installation")
	}

	// 4. Install default pixie - pl
	if viper.Get(internal.PxKeyFlag) == true {
		err = install.InstallPXOperator(ebpfMemory)
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping ebpf installation")
	}

	// 5. Install zeroK services
	if viper.Get(internal.ZksKeyFlag) == true {
		err = install.InstallZKServices(apiKey, clusterKey, clusterName, zkHelmVersion)
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping zk client services installation")
	}

	// 6. install the kustomization for vizier over the default code -- doing it later as it needs the redis instance to come up
	if viper.Get(internal.EbpfKeyFlag) == true {
		err = install.InstallVizier()
		if err != nil {
			return err
		}
	} else {
		ui.GlobalWriter.Println("Skipping ebpf probe installation")
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

func LoadAndValidateFlags() error {

	//Get API Key
	apiKey = install.GetAPIKey()
	if apiKey == "" {
		return errors.New(install.ApiKeyErrorMessage)
	}

	//Ebpf Memory
	ebpfMemory = install.GetEbpfMemory()

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

	return nil
}
