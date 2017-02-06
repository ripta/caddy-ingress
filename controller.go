package main

import (
	"k8s.io/client-go/kubernetes"
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
	_ = port
}

func (c *CaddyIngressController) WatchNamespace(namespace string) {
}
