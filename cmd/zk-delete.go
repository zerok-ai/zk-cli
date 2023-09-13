package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"zkctl/cmd/internal"
	"zkctl/cmd/internal/shell"
	"zkctl/cmd/pkg/k8s"
	"zkctl/cmd/pkg/ui"
	"zkctl/cmd/pkg/utils"
)

const (
	delSpinnerText = "deleting zerok daemon and associated CRDs"
	delSuccessText = "zerok removed from the cluster successfully"
	delFailureText = "failed to delete zerok daemon"

	zkUninstallOperator string = "/operator/uninstall.sh"
)

var (
	namespacesToDelete  = []string{"pl", "zk-client", "px-operator"}
	crdsToDelete        = []string{"zerokops.operator.zerok.ai", "zerokinstrumentations.operator.zerok.ai"}
	clusterRoles        = []string{"zk-daemonset", "zk-operator-metrics-reader", "zk-operator-proxy-role", "zk-operator-role"}
	clusterRoleBindings = []string{"zk-daemonset", "zk-operator-proxy-rolebinding", "zk-operator-rolebinding"}

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

func deleteClusterRolesSilent(client *k8s.Client, names ...string) {
	err := deleteClusterRoles(client, names...)
	if err != nil {
		return
	}
}
func deleteClusterRoles(client *k8s.Client, names ...string) error {
	policyClient := client.RbacV1().ClusterRoles()

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

func deleteNamespacesSilent(client *k8s.Client, names ...string) {
	err := deleteNamespaces(client, names...)
	if err != nil {
		return
	}
}

func deleteNamespaces(client *k8s.Client, names ...string) error {
	namespacesSet := client.CoreV1().Namespaces()

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

func deleteCRDSilent(name []string) {
	err := deleteCRD(name)
	if err != nil {
		return
	}
}
func deleteCRD(names []string) error {
	for _, name := range names {
		cmd := "kubectl delete crd " + name
		_, cmdErr := shell.Shellout(cmd)
		if !errors.IsNotFound(cmdErr) {
			return fmt.Errorf("failed to delete CRD %q: %v", name, cmdErr)
		}
	}
	return nil
}

func deleteClusterRoleBindingsSilent(client *k8s.Client, names ...string) {
	err := deleteClusterRoleBindings(client, names...)
	if err != nil {
		return
	}
}
func deleteClusterRoleBindings(client *k8s.Client, names ...string) error {
	authClient := client.RbacV1().ClusterRoleBindings()

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

	kubeConfig := viper.GetString(KUBECONFIG_FLAG)
	kubeContext := viper.GetString(KUBECONTEXT_FLAG)
	kubeClient, err := k8s.NewKubeClient(kubeConfig, kubeContext)
	if err != nil {
		return err
	}

	// Step0: Getting confirmation from user for deleting ZeroK from cluster
	clusterName, err = kubeClient.GetClusterName()
	promptMessage := fmt.Sprintf(
		"Remove Zerok from cluster (cluster: %s)", clusterName,
	)

	yesFlag := viper.GetBool(internal.YesFlag)
	if !ui.GlobalWriter.YesNoPrompt(promptMessage, yesFlag) {
		return internal.ErrExecutionAborted
	}

	// Step1: Uninstalling backend cli daemon
	cmd := utils.GetBackendCLIPath() + " delete"
	out, cmdErr := shell.ShelloutWithSpinner(cmd, delSpinnerText, delSuccessText, delFailureText)
	if cmdErr != nil {
		internal.DumpErrorAndPrintLocation(out)
		return cmdErr
	}

	// Step2: Adding finalizer to zk-client namespace so that zk-operator can properly shut itself down
	//finalizer := "operator/cleanup-pods"
	//finalizerAlreadyPresent := false
	//namespace, err := kubeClient.CoreV1().Namespaces().Get(context.TODO(), "zk-client", metav1.GetOptions{})
	//if err != nil {
	//	return err
	//}
	//for _, f := range namespace.ObjectMeta.Finalizers {
	//	if f == finalizer {
	//		finalizerAlreadyPresent = true
	//		break
	//	}
	//}
	//if !finalizerAlreadyPresent {
	//	namespace.ObjectMeta.Finalizers = append(namespace.ObjectMeta.Finalizers, finalizer)
	//	_, err = kubeClient.CoreV1().Namespaces().Update(context.TODO(), namespace, metav1.UpdateOptions{})
	//	if err != nil {
	//		return err
	//	}
	//}

	// Step3: Delete namespaces, CRDs, clusterRoles and clusterRoleBindings
	deleteNamespacesSilent(kubeClient, namespacesToDelete...)
	deleteCRDSilent(crdsToDelete)
	deleteClusterRolesSilent(kubeClient, clusterRoles...)
	deleteClusterRoleBindingsSilent(kubeClient, clusterRoleBindings...)

	return nil
}
