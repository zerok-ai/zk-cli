package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/k8"
)

var (
	labelActivateCmd = &cobra.Command{
		Use:   "activate",
		Short: "activate a namespace for zerok",
		RunE:  RunActivateNamespaceCmd,
	}
	labelDeactivateCmd = &cobra.Command{
		Use:   "deactivate",
		Short: "deactivate a namespace for zerok",
		RunE:  RunDeactivateNamespaceCmd,
	}
	rollingRestart  bool
	namespaceString string
)

const (
	ZKMARKER_LABEL_KEY   string = "zk-injection"
	ZKMARKER_LABEL_VALUE string = "enabled"
)

func init() {
	addToRootAndSetFlags(labelActivateCmd)
	addToRootAndSetFlags(labelDeactivateCmd)
}

func addToRootAndSetFlags(c *cobra.Command) {
	rootCmd.AddCommand(c)
	c.Flags().BoolVarP(&rollingRestart, "restart", "r", false, "Use this if you also want to do rolling restart post activation")
	c.Flags().StringVarP(&namespaceString, "namespace", "n", "default", "Source directory to read from")
}

func RunActivateNamespaceCmd(cmd *cobra.Command, args []string) error {

	zkNamespace := k8.NewZkNamespace(namespaceString)
	errorInMarking := zkNamespace.AddLabel(ZKMARKER_LABEL_KEY, ZKMARKER_LABEL_VALUE)
	if errorInMarking == nil && rollingRestart {
		return zkNamespace.DoRollingRestart(namespaceString)
	}

	return errorInMarking
}

func RunDeactivateNamespaceCmd(cmd *cobra.Command, args []string) error {

	zkNamespace := k8.NewZkNamespace(namespaceString)
	errorInMarking := zkNamespace.RemoveLabel(ZKMARKER_LABEL_KEY)
	if errorInMarking == nil && rollingRestart {
		return zkNamespace.DoRollingRestart(namespaceString)
	}

	return errorInMarking
}