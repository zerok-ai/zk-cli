package install

import (
	"context"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
)

var ZkOperatorCmd = &cobra.Command{
	Use:     "operator",
	Aliases: []string{"op", "optr"},
	Short:   "Install ZeroK operator",
	RunE:    RunZKOperatorCmd,
}

const (
	zkInstallOperator string = "/operator/install-demo.sh"
)

func init() {

}

func RunZKOperatorCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return zkOperatorSetup(ctx)
}

func zkOperatorSetup(ctx context.Context) error {
	return logic.ExecOnShellM(logic.GetPWD() + zkInstallOperator, "ZeroK operator installed successfully")
}
