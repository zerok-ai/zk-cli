package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"zkctl/cmd/internal"
	"zkctl/cmd/pkg/shell"
)

const (
	VizierTagFlag         = "VIZIER_TAG"
	VizierSetupKeyFlag    = "vz-setup"
	VizierSetupKeyEnvFlag = "VZ_SETUP"

	pxVizierDevModeSetup string = "/zpx/scripts/setup-vizier-export.sh"
	pxVizierYaml         string = "/zpx/scripts/modified/vizier/exported-vizier-pem.yaml"
	cliVizierYaml        string = "/vizier/vizier.yaml"
)

var (
	vizierDockerTag string
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Internal use",
	RunE:  RunDevCmd,
}

func init() {
	RootCmd.AddCommand(devCmd)

	// add flags
	internal.AddBoolFlag(RootCmd, VizierSetupKeyFlag, VizierSetupKeyEnvFlag, "", false, "for internal use only", true)
}

func createVizierYaml(vizierTag string) error {
	out, err := shell.Shellout("VIZIER_TAG="+vizierTag+" ./"+pxVizierDevModeSetup, false)
	if err != nil {
		internal.DumpErrorAndPrintLocation(fmt.Sprintf("vizier export failed, %s", out))
		return err
	}

	out, err = shell.Shellout("cp ./"+pxVizierYaml+" ./"+cliVizierYaml, false)
	if err != nil {
		internal.DumpErrorAndPrintLocation(fmt.Sprintf("vizier copy failed %s", out))
	}
	return err
}

func RunDevCmd(cmd *cobra.Command, args []string) error {

	vizierTagInterface := viper.Get(VizierTagFlag)
	if vizierTagInterface != nil {
		vizierDockerTag = vizierTagInterface.(string)
	} else {
		vizierDockerTag = "0.1.0-redishost"
	}

	if viper.Get(VizierSetupKeyFlag) == true {
		err := createVizierYaml(vizierDockerTag)
		if err != nil {
			return err
		}
	}

	return nil
}
