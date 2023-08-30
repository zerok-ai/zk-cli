package cmd

import (
	"errors"
	"fmt"
	"os"
	"zkctl/cmd/internal"

	"zkctl/cmd/internal/install"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	VizierTagFlag = "VIZIER_TAG"

	CONTEXT                      = "ctx"
	SHOULD_PRINT_SUCCESS_MESSAGE = "should_print_success_msg"

	zkInstallStores   string = "/db-helm-charts/install-db.sh"
	zkUnInstallClient string = "/helm-charts/uninstall.sh"
)

var (
	authAddress string
	apiKey      string
	clusterName string
	clusterKey  string
	vizierTag   string
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
	internal.AddBoolFlag(RootCmd, install.DevKeyFlag, install.DevKeyEnvFlag, "z", false, "for internal use only", true)
	internal.AddBoolFlag(RootCmd, install.VizierSetupKeyFlag, install.VizierSetupKeyEnvFlag, "", false, "for internal use only", true)
	internal.AddStringFlag(RootCmd, install.ApiKeyFlag, install.ApiKeyEnvFlag, "", "", "api key. This can also be set through environment variable "+install.ApiKeyEnvFlag+" instead of passing the parameter", false)
	internal.AddStringFlag(RootCmd, install.VersionKeyFlag, install.VersionKeyEnvFlag, "", "", "version of the installation", false)
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

	ctx := cmd.Context()
	var err error

	// 1. login to px
	if clusterKey, err = install.LoginToPX(authAddress, apiKey, clusterName); err != nil {
		// send to sentry and print
		return err
	}

	// 2. Install zk-client data stores
	err = internal.ExecuteShellFile(shell.GetPWD()+zkInstallStores, " APP_NAME=zk-stores ",
		"failed to install zk_stores", "zk_stores installed successfully")
	if err != nil {
		return err
	}

	// 3. Install default pixie - pl
	err = install.InstallPXOperator(ctx, apiKey)
	if err != nil {
		return err
	}

	// 4. Install zeroK services
	err = install.InstallZKServices(apiKey, clusterKey)
	if err != nil {
		return err
	}

	// 5. install the kustomization for vizier over the default code -- doing it later as it needs the redis instance to come up
	err = install.InstallVizier(vizierTag)
	if err != nil {
		return err
	}

	// 6. print success message
	if err == nil {
		ui.GlobalWriter.PrintlnSuccessMessageln("installation done")
	}
	return err
}

func LoadAndValidateFlags() error {

	//Get API Key
	apiKey = install.GetAPIKey()
	if apiKey == "" {
		return errors.New(install.ApiKeyErrorMessage)
	}

	// set cloud address as an env variable for shell scripts
	zkCloudAddr := install.GetCloudAddress()

	// set vizier tag
	vizierTagInterface := viper.Get(VizierTagFlag)
	if vizierTagInterface != nil {
		vizierTag = vizierTagInterface.(string)
	}

	// set auth address
	plCloudAddr := "px." + zkCloudAddr
	os.Setenv("PL_CLOUD_ADDR", plCloudAddr)
	authAddress = fmt.Sprintf("https://api.%s", zkCloudAddr)

	return nil
}
