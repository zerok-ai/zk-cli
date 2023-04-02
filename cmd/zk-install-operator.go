package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"time"

	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ZkOperatorCmd = &cobra.Command{
		Use:     "operator",
		Aliases: []string{"op", "optr"},
		Short:   "Install ZeroK operator",
		RunE:    runZkOperatorCmd,
		PreRunE: RunInstallPreCmd,
	}
	// wg sync.WaitGroup
)

const (
	waitTimeForPodsInSeconds         = 100
	waitTimeForService               = 60
	waitTimeForInstallationInSeconds = 240

	zkInstallOperator  string = "/operator/buildAndInstall.sh"
	pxDownloadAuthJson string = "/zpx/scripts/setup-px-auth-json.sh"
	pxSetupDomain      string = "/zpx/scripts/setup-domain.sh"
	pxSetupIngress     string = "/zpx/scripts/setup-ingress.sh"

	pxCreateAuthJson string = "/zpx/scripts/pixie-ui-cli.sh -c authjson -d"

	diSpinnerText = "installing zerok daemon and associated CRDs"
	diSuccessText = "zerok daemon installed successfully"
	diFailureText = "failed to install zerok daemon"

	//TODO: Change it to prod url later
	ZKCLOUD_ADDRESS = "http://auth.avinpx01.getanton.com:80"
)

func runZkOperatorCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return LogicZkOperatorSetup(ctx)
}

func LogicZkOperatorSetup(ctx context.Context) error {

	_, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+zkInstallOperator, "zeroK operator installed successfully")
	if err != nil {
		ui.LogAndPrintError(fmt.Errorf("failed to install zkoperator: %v", err))
		return err
	}

	return installPXOperator(ctx)
}

func loginToPX(ctx context.Context, apiKey string) error {

	ui.GlobalWriter.PrintflnWithPrefixArrow("trying to authenticate")

	path := "v1/p/auth/login"
	urlToBeCalled := fmt.Sprintf("%s/%s?apikey=%s", ZKCLOUD_ADDRESS, path, apiKey)
	shellOut, err := shell.Shellout("curl -lsS '" + urlToBeCalled + "'", false)
	if (err != nil){
		return err
	}
	jsonResponse := map[string]any{}
	json.Unmarshal([]byte(shellOut), &jsonResponse)
	if (int(jsonResponse["status"].(float64)) != 200){
		//TODO: A better error?
		return errors.New("Something went wrong")
	}

	authPath := utils.GetBackendAuthPath()
	os.Remove(authPath)
	
	authTokenPayload := jsonResponse["payload"]
	authTokenPayloadBytes, err := json.Marshal(authTokenPayload)
	if err == nil {
		return utils.WriteTextToFile(string(authTokenPayloadBytes), authPath)
	}
	return err
}

func installPXOperator(ctx context.Context) (err error) {

	apiKey := viper.Get("apikey").(string)

	// login to px
	if err = loginToPX(ctx, apiKey); err != nil {
		// send to sentry and print
		return ui.GlobalWriter.Errorf("===loginToPX %v", err)
	}

	// wg.Add(1)

	go func() {

		// defer wg.Done()
		// go func() error {
		// deploy px operator
		ui.GlobalWriter.PrintflnWithPrefixArrow("installing operator for managing data store")

		var out string
		/*/
		cmd := utils.GetBackendCLIPath() + " deploy"
		out, err := shell.ExecWithLogsDurationAndSuccessM(cmd, "Zerok daemon installed successfully")
		/*/
		cmd := utils.GetBackendCLIPath() + " deploy"
		out, err = shell.ShelloutWithSpinner(cmd, diSpinnerText, diSuccessText, diFailureText)
		/**/

		filePath, _ := utils.DumpError(out)
		if err != nil {
			// send to sentry and print
			ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
			// wg.Done()
		}
		// return nil
		// }()

		// ui.GlobalWriter.PrintflnWithPrefixArrow("Waiting for %d sec for the installation", waitTimeForInstallationInSeconds)
		// time.Sleep(waitTimeForInstallationInSeconds * time.Second)
	}()

	if err != nil {
		return err
	}

	// set ingress to post deploy settings
	// ui.GlobalWriter.PrintflnWithPrefixArrow("Waiting for %d sec for the pods to get deployed", waitTimeForPodsInSeconds)

	for i := 0; i <= waitTimeForPodsInSeconds; i++ {
		time.Sleep(1 * time.Second)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	for i := 0; i <= waitTimeForService; i++ {
		time.Sleep(1 * time.Second)
		if err != nil {
			return err
		}
	}

	return nil
}
