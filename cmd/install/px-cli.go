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
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxInstallCLI    string = "/scripts/install.sh"
)

var PixieCLICmd = &cobra.Command{
	Use: "cli",
	// Aliases: []string{"op", "optr"},
	Short: "Install PX CLI",
	// Args:    cobra.ExactArgs(1),
	RunE: RunPxCLICmd,
}

func init() {
	PixieCmd.AddCommand(PixieCLICmd)
}	

func RunPxCLICmd(cmd *cobra.Command, args []string) error {

	ui.GlobalWriter.PrintNoticeMessage("Installing PX cli")

	var err error
	ctx := cmd.Context()

	// Install px cli
	if err = checkAndInstallPxCLI(ctx); err != nil {
		return err
	}

	ui.GlobalWriter.PrintSuccessMessageln("PX cli installed successfully")

	return nil
}

func checkAndInstallPxCLI(ctx context.Context) error {
	logic.ExecOnShellsE(logic.GetPWD()+pxInstallCLI)
	return nil
}