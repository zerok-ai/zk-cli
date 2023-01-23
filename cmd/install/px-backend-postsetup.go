package install

import (
	// "bytes"
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxInstallBackendPostsetup string = "/zpx/scripts/postsetup-backend.sh"
)

var PixieBackendPostsetupCmd = &cobra.Command{
	Use:   "postsetup",
	Short: "Install PX Backend",
	RunE:  RunPxBackendPostsetupCmd,
}

func init() {
	PixieBackendCmd.AddCommand(PixieBackendPostsetupCmd)
}

func RunPxBackendPostsetupCmd(cmd *cobra.Command, args []string) error {

	var err error
	ctx := cmd.Context()

	ui.GlobalWriter.PrintNoticeMessage("Doing PX backend postsetup")

	// Install px backend
	if err = pxBackendPostsetup(ctx); err != nil {
		return err
	}

	ui.GlobalWriter.PrintSuccessMessageln("PX backend postsetup done")

	return nil
}

func pxBackendPostsetup(ctx context.Context) error {
	logic.ExecOnShellsE(logic.GetPWD() + pxInstallBackendPostsetup)
	return nil
}
