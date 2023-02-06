package install

import (
	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
)

const (
	pxInstallBackendPostsetup string = "/zpx/scripts/postsetup-backend.sh"
)

var PixieBackendPostsetupCmd = &cobra.Command{
	Use:   "postsetup",
	Short: "Install PX Backend",
	RunE:  RunPxBackendPostsetupCmd,
}

func init() {
	ZKBackendCmd.AddCommand(PixieBackendPostsetupCmd)
}

func RunPxBackendPostsetupCmd(cmd *cobra.Command, args []string) error {
	_, err := logic.ExecWithLogsDurationAndSuccessM(logic.GetPWD()+pxInstallBackendPostsetup, "PX backend postsetup done")
	return err
}
