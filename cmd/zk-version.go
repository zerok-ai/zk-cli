package cmd

import (
	"github.com/spf13/cobra"
	"zkctl/cmd/pkg/ui"
)

var (
	// this is a placeholder value which will be overriden by the build process
	BinaryVersion string
)

func init() {
	RootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get zerok cli version",
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.GlobalWriter.Println("version = " + BinaryVersion)
		return nil
	},
}
