package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

var (
	labelNamespaceCmd = &cobra.Command{
		Use:   "activate",
		Short: "activate a namespace for zerok",
		RunE:  RunLabelNamespaceCmd,
		// Args:    cobra.ExactArgs(1),
	}
	rollingRestart bool
)

const (
	labelNamespaceCmdString string = "kubectl label namespace %s zk-injection=enabled"
	labelSuccessMessage     string = "Namespace %s is activated for ZeroK"

	restartCmdString      string = "kubectl rollout restart deployment -n %s"
	restartNeededMessage  string = "Rolling restart is needed for all the deployments in Namespace %s for activation to take affect"
	restartSuccessMessage string = "Rolling restart initiated for Namespace %s"
)

func init() {
	rootCmd.AddCommand(labelNamespaceCmd)

	labelNamespaceCmd.PersistentFlags().BoolVarP(&rollingRestart, "restart", "r", false, "Use this if you also want to do rolling restart")

}

func RunLabelNamespaceCmd(cmd *cobra.Command, args []string) error {

	var err error
	ctx := cmd.Context()

	if len(args) < 1 {
		ui.GlobalWriter.PrintErrorMessageln("Namespace missing in the command")
		return nil
	}

	namespace := args[0]

	if err = labelNamespaceAndRestart(ctx, namespace); err != nil {
		return err
	}
	return nil
}

func labelNamespaceAndRestart(ctx context.Context, namespace string) error {
	command := fmt.Sprintf(labelNamespaceCmdString, namespace)
	err := logic.ExecOnShellM(command, fmt.Sprintf(labelSuccessMessage, namespace))

	if (err !=nil){
		return err
	}

	if (rollingRestart) {
		command := fmt.Sprintf(restartCmdString, namespace)
		return logic.ExecOnShellM(command, fmt.Sprintf(restartSuccessMessage, namespace))
	} else {
		ui.GlobalWriter.PrintWarningMessageln(fmt.Sprintf(restartNeededMessage, namespace))
	}

	return nil
}
