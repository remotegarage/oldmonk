

# Oldmonk - a QueueAutoScaler Operator for Kubernetes


![Oldmonk](https://github.com/evalsocket/oldmonk/workflows/Oldmonk/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/evalsocket/oldmonk)](https://goreportcard.com/report/github.com/evalsocket/oldmonk)
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/evalsocket/oldmonk)
[![Netlify Status](https://api.netlify.com/api/v1/badges/8eba3823-ba99-4c4a-bca3-8d81d347f80d/deploy-status)](https://app.netlify.com/sites/oldmonk/deploys)

[![Docker Repository on Quay](https://dockeri.co/image/evalsocket/oldmonk "Docker Repository on Dockerhub")](https://dockeri.co/evalsocket/oldmonk)

NOTE :  Not production ready.

## Oldmonk

According to Wikipedia:

>Old Monk : Old Monk Rum is an iconic vatted Indian dark rum, launched in 1954. It is blended and aged for a minimum of 7 years. It is a dark rum with a distinct vanilla flavour, with an alcohol content of 42.8%.
>In 1855, an entrepreneurial Scotsman named Edward Abraham Dyer, father of Colonel Reginald Edward Harry Dyer of Jallianwala Bagh massacre, set up a brewery in **Kasauli, Himachal Pradesh** to cater to the British requirement for cheap beer. This brewery changed hands and became a distillery by the name of Mohan Meakin Pvt. Ltd
>Mohan named the rum in honour of the Benedictine religious order of Christianity, with which he was fascinated.
>Some believe that Ved Rattan Mohan met "**Gumnami baba**" few times during this period and he named his creation after Gumnami baba, who was an old monk and believed to be as **Netaji Subhas Chandra Bose** by a group of Netaji enthusiasts and researchers, though there is no particular proof behind this. This doctrine also supports the fact that **Netaji Subhas Chandra Bose** never died out of the plane crash and came back to India as a monk and lived a spiritual life in disguise.
>
## What is Oldmonk

Scale kubernetes pods based on the Queue length of a queue in a Message Queueing Service. oldmonk automatically scales the number of pods in a deployment based on observed queue length.

## Purpose

Kubernetes does support custom metric scaling using Horizontal Pod Autoscaler. Before making this everyone was using HPA to scale our worker pods. Below are the reasons for moving away from HPA and making a custom resource:

**TLDR;** Don't want to write and maintain custom metric exporters? Use Oldmonk to quickly start scaling your pods based on queue length with minimum effort (few kubectl commands and you are done !)

1. **No need to write and maintain custom metric exporters**: In case of HPA with custom metrics, the users need to write and maintain the custom metric exporters. This makes sense for HPA to support all kinds of use cases. Oldmonk comes with queue metric exporters(pollers) integrated and the whole setup can start working with 2 kubectl commands.

2. **Fast Scaling**: Everyone wanted to achieve super fast near real time scaling. As soon as a job comes in queue the containers should scale if needed. The concurrency, speed and interval of sync have been made configurable to keep the API calls to minimum.

3. **Platform Independent**: Oldmonk does not care if you have a plain Kubernetes, a cloud based Kubernetes (like GKE), or a complete PaaS platform based on Kubernetes (like OpenShift). Oldmonk also does not care how you want to structure your data, how many queue you want to use.

## Example

The configuration is described in the QueueAutoScaler CRD, here is an example:

```yaml
apiVersion: oldmonk.evalsocket.in/v1
kind: QueueAutoScaler
metadata:
  name: lifecycle
spec:
  type : "BEANSTALKD"
  option :
    tube: 'default'
    key : 'current-jobs-ready'
  secrets : 'rabbitmq'
  maxPods : 6
  minPods : 1
  scaleDown :
    amount : 1
    threshold : 3
  scaleUp :
    amount : 1
    threshold : 4
  deployment : 'demo-app'
  autopilot : false
```

### Support
 - [x] Rabbitmq
 - [x] SQS
 - [x] Beanstalk
 - [x] Nats.io (Not Tested on scale)
 - [ ] Kafka

## Get started in 5 Min
[docs/build.md](/docs/build.md)

### Install
Running the below script will create the Oldmonk [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and start the controller. The controller watches over all the specified CRD and scales the Kubernetes deployments based on the specification.

```bash

minikube start --kubernetes-version v1.15.0 

make install
```

### Verify Installation
Check the QueueAutoScaler resource is accessible using kubectl

```bash
kubectl get QueueAutoScaler
```


### Examples / Demos
Do install the controller before going with the example. We've created several examples for you to test out Oldmonk. See [EXAMPLES](example/) for details.

- Create Deployment that needs to scale based on queue length.
```bash
kubectl create -f example/demo/worker.yaml
```

- Create `Oldmonk object (lifecycle)` that will start scaling the `demoapp` based on SQS queue length.
```bash
kubectl create -f example/sqs.yaml
```

[![asciicast](https://asciinema.org/a/yh718d1AAyhiVAS9CyqstecAz.svg)](https://asciinema.org/a/yh718d1AAyhiVAS9CyqstecAz)


## Monitoring

### Monitoring with Prometheus

[Prometheus](https://prometheus.io/) is an open-source systems monitoring and alerting toolkit.

Prometheus collects metrics from monitored targets by scraping metrics HTTP endpoints.

- [configuring-prometheus](https://prometheus.io/docs/introduction/first_steps/#configuring-prometheus)
- `scrape_configs` controls what resources Prometheus monitors.
- `kubernetes_sd_configs` Kubernetes SD configurations allow retrieving scrape targets. Please see [kubernetes_sd_configs](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#endpoints) for details.
- Additionally, `relabel_configs` allow advanced modifications to any target and its labels before scraping.

By default, the metrics in Operator SDK are exposed on `0.0.0.0:8383/metrics`

For more information, see [Metrics in Operator SDK](https://github.com/operator-framework/operator-sdk/blob/v0.8.1/doc/user/metrics/README.md)

#### Usage:

```
scrape_configs:
  - job_name: 'kubernetes-service-endpoints'
    kubernetes_sd_configs:
    - role: endpoints
    relabel_configs:
      - source_labels: [__meta_kubernetes_namespace]
        action: keep
        regex: test-oldmonk-operator
```
You can find additional examples on their [GitHub page](https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus-kubernetes.yml).

#### Verify metrics port:
kubectl exec `POD-NAME` curl localhost:8383/metrics  -n `NAMESPACE`

(e.g. `kubectl exec oldmonk-operator-5b9b664cfc-6rdrh curl localhost:8383/metrics  -n oldmonk-operator`)


**Got a question?**

Please feel free to open issues in [Github](https://github.com/evalsocket/oldmonk/issues) if you have any questions or concerns. 

We hope Oldmonk is of use to you! Made with :heart: for Open source 

## Thanks

Thanks to kubernetes team for making [crds](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/), [Operator-sdk](https://github.com/operator-framework/operator-sdk)
and [Uswitch/sqs-autoscaler-controller](https://github.com/uswitch/sqs-autoscaler-controller) 
