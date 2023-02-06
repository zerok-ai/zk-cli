package install

import (
	// "bytes"
	"context"
	// "os"

	// // "errors"
	// "fmt"
	// "log"
	// "os/exec"
	// "time"

	"github.com/spf13/cobra"
	// // "github.com/spf13/viper"
	// // "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/k8s"
	// install "github.com/zerok-ai/zk-cli/zkctl/cmd"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
)

const (
	pxInstallCLI string = "/scripts/install-px-cli.sh"
)

var ZkBackendCLICmd = &cobra.Command{
	Use: "cli",
	// Aliases: []string{"op", "optr"},
	Short: "Install PX CLI",
	RunE:  RunBackendCLICmd,
}

func init() {
	ZKBackendCmd.AddCommand(ZkBackendCLICmd)
}

func RunBackendCLICmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return checkAndInstallPxCLI(ctx)
}

func checkAndInstallPxCLI(ctx context.Context) error {
	_, err := logic.ExecWithLogsDurationAndSuccessM(logic.GetPWD()+pxInstallCLI, "PX cli installed successfully")
	return err
}
