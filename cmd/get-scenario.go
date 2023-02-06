package cmd

import (
	// "bytes"
	"strings"
	client "zk-api-server"

	// "os"

	// // "errors"
	"fmt"
	// "log"
	// "os/exec"
	// "time"

	"github.com/spf13/cobra"
	// // "github.com/spf13/viper"
	// // "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/k8s"
	// install "github.com/zerok-ai/zk-cli/zkctl/cmd"
	logic "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
	// "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
)

const (
	pxConstants string = "zpx/scripts/pixie-ui-cli.sh"
)

var PixieClientCmd = &cobra.Command{
	Use: "datadump",
	// Aliases: []string{"op", "optr"},
	Short: "Interact with PX server through this client",
	RunE:  RunPxClientCmd,
}

func init() {
	rootCmd.AddCommand(PixieClientCmd)
}

func RunPxClientCmd(cmd *cobra.Command, args []string) error {

	ui.GlobalWriter.Println("--> Authenticating user")
	cloudAddress, err := logic.Shellout(fmt.Sprintf("%s/%s -c domain", logic.GetPWD(), pxConstants), false)
	if err != nil {
		return err
	}

	clusterId, err := logic.Shellout(fmt.Sprintf("%s/%s -c cluster", logic.GetPWD(), pxConstants), false)
	if err != nil {
		return err
	}

	ui.GlobalWriter.PrintSuccessMessageln("User authenticated")
	ui.GlobalWriter.Println("")

	apikey, err := logic.Shellout(fmt.Sprintf("%s/%s -c apikey", logic.GetPWD(), pxConstants), false)
	if err != nil {
		return err
	}

	ui.GlobalWriter.Println("--> Fetching data ")

	// ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf("cloudAddress=%s clusterId=%s apikey=%s\n", cloudAddress, clusterId, apikey))

	// cloudAddress := "pxtest1.getanton.com:443"
	// clusterId := "a95919ea-a65d-4349-907f-783193558fc0"
	// apikey := "px-api-d16c812a-4b44-4aa6-a79f-70825ef85003"

	pxlPath := fmt.Sprintf("%s/zk-api-server/getNamespaceHTTPTraffic.pxl", logic.GetPWD())
	return client.PixieClient(strings.TrimSpace(apikey), strings.TrimSpace(clusterId), strings.TrimSpace(cloudAddress)+":443", pxlPath)
}
