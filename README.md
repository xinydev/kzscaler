# KZScaler(Kubernetes Zero Scaler)

KZScaler can enable scaling to/from zero feature for any HTTP service in Istio enabled Kubernetes clusters without any
modification

The highlight of KZScaler is that there is no need to modify the existing services to enable the above feature

[中文简介](README-CN.md)

## Components

* Outbound-Proxy:
  envoy's [wasm extension](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/wasm/v3/wasm.proto)
  , which analyzes the outbound request of the user and determines whether the target service is scaled to 0

* KZScaler-Controller: CRD controller and web server

### 0->1

Hijack all outbound traffic, and analyze whether the number of instances of target service is 0. If it is 0, scale to 1.

### 1->0(WIP)

Hijack all inbound traffic and expose a metric, KZScaler controller

## Suitable scenarios

There are a large number of inactive microservices in the cluster. Microservices communicate with each other very
**infrequently**. Admin hopes to reduce the resource which microservices using

## Limit

If the clients of the service are mainly outside the cluster, it is recommended to directly use Keda-httpaddon or other
projects

At present, it is only a demo project and is under active development

## Installation

### Prepare the environment

#### kubernetes

```shell

brew install kind

kind create cluster --name k1 --image kindest/node:v1. twenty-three point three

```

#### istio

```shell

brew install istioctl

istioctl install --set profile=demo -y

```

#### User services for testing

```shell

kubectl create ns testns

kubectl label namespace testns istio-injection=enabled

kubectl apply -f example/userservices. yaml -n testns

```

### Install KZScaler

```shell

ko apply -f config/

kubectl apply -f config/600-envoyconfig. yaml -n testns

kubectl apply -f example/zeroscaler. yaml -n testns

```

### Verify

Change the number of demo server instances to 0, and then access the service to see if the number of instances will be
automatically increased

```shell

kubectl scale deployment -n testns demo-server --replicas 0
```

```shell

kubectl get deployments. apps -n testns demo-server
```

```
NAME          READY   UP-TO-DATE   AVAILABLE   AGE
demo-server   1/1      1           1           128m
```

mock requests

```shell
kubectl exec demo-client-746b5998bc-75h4n -n tetsns -it -- bash

apt update
apt install -y curl

curl -I http://demo-server.testns
```

You will see that the request will receive the result normally, and the number of instances has been increased to 1

```shell

kubectl get deployments. apps -n testns demo-server

```

```
NAME          READY   UP-TO-DATE   AVAILABLE   AGE
demo-server   1/1      1           1           128m
```

## Uninstall

```shell

kubectl delete -f config/600-envoyconfig. yaml -n testns

kubectl delete ns KZScaler testns

```

## Roadmap

- [ ] supports automatic scaling of idle instances to 0

- [ ] grpc supported

- [ ] reduce outbound proxy and KZScaler controller requests

## Contributing

If you are interested, have questions or found some bugs, etc., you can use Issues or e-mail(xinydev@gmail.com) to
provide feedback or **participate in project development**