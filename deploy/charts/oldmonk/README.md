# oldmonk

oldmonk scale kubernetes pods based on the Queue length of a queue in a Message Queueing Service. oldmonk automatically scales the number of pods in a deployment based on observed queue length.

## Prerequisites

- Kubernetes 1.12+

## Installing the Chart

To install the chart with the release name `my-release`:

```console
## IMPORTANT: you MUST install the oldmonk CRD **before** installing the oldmonk Helm chart.
$ kubectl apply --validate=false \
    -f https://raw.githubusercontent.com/evalsocket/oldmonk/master/deploy/crds/oldmonk.evalsocket.in_queueautoscalers_crd.yaml

## Add the Jetstack Helm repository
$ helm repo add oldmonk https://XXXXXXXX

## Install the oldmonk helm chart
$ helm install --name my-release --namespace kube-system oldmonk/oldmonk
```
> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the oldmonk chart and their default values.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `image` | Image name | `evalsocket/oldmonk` |
| `image.tag` | Image tag | `latest` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `serviceAccount.create` | If `true`, create a new service account | `true` |
| `serviceAccount.name` | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template |
| `rbac.create` | If `true`, create role/rolebinding | `true` |
| `resources` | CPU/memory resource requests/limits | `{}` |
| `extraEnv` | Optional environment variables for oldmonk | `[]` |
| `podAnnotations` | Annotations to add to the oldmonk pod | `{}` |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `affinity` | Node affinity for pod assignment | `{}` |
| `tolerations` | Node tolerations for pod assignment | `[]` |
| `queues` | Queues to autoscaling watch | `[]` |
| `options` | Oldmonk options | `[]` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
$ helm install --name my-release -f values.yaml .
```
> **Tip**: You can use the default [values.yaml](values.yaml)