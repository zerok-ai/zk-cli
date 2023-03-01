package cmd

import (
	"context"
	"os"
	"sync"
	"time"

	// "zkctl/cmd/pkg/utils"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
	// "zkctl/cmd/pkg/utils"
)

var (
	ZkOperatorCmd = &cobra.Command{
		Use:     "operator",
		Aliases: []string{"op", "optr"},
		Short:   "Install ZeroK operator",
		RunE:    runZkOperatorCmd,
		PreRunE: RunInstallPreCmd,
	}
	wg sync.WaitGroup
)

const (
	waitTimeForPodsInSeconds         = 150
	waitTimeForInstallationInSeconds = 240

	zkInstallOperator  string = "/operator/install-demo.sh"
	pxDownloadAuthJson string = "/zpx/scripts/setup-px-auth-json.sh"
	pxSetupDomain      string = "/zpx/scripts/setup-domain.sh"
	pxSetupIngress     string = "/zpx/scripts/setup-ingress.sh"

	pxCreateAuthJson string = "/zpx/scripts/pixie-ui-cli.sh -c authjson -d avinpx01.getanton.com"
)

func runZkOperatorCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return LogicZkOperatorSetup(ctx)
}

func LogicZkOperatorSetup(ctx context.Context) error {
	// var err error
	// _, err = shell.ExecWithLogsDurationAndSuccessM(shell.GetPWD()+zkInstallOperator, "ZeroK operator installed successfully")

	return installPXOperator(ctx)
}

func loginToPX(ctx context.Context) error {

	ui.GlobalWriter.PrintflnWithPrefixArrow("trying to authenticate")

	/**/
	authPath := utils.GetBackendAuthPath()
	os.Remove(authPath)
	auth, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+pxCreateAuthJson, "Credential file created for operator")

	if err != nil {
		return err
	}
	utils.WriteTextToFile(auth, authPath)
	// AUTHJSON=$($SCRIPTS_DIR/pixie-ui-cli.sh -c authjson)
	// echo $AUTHJSON >> ~/.pixie/auth.json

	/*/

	_, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+pxDownloadAuthJson, "Credential file created for operator")
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("==== %v", err)
		return err
	}
	/**/
	return nil
}

func installPXOperator(ctx context.Context) error {

	zerokcloudContext := "gke_zerok-dev_us-west1-b_avinpx01"
	var err error

	// deploy setup connection with zerok backend
	out, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+pxSetupDomain, "setting up the connection with zerok backend")
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===setup domain %v err=%s", err, out)
	}

	// set ingress to predeploy settings
	err = setIngress(zerokcloudContext, "predeploy")
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===predeploy %v", err)
		return err
	}

	// login to px
	err = loginToPX(ctx)
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===loginToPX %v", err)
		return err
	}

	wg.Add(1)

	go func() {

		go func() {
			// deploy px operator
			ui.GlobalWriter.PrintflnWithPrefixArrow("installing operator for managing data store")
			_, err = shell.ExecWithLogsDurationAndSuccessM(utils.GetBackendCLIPath()+" deploy", "Zerok daemon installed successfully")
			if err != nil {
				// send to sentry and print
				ui.GlobalWriter.Errorf("===deploy %v", err)
			}
		}()

		ui.GlobalWriter.PrintflnWithPrefixArrow("Waiting for %d sec for the installation", waitTimeForInstallationInSeconds)
		time.Sleep(waitTimeForInstallationInSeconds * time.Second)
		wg.Done()

	}()

	// set ingress to post deploy settings
	ui.GlobalWriter.PrintflnWithPrefixArrow("Waiting for %d sec for the pods to get deployed", waitTimeForPodsInSeconds)
	time.Sleep(waitTimeForPodsInSeconds * time.Second)
	err = setIngress(zerokcloudContext, "deploy")
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===deploy %v", err)
		return err
	}

	wg.Wait()

	ui.GlobalWriter.PrintSuccessMessageln("\ninstallation done. back to room temperature")

	return err
}

func setIngress(zerokcloudContext, command string) error {

	var err error

	// get current context
	curClusterContext, err := shell.Shellout("kubectl config current-context", false)
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===current-context %v", err)
	}

	// set current context to zerokcloud
	_, err = shell.Shellout("kubectl config use-context "+zerokcloudContext, false)
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===use-context %v", err)
	}
	// predeploy px operator
	_, err = shell.Shellout(shell.GetPWD()+pxSetupIngress+" "+command, false)
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===postdeploy %v", err)
	}
	// set current context to zerokcloud
	_, err = shell.Shellout("kubectl config use-context "+curClusterContext, false)
	if err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===use-context %v", err)
	}

	return err
}
