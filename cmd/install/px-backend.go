package install

import (
	// "bytes"
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxInstallBackend string = "/zpx/scripts/setup-backend.sh"
)

var PixieBackendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Install PX Backend",
	RunE:  RunPxBackendCmd,
}

func init() {
	PixieCmd.AddCommand(PixieBackendCmd)
}

func RunPxBackendCmd(cmd *cobra.Command, args []string) error {

	ui.GlobalWriter.PrintNoticeMessage("Installing PX backend")

	var err error
	ctx := cmd.Context()

	// Install px backend
	if err = checkAndInstallPxBackend(ctx); err != nil {
		return err
	}

	ui.GlobalWriter.PrintSuccessMessageln("PX backend installed successfully")

	return nil
}

func checkAndInstallPxBackend(ctx context.Context) error {
	logic.ExecOnShellsE(logic.GetPWD() + pxInstallBackend)
	return nil
}
