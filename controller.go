package main

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	extv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/healthz"
)

type CaddyIngressController struct {
	clientset   *kubernetes.Clientset
	svcFallback string
}

func NewIngressController(clientset *kubernetes.Clientset, svcFallback string) (*CaddyIngressController, error) {
	return &CaddyIngressController{
		clientset:   clientset,
		svcFallback: svcFallback,
	}, nil
}

func (c *CaddyIngressController) ListIngresses(namespace string, opts v1.ListOptions) (*extv1beta1.IngressList, error) {
	return c.clientset.ExtensionsV1beta1().Ingresses(namespace).List(opts)
}

func (c *CaddyIngressController) RunHealthz(port int) {
	mux := http.NewServeMux()
	healthz.InstallHandler(mux)

	mux.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		c.Stop()
	})

	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Caddy Ingress Controller version %v", version)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
	}

	glog.Fatal(server.ListenAndServe())
}

func (c *CaddyIngressController) Stop() error {
	return nil
}

func (c *CaddyIngressController) WatchNamespace(namespace string) {
	glog.Infof("Listing ingresses in namespace %v", namespace)

	ingresses, err := c.ListIngresses(namespace, v1.ListOptions{})
	if err != nil {
		glog.Errorf("Could not list ingresses: %v", err)
		return
	}

	glog.Infof("Found %d ingresses (ResourceVersion %v)", len(ingresses.Items), ingresses.ListMeta.ResourceVersion)
	for ingress, idx := range ingresses.Items {
		glog.Infof(" #%d: %v", idx, ingress)
	}
}
