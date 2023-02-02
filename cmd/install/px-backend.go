package install

import (
	// "bytes"
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
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
	return checkAndInstallPxBackend(cmd.Context())
}

func checkAndInstallPxBackend(ctx context.Context) error {
	return logic.ExecOnShellM(logic.GetPWD() + pxInstallBackend, "PX backend installed successfully")
}
