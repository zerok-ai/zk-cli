package install

import (
	// "bytes"
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
)

const (
	pxInstallBackendPresetup string = "/zpx/scripts/presetup-backend.sh"
)

var PixieBackendPresetupCmd = &cobra.Command{
	Use:   "presetup",
	Short: "Tasks to run before setup of PX Backend",
	RunE:  RunPxBackendPresetupCmd,
}

func init() {
	PixieBackendCmd.AddCommand(PixieBackendPresetupCmd)
}

func RunPxBackendPresetupCmd(cmd *cobra.Command, args []string) error {

	if err := RunPxCLICmd(cmd, args); err != nil {
		return err
	}

	return pxBackendPresetup(cmd.Context())
}

func pxBackendPresetup(ctx context.Context) error {
	_, err := logic.ExecOnShellM(logic.GetPWD()+pxInstallBackendPresetup, "PX backend presetup done")
	return err
}
