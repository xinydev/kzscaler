# KZScaler(Kubernetes Zero Scaler)

KZScaler可以在安装了Istio的Kubernetes集群内，为存量HTTP服务支持 闲时缩放到零实例，忙时自动从零实例恢复的功能。

KZScaler的亮点是启用上述的功能不需要对现有服务进行修改

## 原理

主要使用了Envoy的[Wasm扩展](https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/wasm/v3/wasm.proto)

组件

* outbound-proxy: wasm扩展，分析用户的outbound请求，判断目标服务是否以及缩放到0了
* kzscaler-controller: crd controller和web server

### 0->1

劫持所有OUTBOUND流量，通过请求地址分析该地址的服务实例数是否为0，如果为0，则调用

### 1->0

正在实现中  
劫持所有INBOUND流量，暴露一个Metric，kzscaler-controller

## 适合场景

集群内有大量不活跃的微服务，微服务之间会互相通用，希望减少这里面的资源消耗

## 限制

如果服务的用户主要在集群外，推荐直接使用keda-httpaddon等方案

当前只是一个demo项目，正在积极开发中

## 安装

### 准备环境

#### kubernetes

```shell
brew install kind
kind create cluster --name k1 --image kindest/node:v1.23.3
```

#### istio

```shell
brew install istioctl
istioctl install --set profile=demo -y
```

#### 测试用的服务

```shell
kubectl create ns testns
kubectl label namespace testns istio-injection=enabled
kubectl apply -f example/userservices.yaml -n testns
```

### 安装KZScaler

```shell
kubectl apply -f https://github.com/kzscaler/kzscaler/releases/download/v0.0.1-alpha/release.yaml

# envoy配置
kubectl apply -f https://github.com/kzscaler/kzscaler/releases/download/v0.0.1-alpha/release-wasm.yaml -n testns

kubectl apply -f example/zeroscaler. yaml -n testns
```

### 验证

将demo-server实例数改为0，然后访问这个服务，看实例数会不会自动被增加

```shell
kubectl scale deployment -n testns demo-server --replicas 0
```

```shell
kubectl get deployments.apps -n testns demo-server
```

```
NAME          READY   UP-TO-DATE   AVAILABLE   AGE
demo-server   0/0      0           0           128m
```

模拟请求

```shell
kubectl exec demo-client-746b5998bc-75h4n -n tetsns -it -- bash
apt update
apt install -y curl
curl -I http://demo-server.testns
```

将会看到这个请求会正常收到结果,实例数已经被增加到了1个

```shell
kubectl get deployments.apps -n testns demo-server
```

```
NAME          READY   UP-TO-DATE   AVAILABLE   AGE
demo-server   1/1      1           1           128m
```

## 删除

```shell
kubectl delete -f config/600-envoyconfig.yaml -n testns
kubectl delete ns kzscaler testns  
```

## Roadmap

- [ ] 支持自动将空闲实例缩放到0
- [ ] 支持gRPC
- [ ] 减少outbound-proxy与kzscaler-controller的请求

## Contributing

感兴趣，有疑问，发现bug等等都可以通过Issues,邮件(xinydev@gmail.com)进行反馈，或者参与到项目开发