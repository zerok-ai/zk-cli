package install

import (
	// "bytes"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
)

const (
	pxInstallOperator string = "/zpx/scripts/setup-operator.sh"
)

var ZkBackendOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Install PX operator",
	RunE:  RunPxOperatorCmd,
}

func init() {
	ZKBackendCmd.AddCommand(ZkBackendOperatorCmd)
}

func RunPxOperatorCmd(cmd *cobra.Command, args []string) error {
	_, err := logic.ExecWithLogsDurationAndSuccessM(logic.GetPWD()+pxInstallOperator, "PX operator installed successfully")
	return err
}
