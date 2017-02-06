# Kubernetes Ingress Controller for Caddy Webserver

This is my first dip into writing [Kubernetes](https://k8s.io/) [Ingress
Controllers](https://github.com/kubernetes/contrib/tree/master/ingress/controllers).
Specifically, I'm interested in:

* hosting the ingress controller outside of the cluster;
* utilizing NodePort as an entrypoint into the cluster;
* running [Caddy webserver](https://caddyserver.com/) with its built-in support
  for HTTP/2 and [ACME](https://letsencrypt.github.io/acme-spec/).

