package install

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"log"
	"math/big"
	"os"
	"time"
	"zkctl/cmd/internal/shell"
	"zkctl/cmd/pkg/k8s"
)

var x509Name = pkix.Name{
	Organization: []string{"Pixie Labs Inc."},
	Country:      []string{"US"},
	Province:     []string{"California"},
	Locality:     []string{"San Francisco"},
}

const bitsize = 4096

type certGenerator struct {
	ca    *x509.Certificate
	caKey *rsa.PrivateKey
}

func getCloudDNSNamesForNamespace(namespace string) []string {
	return []string{
		fmt.Sprintf("*.%s", namespace),
		fmt.Sprintf("*.%s.svc.cluster.local", namespace),
		fmt.Sprintf("*.%s.pod.cluster.local", namespace),
		fmt.Sprintf("*.pl-nats.%s.svc", namespace),
		"*.pl-nats",
		"pl-nats",
		"*.local",
		"localhost",
	}
}

func newCertGenerator() (*certGenerator, error) {
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(1653),
		Subject:               x509Name,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caKey, err := rsa.GenerateKey(rand.Reader, bitsize)
	if err != nil {
		return nil, err
	}

	return &certGenerator{
		ca:    ca,
		caKey: caKey,
	}, nil
}

func (cg *certGenerator) generateSignedCertAndKey(dnsNames []string) ([]byte, []byte, error) {
	cert := &x509.Certificate{
		SerialNumber:          big.NewInt(1658),
		Subject:               x509Name,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 0, 0),
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, bitsize)
	if err != nil {
		return nil, nil, err
	}

	return cg.signCertAndKey(cert, privateKey)
}

func (cg *certGenerator) signCertAndKey(cert *x509.Certificate, privateKey *rsa.PrivateKey) ([]byte, []byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cg.ca, &privateKey.PublicKey, cg.caKey)
	if err != nil {
		return nil, nil, err
	}

	certData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	if err != nil {
		return nil, nil, err
	}

	keyData := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return nil, nil, err
	}

	return certData, keyData, nil
}

func (cg *certGenerator) signedCA() ([]byte, error) {
	caCertData, _, err := cg.signCertAndKey(cg.ca, cg.caKey)
	if err != nil {
		return nil, err
	}
	return caCertData, nil
}

func CreateGenericSecretFromLiterals(namespace, name string, fromLiterals map[string]string) (*v1.Secret, error) {
	secret := &v1.Secret{}
	secret.SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("Secret"))

	secret.Name = name
	secret.Data = map[string][]byte{}
	secret.Namespace = namespace

	for k, v := range fromLiterals {
		secret.Data[k] = []byte(v)
	}

	return secret, nil
}

func ConvertResourceToYAML(obj runtime.Object) (string, error) {
	buf := new(bytes.Buffer)
	e := jsonserializer.NewYAMLSerializer(jsonserializer.DefaultMetaFactory, nil, nil)
	err := e.Encode(obj, buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type Resource struct {
	Object *unstructured.Unstructured
	GVK    *schema.GroupVersionKind
}

func GetResourcesFromYAML(yamlFile io.Reader) ([]*Resource, error) {
	resources := make([]*Resource, 0)

	decodedYAML := yaml.NewYAMLOrJSONDecoder(yamlFile, 4096)

	for {
		ext := runtime.RawExtension{}
		err := decodedYAML.Decode(&ext)

		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if ext.Raw == nil {
			continue
		}

		_, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var unstructRes unstructured.Unstructured
		unstructRes.Object = make(map[string]interface{})
		var unstructBlob interface{}

		err = json.Unmarshal(ext.Raw, &unstructBlob)
		if err != nil {
			return nil, err
		}

		unstructRes.Object = unstructBlob.(map[string]interface{})

		resources = append(resources, &Resource{
			Object: &unstructRes,
			GVK:    gvk,
		})
	}

	return resources, nil
}

func parseSecretYAML(yamlString string) *v1.Secret {
	decoder := serializer.NewCodecFactory(scheme.Scheme).UniversalDeserializer()

	obj, _, err := decoder.Decode([]byte(yamlString), nil, nil)
	if err != nil {
		fmt.Printf("Error decoding YAML: %v\n", err)
		os.Exit(1)
	}

	secret, ok := obj.(*v1.Secret)
	if !ok {
		fmt.Println("Error: Not a valid Secret YAML")
		os.Exit(1)
	}

	return secret
}

func GenerateCloudCertYAMLs(namespace string) (string, error) {

	//Check if secret already exists
	cmd := "kubectl get secret service-tls-certs -n " + namespace + " -o jsonpath='{.data.ca.crt}'"
	_, secretExistsErr := shell.Shellout(cmd)
	if secretExistsErr == nil {
		//secret already exists, so returning without doing anything
		return "", nil
	}

	cg, err := newCertGenerator()
	if err != nil {
		return "", err
	}

	clientCert, clientKey, err := cg.generateSignedCertAndKey(getCloudDNSNamesForNamespace(namespace))
	if err != nil {
		return "", err
	}
	serverCert, serverKey, err := cg.generateSignedCertAndKey(getCloudDNSNamesForNamespace(namespace))
	if err != nil {
		return "", err
	}
	caCert, err := cg.signedCA()
	if err != nil {
		return "", err
	}

	tlsCert, err := CreateGenericSecretFromLiterals(namespace, "service-tls-certs", map[string]string{
		"server.key": string(serverKey),
		"server.crt": string(serverCert),
		"ca.crt":     string(caCert),
		"client.key": string(clientKey),
		"client.crt": string(clientCert),
	})
	if err != nil {
		return "", err
	}
	yaml, err := ConvertResourceToYAML(tlsCert)
	if err != nil {
		return "", err
	}

	yamlString := fmt.Sprintf("---\n%s\n", yaml)
	yamlSecret := parseSecretYAML(yamlString)
	clientSet, err := k8s.GetKubeClientSet()
	if err != nil {
		fmt.Printf("Error creating secret: %v\n", err)
		os.Exit(1)
	}
	_, err = clientSet.CoreV1().Secrets(namespace).Create(context.TODO(), yamlSecret, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error creating secret: %v\n", err)
		os.Exit(1)
	}

	return yamlString, nil
}
