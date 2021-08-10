# shelly-exporter

Prometheus multi-target exporter that collects status data from the
[Shelly device API](https://shelly-api-docs.shelly.cloud).

![Build Docker image](https://github.com/ianunruh/shelly-exporter/actions/workflows/docker-build.yml/badge.svg)

Tested with the Shelly Plug US.

[Example metrics](docs/example-metrics.txt)

## Features

* Handles multiple devices with a single exporter, multiple relays per device
* Lightweight Docker image (less than 15MB)

## Usage

```bash
# Only required if device has restricted login
export SHELLY_USERNAME="admin"
export SHELLY_PASSWORD="changeme"

docker compose up -d --build
docker compose ps
docker compose logs exporter

curl -s "http://localhost:9090/probe?target=192.168.1.x"
```

## Deployment

Use Kustomize to deploy the exporter to a Kubernetes cluster.

```bash
kubectl -n monitoring create secret generic shelly-exporter \
    --from-literal=SHELLY_USERNAME=${SHELLY_USERNAME} \
    --from-literal=SHELLY_PASSWORD=${SHELLY_PASSWORD}

kubectl kustomize "https://github.com/ianunruh/shelly-exporter.git/deploy/basic?ref=v1.0.0" | \
    kubectl apply -n monitoring -f-
```

Then create the necessary resources for configuring Prometheus. Note the label used on
the service, which the service monitor is configured to select.

```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: shelly-plug
  labels:
    app.kubernetes.io/name: shelly-exporter-target
spec:
  type: ClusterIP
  clusterIP: None
  ports:
  - name: http
    port: 80
---
apiVersion: v1
kind: Endpoints
metadata:
  name: shelly-plug
subsets:
- addresses:
  - ip: 192.168.1.186
  - ip: 192.168.1.188
  ports:
  - name: http
    port: 80
```
