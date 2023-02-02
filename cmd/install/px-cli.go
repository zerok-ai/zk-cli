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
	pxInstallCLI string = "/scripts/install.sh"
)

var PixieCLICmd = &cobra.Command{
	Use: "cli",
	// Aliases: []string{"op", "optr"},
	Short: "Install PX CLI",
	RunE: RunPxCLICmd,
}

func init() {
	PixieCmd.AddCommand(PixieCLICmd)
}

func RunPxCLICmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return checkAndInstallPxCLI(ctx);
}

func checkAndInstallPxCLI(ctx context.Context) error {
	return logic.ExecOnShellM(logic.GetPWD() + pxInstallCLI, "PX cli installed successfully")
}
