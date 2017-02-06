package main

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
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

func (c *CaddyIngressController) WatchNamespace(namespace string) {
}

func (c *CaddyIngressController) Stop() {
}
