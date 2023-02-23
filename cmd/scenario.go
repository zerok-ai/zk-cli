package cmd

import (
	// "bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	client "github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/scenarios"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/shell"
	"github.com/zerok-ai/zk-cli/zkctl/cmd/pkg/ui"
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
	RootCmd.AddCommand(PixieClientCmd)
}

func RunPxClientCmd(cmd *cobra.Command, args []string) error {

	ui.GlobalWriter.Println("--> Authenticating user")
	cloudAddress, err := shell.Shellout(fmt.Sprintf("%s/%s -c domain", shell.GetPWD(), pxConstants), false)
	if err != nil {
		return err
	}

	clusterId, err := shell.Shellout(fmt.Sprintf("%s/%s -c cluster", shell.GetPWD(), pxConstants), false)
	if err != nil {
		return err
	}

	ui.GlobalWriter.PrintSuccessMessageln("User authenticated")
	ui.GlobalWriter.Println("")

	apikey, err := shell.Shellout(fmt.Sprintf("%s/%s -c apikey", shell.GetPWD(), pxConstants), false)
	if err != nil {
		return err
	}

	ui.GlobalWriter.Println("--> Fetching data ")

	// ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf("cloudAddress=%s clusterId=%s apikey=%s\n", cloudAddress, clusterId, apikey))

	// cloudAddress := "pxtest1.getanton.com:443"
	// clusterId := "a95919ea-a65d-4349-907f-783193558fc0"
	// apikey := "px-api-d16c812a-4b44-4aa6-a79f-70825ef85003"

	pxlPath := fmt.Sprintf("%s/zk-api-server/getNamespaceHTTPTraffic.pxl", shell.GetPWD())
	return client.PixieClient(strings.TrimSpace(apikey), strings.TrimSpace(clusterId), strings.TrimSpace(cloudAddress)+":443", pxlPath)
}
