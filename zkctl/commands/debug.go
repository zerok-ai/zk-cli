package zkctl

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	zkctl "github.com/zerok-ai/zk-cli/zkctl/backend"
)

var onlyDigits bool
var inspectCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug a request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		i := args[0]
		res, kind := zkctl.Inspect(i, onlyDigits)

		pluralS := "s"
		if res == 1 {
			pluralS = ""
		}
		fmt.Printf("'%s' has %d %s%s.\n", i, res, kind, pluralS)
		fmt.Println("calling-birds:=", viper.GetString("calling-birds"))
	},
}

func init() {
	inspectCmd.Flags().BoolVarP(&onlyDigits, "digits", "d", false, "Count only digits")
	rootCmd.AddCommand(inspectCmd)
}
