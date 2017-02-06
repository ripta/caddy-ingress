package main

import (
	_ "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
	_ "k8s.io/kubernetes/pkg/healthz"
)

type CaddyIngressController struct {
	client      *unversioned.Client
	svcFallback string
}

func NewIngressController(kubeClient *unversioned.Client, svcFallback string) (*CaddyIngressController, error) {
	return &CaddyIngressController{
		client:      kubeClient,
		svcFallback: svcFallback,
	}, nil
}

func (c *CaddyIngressController) RunHealthz(port int) {
	_ = port
}

func (c *CaddyIngressController) WatchNamespace(namespace string) {
}
