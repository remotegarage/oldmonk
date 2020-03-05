# Oldmonk - Build from scratch

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

## Build Docker Image

```bash
#Update Image name is Make file
make docker-build
make docker-push
```