package cmd

import (
	"context"
	"fmt"

	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
	
)

const (
	delSpinnerText = "deleting zerok daemon and associated CRDs"
	delSuccessText = "zerok removed from the cluster successfully"
	delFailureText = "failed to delete zerok daemon"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete ZeroK",
	RunE:  RunDeleteCmd,
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

func RunDeleteCmd(cmd *cobra.Command, args []string) error {
	err := logicZkDelete(cmd.Context())
	return err
}

func logicZkDelete(ctx context.Context) error {

	var out string

	cmd := utils.GetBackendCLIPath() + " delete 2>&1 > ./dump"

	out, err := shell.ShelloutWithSpinner(cmd, delSpinnerText, delSuccessText, delFailureText)

	filePath, _ := utils.DumpError(out)
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
		// wg.Done()
		return err
	}

	return err
}
