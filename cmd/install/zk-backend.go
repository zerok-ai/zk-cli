package install

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
	// // "github.com/spf13/viper"
)

var ZKBackendCmd = &cobra.Command{
	Use:   "backend",
	Short: "Install ZeroK components",
	RunE:  RunZKBackendCmd,
}

func init() {

}

func RunZKBackendCmd(cmd *cobra.Command, args []string) error {

	var err error

	// if err = RunZkOperatorCmd(cmd, args); err != nil {
	// 	return err
	// }

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

	// if err = RunPxBackendPostsetupCmd(cmd, args); err != nil {
	// 	return err
	// }

	return nil
}
