package install

import (
	// "bytes"
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxInstallOperator string = "/zpx/scripts/setup-operator.sh"
)

var PixieOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Install PX operator",
	RunE:  RunPxOperatorCmd,
}

func init() {
	PixieCmd.AddCommand(PixieOperatorCmd)
}

func RunPxOperatorCmd(cmd *cobra.Command, args []string) error {

	ui.GlobalWriter.PrintNoticeMessage("Installing PX operator")

	var err error
	ctx := cmd.Context()

	// Install px operator
	if err = installPxOperator(ctx); err != nil {
		return err
	}

	ui.GlobalWriter.PrintSuccessMessageln("PX operator installed successfully")
	return nil
}

func installPxOperator(ctx context.Context) error {
	logic.ExecOnShellsE(logic.GetPWD() + pxInstallOperator)
	return nil
}
