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

func deleteClusterRoles(clientset *k8s.Client, names ...string) error {
	policyClient := clientset.RbacV1().ClusterRoles()

	for _, name := range names {
		deletePropagation := metav1.DeletePropagationForeground
		deleteOptions := &metav1.DeleteOptions{
			PropagationPolicy: &deletePropagation,
		}

		err := policyClient.Delete(context.TODO(), name, *deleteOptions)
		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("failed to delete ClusterRole %q: %v", name, err)
			}
		}
	}

	return nil
}

func deleteNamespaces(clientset *k8s.Client, names ...string) error {
	namespacesSet := clientset.CoreV1().Namespaces()

	for _, namespaceName := range names {
		err := namespacesSet.Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	return nil
}

func deleteCRD(name string) error {
	cmd := "kubectl delete crd " + name
	_, cmdErr := shell.Shellout(cmd, false)
	if cmdErr != nil {
		return cmdErr
	}
	return nil
}

func deleteClusterRoleBindings(clientset *k8s.Client, names ...string) error {
	authClient := clientset.RbacV1().ClusterRoleBindings()

	for _, name := range names {
		deletePropagation := metav1.DeletePropagationForeground
		deleteOptions := &metav1.DeleteOptions{
			PropagationPolicy: &deletePropagation,
		}

		err := authClient.Delete(context.TODO(), name, *deleteOptions)
		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("failed to delete ClusterRoleBinding %q: %v", name, err)
			}
		}
	}

	return nil
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
	deleteNamespaces(kubeClient, "pl", "olm", "zk-client", "px-operator")
	//namespaces := []string{"pl", "olm", "zk-client", "px-operator"}
	//for _, namespaceName := range namespaces {
	//	err := kubeClient.CoreV1().Namespaces().Delete(context.TODO(), namespaceName, metav1.DeleteOptions{})
	//	if err != nil {
	//		if !errors.IsNotFound(err) {
	//			return err
	//		}
	//	}
	//}

	//Deleting CRD
	deleteCRD("zerokops.operator.zerok.ai")

	//Deleting clusterRoleBindings
	deleteClusterRoles(kubeClient, "zk-daemonset", "zk-operator-metrics-reader", "zk-operator-proxy-role", "zk-operator-role")
	deleteClusterRoleBindings(kubeClient, "zk-daemonset", "zk-operator-proxy-rolebinding", "zk-operator-rolebinding")

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
