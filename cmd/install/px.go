package install

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	// // "github.com/spf13/viper"
)

var PixieCmd = &cobra.Command{
	Use:   "px",
	Short: "Install PX",
	RunE:  RunPxCmd,
}

func init() {

}

func RunPxCmd(cmd *cobra.Command, args []string) error {

	var err error

	if err = RunPxBackendPresetupCmd(cmd, args); err != nil {
		return err
	}

	if err = RunPxBackendCmd(cmd, args); err != nil {
		return err
	}

	for x := 0; x < 5; x++ {
		time.Sleep(30 * time.Second)

		// check whether pods are up
		fmt.Println("Checking whether the pods of backend are up")
	}

	if err = RunPxOperatorCmd(cmd, args); err != nil {
		return err
	}

	// if err = RunPxBackendPostsetupCmd(cmd, args); err != nil {
	// 	return err
	// }

	return nil
}
