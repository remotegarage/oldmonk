

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

Oldmonk is all about turning day 2 operations into code! Not just that, it means you start thinking about day 2 on day 1. This is a dream come true for any Operations team!
Oldmonk leverages the strength of automation and combines it with the power of queue based workflows.

## Purpose

The Oldmonk provides the ability to implement these queue based deployment for any resources in Kubernetes. oldmonk does not care if you have a plain Kubernetes, a cloud based Kubernetes (like GKE), or a complete PaaS platform based on Kubernetes (like OpenShift). Oldmonk also does not care how you want to structure your data, how many queue you want to use.

Oldmonk can handle straight-up (static) yaml files with the complete definition or create dynamic ones based on your templating engine. Oldmonk supports *Helm Charts*, but can easily be extended to support others.

These templates will be merged and processed with a set of environment-specific parameters to get a list of resource manifests. Then these manifest can be created/updated/deleted in Kubernetes.

## Docs

 - [Oldmonk CRD API Reference](https://oldmonk.netlify.com/)

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

 ## Prerequisites

  - [git](#)
  - [go](#) version v1.13+.
  - [mercurial](#) version 3.9+
  - [docker](#) version 17.03+.
  - [kubectl](#) version v1.12.0+.
  - Access to a Kubernetes v1.12.0+ cluster.
  - [Operator-SDK](#)

## Get started in 5 Min
Install [Operator sdk](https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md), create a [GKE cluster](https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-cluster) and connect cluster
 ```bash
 cd ./demo
 # Start beanstalkd container using docker-compose
 docker-compose up -d
 cd ../
 make manager
 make install
 operator-sdk up local
 # Run it and Open a new terminal
 # Terminal  2
 watch kubectl get deployment
 # Terminal  3
 watch kubectl get QueueAutoScaler
 # Terminal  4
 kubectl apply -f example/demo/demo.yaml #custom resource for demo we are only deploying nginx image but in real case it's your worker image
 # Open http://127.0.0.1 and add job in default tube. cross the thresold of 4 and check terminal 1 log and also get deployment watch
 # After verifing the autoscale delete all jobs from default tube and again watch terminal1 and deployment
 kubectl delete -f example/demo/demo.yaml
 # It will also delete the deployment with crd defination

 ```

[![asciicast](https://asciinema.org/a/yh718d1AAyhiVAS9CyqstecAz.svg)](https://asciinema.org/a/yh718d1AAyhiVAS9CyqstecAz)



## Examples / Demos

We've created several examples for you to test out Oldmonk. See [EXAMPLES](example/) for details.

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

## Deploying to Kubernetes

```
#Replace image from Makefile
make docker-build
make docker-push
make install

```

## Deploying CRD

```
kubectl apply -f config/beanstalk.yaml

```


### Vanilla Manifests

You have to first clone or download the repository contents. The kubernetes deployment and files are provided inside `deploy/crd/` folder.

## Help

**Got a question?**
File a GitHub [issue](https://github.com/evalsocket/oldmonk/issues), or send us an [email](mailto:evalsocket@protonmail.com).

## Ref 

  - [Uswitch/sqs-autoscaler-controller](https://github.com/uswitch/sqs-autoscaler-controller)
