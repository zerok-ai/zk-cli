package cmd

import (
	"context"
	"fmt"
	"os"

	// "sync"
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

	zkInstallOperator  string = "/operator/install-demo.sh"
	pxDownloadAuthJson string = "/zpx/scripts/setup-px-auth-json.sh"
	pxSetupDomain      string = "/zpx/scripts/setup-domain.sh"
	pxSetupIngress     string = "/zpx/scripts/setup-ingress.sh"

	pxCreateAuthJson string = "/zpx/scripts/pixie-ui-cli.sh -c authjson -d"

	diSpinnerText = "installing zerok daemon and associated CRDs"
	diSuccessText = "zerok daemon installed successfully"
	diFailureText = "failed to install zerok daemon"
)

func runZkOperatorCmd(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	return LogicZkOperatorSetup(ctx)
}

func LogicZkOperatorSetup(ctx context.Context) error {
	// var err error
	_, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+zkInstallOperator, "zeroK operator installed successfully")
	if err != nil {
		return err
	}

	return installPXOperator(ctx)
}

func loginToPX(ctx context.Context) error {

	ui.GlobalWriter.PrintflnWithPrefixArrow("trying to authenticate")

	/**/
	authPath := utils.GetBackendAuthPath()
	os.Remove(authPath)

	cmd := fmt.Sprintf("%s%s %s", shell.GetPWD(), pxCreateAuthJson, viper.Get("PL_CLOUD_ADDR"))

	auth, err := shell.ExecWithDurationAndSuccessM(cmd, "Credential file created for operator")
	if err == nil {
		return utils.WriteTextToFile(auth, authPath)
	}
	return err
}

func installPXOperator(ctx context.Context) (err error) {

	zerokcloudContext := viper.Get("PLC_CLUSTER").(string)

	// deploy setup connection with zerok backend
	if out, err := shell.ExecWithDurationAndSuccessM(shell.GetPWD()+pxSetupDomain, "setting up the connection with zerok backend"); err != nil {
		// send to sentry and print
		ui.GlobalWriter.Errorf("===setup domain %v err=%s", err, out)
	}

	// set ingress to predeploy settings
	if err = setIngress(zerokcloudContext, "predeploy"); err != nil {
		// send to sentry and print
		return ui.GlobalWriter.Errorf("===predeploy %v", err)
	}

	// login to px
	if err = loginToPX(ctx); err != nil {
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

	if err = setIngress(zerokcloudContext, "deploy"); err != nil {
		// send to sentry and print
		return ui.GlobalWriter.Errorf("===deploy %v", err)
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

func setIngress(zerokcloudContext, command string) (err error) {

	// get current context
	curClusterContext, err := shell.Shellout("kubectl config current-context", false)
	if err != nil {
		// send to sentry and print
		return ui.GlobalWriter.Errorf("===current-context %v", err)
	}

	defer func() {
		// set current context to zerokcloud
		if out, err1 := shell.Shellout("kubectl config use-context "+curClusterContext, false); err1 != nil {
			// send to sentry and print
			if err == nil {
				err = ui.GlobalWriter.Errorf("===use-context %v %s", err1, out)
			}
		}
	}()

	// set current context to zerokcloud
	if out, err := shell.Shellout("kubectl config use-context "+zerokcloudContext, false); err != nil {
		// send to sentry and print
		return ui.GlobalWriter.Errorf("===use-context %v %s", err, out)
	}

	// px operator predeploy/deploy
	if out, err := shell.Shellout(shell.GetPWD()+pxSetupIngress+" "+command, false); err != nil {
		// send to sentry and print
		return ui.GlobalWriter.Errorf("===command %v \n%s", err, out)
	}

	return nil
}
