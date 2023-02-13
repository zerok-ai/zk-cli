package install

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/shell"
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

func RunZkOperatorCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return zkOperatorSetup(ctx)
}

func zkOperatorSetup(ctx context.Context) error {
	_, err := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+zkInstallOperator, "ZeroK operator installed successfully")
	return err
}
