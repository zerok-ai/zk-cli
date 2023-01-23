package install

import (
	// "bytes"
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxInstallBackendPresetup string = "/zpx/scripts/presetup-backend.sh"
)

var PixieBackendPresetupCmd = &cobra.Command{
	Use:   "presetup",
	Short: "Install PX Backend",
	RunE:  RunPxBackendPresetupCmd,
}

func init() {
	PixieBackendCmd.AddCommand(PixieBackendPresetupCmd)
}

func RunPxBackendPresetupCmd(cmd *cobra.Command, args []string) error {

	var err error
	ctx := cmd.Context()

	// Install px cli
	if err = RunPxCLICmd(cmd, args); err != nil {
		return err
	}

	ui.GlobalWriter.PrintNoticeMessage("Doing PX backend presetup")

	// Install px backend
	if err = pxBackendPresetup(ctx); err != nil {
		return err
	}

	ui.GlobalWriter.PrintSuccessMessageln("PX backend presetup done")

	return nil
}

func pxBackendPresetup(ctx context.Context) error {
	logic.ExecOnShellsE(logic.GetPWD() + pxInstallBackendPresetup)
	return nil
}
