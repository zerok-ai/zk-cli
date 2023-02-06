package install

import (
	// "bytes"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxInstallBackendPresetup string = "/zpx/scripts/presetup-backend.sh"
	pxInstallBackend         string = "/zpx/scripts/setup-backend.sh"
)

var ZKDatastore = &cobra.Command{
	Use:   "datastore",
	Short: "Install zerok datastore",
	RunE:  RunZkDatastoreCmd,
}

func init() {
	ZKBackendCmd.AddCommand(ZKDatastore)
}

func RunZkDatastoreCmd(cmd *cobra.Command, args []string) error {
	ui.GlobalWriter.PrintNoticeMessage("Running pre-setup scripts ")
	_, err := logic.ExecWithDurationAndSuccessM(logic.GetPWD()+pxInstallBackendPresetup, "PX backend presetup done")
	if err != nil {
		return err
	}

	_, err1 := logic.ExecWithLogsDurationAndSuccessM(logic.GetPWD()+pxInstallBackend, "PX backend installed successfully")
	return err1
}
