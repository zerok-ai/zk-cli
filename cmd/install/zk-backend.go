package install

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/shell"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxInstallCLI              string = "/scripts/install-px-cli.sh"
	pxInstallOperator         string = "/zpx/scripts/setup-operator.sh"
	pxInstallBackendPresetup  string = "/zpx/scripts/presetup-backend.sh"
	pxInstallBackend          string = "/zpx/scripts/setup-backend.sh"
	pxInstallBackendPostsetup string = "/zpx/scripts/postsetup-backend.sh"
)

func init() {
	ZKBackendCmd.AddCommand(ZkBackendCLICmd)
	ZKBackendCmd.AddCommand(ZkBackendOperatorCmd)
	ZKBackendCmd.AddCommand(ZkBackendPostsetupCmd)
	ZKBackendCmd.AddCommand(ZKDatastore)
}

var ZkBackendPostsetupCmd = &cobra.Command{
	Use:   "postsetup",
	Short: "Run backend post setup tasks",
	RunE:  RunZkBackendPostsetupCmd,
}

func RunZkBackendPostsetupCmd(cmd *cobra.Command, args []string) error {
	_, err := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallBackendPostsetup, "PX backend postsetup done")
	return err
}

var ZkBackendCLICmd = &cobra.Command{
	Use:   "cli",
	Short: "Install PX CLI",
	RunE:  RunBackendCLICmd,
}

func RunBackendCLICmd(cmd *cobra.Command, args []string) error {
	_, err := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallCLI, "PX cli installed successfully")
	return err
}

var ZkBackendOperatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "Install PX operator",
	RunE:  RunPxOperatorCmd,
}

func RunPxOperatorCmd(cmd *cobra.Command, args []string) error {
	_, err := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallOperator, "PX operator installed successfully")
	return err
}

var ZKDatastore = &cobra.Command{
	Use:   "datastore",
	Short: "Install zerok datastore",
	RunE:  RunZkDatastoreCmd,
}

func RunZkDatastoreCmd(cmd *cobra.Command, args []string) error {
	ui.GlobalWriter.PrintNoticeMessage("Running pre-setup scripts ")
	_, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+pxInstallBackendPresetup, "PX backend presetup done")
	if err != nil {
		return err
	}

	_, err1 := shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+pxInstallBackend, "PX backend installed successfully")
	return err1
}

var ZKBackendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Install ZeroK components",
	RunE:  RunZKBackendCmd,
}

func RunZKBackendCmd(cmd *cobra.Command, args []string) error {

	var err error

	if err := RunBackendCLICmd(cmd, args); err != nil {
		return err
	}

	if err = RunZkDatastoreCmd(cmd, args); err != nil {
		return err
	}

	for x := 0; x < 5; x++ {
		time.Sleep(30 * time.Second)

		// check whether pods are up
		fmt.Println("Checking readiness of the backend pods")
	}
	ui.GlobalWriter.PrintSuccessMessageln("Pods ready")

	if err = RunPxOperatorCmd(cmd, args); err != nil {
		return err
	}

	// if err = RunZkBackendPostsetupCmd(cmd, args); err != nil {
	// 	return err
	// }

	return nil
}
