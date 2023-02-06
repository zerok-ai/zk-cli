package install

import (
	// "bytes"
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
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
	ctx := cmd.Context()
	return installPxOperator(ctx)
}

func installPxOperator(ctx context.Context) error {
	_, err := logic.ExecOnShellM(logic.GetPWD()+pxInstallOperator, "PX operator installed successfully")
	return err
}
