

# KLoops

This chart bootstraps installation of [KLoops](https://github.com/eddycharly/kloops).

## Installing

- Add kloops helm charts repo

```bash
helm repo add kloops TODO

helm repo update
```

- Install (or upgrade)

```bash
# This will install KLoops in the kloops namespace (with a kloops release name)

# Helm v2
helm upgrade --install kloops --namespace kloops kloops/kloops
# Helm v3
helm upgrade --install kloops --namespace kloops kloops/kloops
```

Look [below](#values) for the list of all available options and their corresponding description.

## Uninstalling

To uninstall the chart, simply delete the release.

```bash
# This will uninstall KLoops in the kloops namespace (assuming a kloops release name)

# Helm v2
helm delete --purge kloops
# Helm v3
helm uninstall kloops --namespace kloops
```

## Version

Current chart version is `0.0.1`

## Values

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| `chatbot.image.pullPolicy` | string | Docker image pull policy | `"IfNotPresent"` |
| `chatbot.image.repository` | string | Docker image repository | `"eddycharly/kloops-chatbot"` |
| `chatbot.image.tag` | string | Docker image tag | `"latest"` |
| `chatbot.ingress.annotations` | object | Ingress annotations | `{}` |
| `chatbot.ingress.enabled` | bool | Enable ingress | `false` |
| `chatbot.ingress.hosts` | list | Ingress host names | `[]` |
| `chatbot.replicaCount` | int | Number of replicas | `1` |
| `chatbot.resources.limits` | object | Pods resource limits | `{"cpu":"100m","memory":"256Mi"}` |
| `chatbot.resources.requests` | object | Pods resource requests | `{"cpu":"80m","memory":"128Mi"}` |
| `chatbot.service` | object | Service settings | `{"annotations":{},"externalPort":80,"internalPort":8090,"type":"ClusterIP"}` |
| `cluster.crds.create` | bool | Create custom resource definitions | `true` |
| `dashboard.image.pullPolicy` | string | Docker image pull policy | `"IfNotPresent"` |
| `dashboard.image.repository` | string | Docker image repository | `"eddycharly/kloops-dashboard"` |
| `dashboard.image.tag` | string | Docker image tag | `"latest"` |
| `dashboard.ingress.annotations` | object | Ingress annotations | `{}` |
| `dashboard.ingress.enabled` | bool | Enable ingress | `false` |
| `dashboard.ingress.hosts` | list | Ingress host names | `[]` |
| `dashboard.replicaCount` | int | Number of replicas | `1` |
| `dashboard.resources.limits` | object | Pods resource requests | `{"cpu":"100m","memory":"256Mi"}` |
| `dashboard.resources.requests.cpu` | string |  | `"80m"` |
| `dashboard.resources.requests.memory` | string |  | `"128Mi"` |
| `dashboard.service` | object | Service settings | `{"annotations":{},"externalPort":80,"internalPort":8090,"type":"ClusterIP"}` |

You can look directly at the [values.yaml](./values.yaml) file to look at the options and their default values.
