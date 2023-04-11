package cmd

import (
	"context"
	"fmt"

	// "zkctl/cmd/pkg/k8s"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	delSpinnerText = "deleting zerok daemon and associated CRDs"
	delSuccessText = "zerok removed from the cluster successfully"
	delFailureText = "failed to delete zerok daemon"

	zkUninstallOperator string = "/operator/uninstall.sh"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete ZeroK",
		RunE:  RunDeleteCmd,
	}
)

func init() {
	RootCmd.AddCommand(deleteCmd)
}

func RunDeleteCmd(cmd *cobra.Command, args []string) error {
	err := logicZkDelete(cmd.Context())
	return err
}

func logicZkDelete(ctx context.Context) error {

	cmd := utils.GetBackendCLIPath() + " delete"

	out, err := shell.ShelloutWithSpinner(cmd, delSpinnerText, delSuccessText, delFailureText)

	if err != nil {
		filePath, _ := utils.DumpError(out)
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
		return err
	}

	// delete namespaces
	// return k8s.DeleteNamespaces([]string{"zerok-injector", "zerok-operator-system"})

	out, err = shell.ExecWithDurationAndSuccessM(shell.GetPWD()+zkUninstallOperator, "zeroK operator uninstalled successfully")
	if err != nil {
		filePath, _ := utils.DumpError(out)
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
		return err
	}

	//TODO: check all the namespaces and deactivate them

	return nil
}
