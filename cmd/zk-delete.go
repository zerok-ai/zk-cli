package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"zkctl/cmd/pkg/k8s"
	"zkctl/cmd/pkg/shell"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	delSpinnerText = "deleting zerok daemon and associated CRDs"
	delSuccessText = "zerok removed from the cluster successfully"
	delFailureText = "failed to delete zerok daemon"

	zkUninstallOperator string = "/operator/uninstall.sh"
)

var (
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete ZeroK",
		RunE:  RunDeleteCmd,
	}
)

func init() {
	RootCmd.AddCommand(deleteCmd)
}

func RunDeleteCmd(cmd *cobra.Command, args []string) error {
	err := logicZkDelete(cmd.Context())
	return err
}

func logicZkDelete(ctx context.Context) error {

	// Step1: Uninstalling backend cli daemon
	cmd := utils.GetBackendCLIPath() + " delete"

	out, cmdErr := shell.ShelloutWithSpinner(cmd, delSpinnerText, delSuccessText, delFailureText)

	if cmdErr != nil {
		filePath, _ := utils.DumpError(out)
		ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
		return cmdErr
	}

	// Step2: Adding finalizer to zk-client namespace so that zk-operator can properly shut itself down
	kubeconfig := viper.GetString(KUBECONFIG_FLAG)
	kubecontext := viper.GetString(KUBECONTEXT_FLAG)

	kubeClient, err := k8s.NewKubeClient(kubeconfig, kubecontext)
	if err != nil {
		return err
	}
	namespace, err := kubeClient.CoreV1().Namespaces().Get(context.TODO(), "zk-client", metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		finalizer := "operator/cleanup-pods"
		finalizerAlreadyPresent := false
		for _, f := range namespace.ObjectMeta.Finalizers {
			if f == finalizer {
				finalizerAlreadyPresent = true
				break
			}
		}
		if !finalizerAlreadyPresent {
			namespace.ObjectMeta.Finalizers = append(namespace.ObjectMeta.Finalizers, finalizer)
			_, err = kubeClient.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}

	// Step3: delete namespaces
	// iterating one by one as deleting all namespaces using k8s is throwing error in case a namespace doesn;t exists
	namespaces := []string{"pl", "olm", "zk-client", "px-operator"}
	for _, namespaceName := range namespaces {
		err := kubeClient.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	//Old code
	//out, err = shell.ExecWithDurationAndSuccessM(shell.GetPWD()+zkUninstallOperator, "zeroK operator uninstalled successfully")
	//if err != nil {
	//	filePath, _ := utils.DumpError(out)
	//	ui.GlobalWriter.PrintErrorMessage(fmt.Sprintf("installation failed, Check %s for details\n", filePath))
	//	return err
	//}

	//TODO: check all the namespaces and deactivate them
	return nil
}
