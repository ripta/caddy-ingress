package main

import (
	"flag"
	"fmt"
	_ "net/http"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	_ "k8s.io/kubernetes/pkg/healthz"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

var (
	version = "xxx:REPLACE:DURING:BUILD"
	flagSet = pflag.NewFlagSet("", pflag.ExitOnError)

	isInCluster = flagSet.Bool("is-in-cluster", false,
		`Whether the controller is run inside k8s or not (default: false)`)

	svcFallback = flagSet.String("fallback-service", "",
		`The name of the fallback service in "namespace/name" format.`)
)

func main() {
	glog.Infof("Caddy Ingress Controller version %v", version)

	flagSet.AddGoFlagSet(flag.CommandLine)
	flagSet.Parse(os.Args)
	clientConfig := util.DefaultClientConfig(flagSet)

	kubeClient, err := createClient(*isInCluster, &clientConfig)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %v", err)
	}

	//watchNamespace, namespaceGiven, err := clientConfig.Namespace()
	//if err != nil {
	//	glog.Fatalf("Unexpected error: %v", err)
	//}

	//if !namespaceGiven {
	//	watchNamespace = api.NamespaceAll
	//}

	ingressController, err := NewIngressController(kubeClient, *svcFallback)
	if err != nil {
		glog.Fatalf("Could not create ingress controller: %v", err)
	}

	ingressController.RunHealthz(8080)
	ingressController.WatchNamespace(api.NamespaceAll)
}

func createClient(isInCluster bool, clientConfig *clientcmd.ClientConfig) (*unversioned.Client, error) {
	// Generate an in-cluster client when possible
	if isInCluster {
		kubeClient, err := unversioned.NewInCluster()
		if err != nil {
			return nil, fmt.Errorf("Could not initialize in-cluster client: %v", err)
		}

		return kubeClient, nil
	}

	// Otherwise, manually initialize one
	kubeConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not create client configuration: %v", err)
	}

	kubeClient, err := unversioned.New(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not initialize generic client: %v", err)
	}

	return kubeClient, nil
}
