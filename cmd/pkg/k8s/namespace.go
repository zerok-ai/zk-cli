package k8s

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"zkctl/cmd/pkg/ui"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	labelSuccessMessage        string = "Namespace '%s' has been activated for ZeroK"
	labelAlreadyPresentMessage string = "Namespace '%s' is already active for ZeroK"
	labelNotPresentMessage     string = "Namespace '%s' is NOT active for ZeroK"
	labelRemovedMessage        string = "Namespace '%s' has been deactivated for ZeroK"

	messageScale             string = "Scale deployment '%s' to %d replicas"
	messageErrorGettingScale string = "Failed to get the scale value for deployment '%s' : %v"
	messageErrorScale        string = "Failed to scale to %d for deployment '%s' : %v"
)

type ZkResource struct {
	namespaceName string
	clientset     *kubernetes.Clientset
}

type ZkNamespace ZkResource

func NewZkNamespace(name string) *ZkNamespace {

	clientsetLocal, errClient := getKubeClientSet()
	if errClient != nil {
		panic(errClient)
	}

	return NewZkNamespaceWithClientSet(name, clientsetLocal)
}

func NewZkNamespaceWithClientSet(name string, clientsetExt *kubernetes.Clientset) *ZkNamespace {

	ns := &ZkNamespace{
		namespaceName: name,
		clientset:     clientsetExt,
	}

	return ns
}

func getKubeClientSet() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return kubernetes.NewForConfig(config)
}

func (ns *ZkNamespace) AddLabel(key string, value string) error {

	namespace, err := ns.clientset.CoreV1().Namespaces().Get(context.Background(), ns.namespaceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// add label iff not exists
	if existingValue, found := namespace.ObjectMeta.Labels[key]; !found || strings.Compare(existingValue, value) != 0 {
		namespace.ObjectMeta.Labels[key] = value
		_, updateErr := ns.clientset.CoreV1().Namespaces().Update(context.Background(), namespace, metav1.UpdateOptions{})
		if updateErr != nil {
			fmt.Printf("Failed to update the namespace: %v", updateErr)
			return updateErr
		}
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf(labelSuccessMessage, ns.namespaceName))
	} else {
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf(labelAlreadyPresentMessage, ns.namespaceName))
	}
	// printAllLabels(namespace)

	return nil
}

func (ns *ZkNamespace) RemoveLabel(key string) error {

	namespace, err := ns.clientset.CoreV1().Namespaces().Get(context.Background(), ns.namespaceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// add label iff not exists
	if _, found := namespace.ObjectMeta.Labels[key]; found {
		delete(namespace.ObjectMeta.Labels, key)
		_, updateErr := ns.clientset.CoreV1().Namespaces().Update(context.Background(), namespace, metav1.UpdateOptions{})
		if updateErr != nil {
			fmt.Printf("Failed to update the namespace: %v", updateErr)
			return updateErr
		}
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf(labelRemovedMessage, ns.namespaceName))
	} else {
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf(labelNotPresentMessage, ns.namespaceName))
	}
	printAllLabels(namespace)

	return nil
}

func printAllLabels(namespace *v1.Namespace) {
	labels := namespace.GetLabels()
	labelCount := len(labels)
	ui.GlobalWriter.Printf("Number of labels of the namespace '%s' are: %d \n", namespace.GetName(), labelCount)

	if labelCount > 0 {
		ui.GlobalWriter.Printf("Labels of the namespace '%s' are: \n", namespace.GetName())
		for key, value := range labels {
			fmt.Printf("%s = %v\n", key, value)
		}
	}
}

func (ns *ZkNamespace) DoRollingRestart(namespaceString string) error {

	depI := ns.clientset.AppsV1().Deployments(namespaceString)
	depList, err := depI.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, deployment := range depList.Items {

		// get current scale
		scale, err := depI.GetScale(context.TODO(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			ui.GlobalWriter.PrintErrorMessageln(fmt.Sprintf(messageErrorGettingScale, deployment.Name, err))
			continue
		}
		var currRepCount int32 = scale.Spec.Replicas

		// scaledown
		errDownScale := scaleDeployment(scale, depI, deployment.Name, 0)
		if errDownScale != nil {
			ui.GlobalWriter.PrintErrorMessageln(fmt.Sprintf(messageErrorScale, 0, deployment.Name, errDownScale))
			continue
		}
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf(messageScale, deployment.Name, 0))

		// scaleback
		errScaleBack := scaleDeployment(scale, depI, deployment.Name, currRepCount)
		if errScaleBack != nil {
			ui.GlobalWriter.PrintErrorMessageln(fmt.Sprintf(messageErrorScale, currRepCount, deployment.Name, errScaleBack))
			continue
		}
		ui.GlobalWriter.PrintSuccessMessageln(fmt.Sprintf(messageScale, deployment.Name, currRepCount))
	}

	return nil
}

func scaleDeployment(scale *autoscalingv1.Scale, depI appsv1.DeploymentInterface, deploymentName string, replicaCount int32) error {
	// scaledown
	scale.Spec.Replicas = replicaCount
	_, errorScale := depI.UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if errorScale != nil {
		return errorScale
	}
	return errorScale
}
