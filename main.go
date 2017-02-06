package main

import (
	"flag"
	"fmt"
	_ "net/http"
	_ "os"

	"github.com/golang/glog"

	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	version    = "xxx:REPLACE:DURING:BUILD"
	kubeconfig = flag.String("kubeconfig", "/etc/kubernetes/kubeconfig", "Absolute path to the kubeconfig file")

	isInCluster = flag.Bool("is-in-cluster", false,
		`Whether the controller is run inside k8s or not (default: false)`)

	svcFallback = flag.String("fallback-service", "",
		`The name of the fallback service in "namespace/name" format.`)
)

func main() {
	glog.Infof("Caddy Ingress Controller version %v", version)
	flag.Parse()

	config, err := getClusterConfig(*isInCluster, *kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to get cluster configuration: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %v", err)
	}

	//watchNamespace, namespaceGiven, err := clientConfig.Namespace()
	//if err != nil {
	//	glog.Fatalf("Unexpected error: %v", err)
	//}

	//if !namespaceGiven {
	//	watchNamespace = v1.NamespaceAll
	//}

	ingressController, err := NewIngressController(clientset, *svcFallback)
	if err != nil {
		glog.Fatalf("Could not create ingress controller: %v", err)
	}

	//go ingressController.RunHealthz(8080)
	ingressController.WatchNamespace(v1.NamespaceAll)
}

func getClusterConfig(isInCluster bool, kubeconfig string) (*rest.Config, error) {
	if isInCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("Could not initialize in-cluster context: %v", err)
		}
		return config, nil
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Could not use current context from %s: %v", kubeconfig, err)
	}

	return config, nil
}
