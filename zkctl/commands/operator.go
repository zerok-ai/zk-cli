package zkctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	logic "github.com/zerok-ai/zk-cli/zkctl/backend"
)

var reverseCmd = &cobra.Command{
	Use:     "operator",
	Aliases: []string{"op", "optr"},
	Short:   "Reverses a string",
	// Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := logic.Reverse(args[0])
		fmt.Println(res)

		fmt.Println("cloudprovider:=", viper.GetString("cloudprovider"))
	},
}

func init() {
	rootCmd.AddCommand(reverseCmd)
}
