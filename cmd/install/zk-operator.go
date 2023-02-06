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
	RunE:    RunZkOperatorCmd,
}

const (
	zkInstallOperator string = "/operator/install-demo.sh"
)

func init() {

}

func RunZkOperatorCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return zkOperatorSetup(ctx)
}

func zkOperatorSetup(ctx context.Context) error {
	_, err := logic.ExecWithLogsDurationAndSuccessM(logic.GetPWD()+zkInstallOperator, "ZeroK operator installed successfully")
	return err
}
