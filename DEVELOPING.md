# Developing Guide

## find istio cluster

```shell
istioctl proxy-config all istio-ingressgateway-69dc56d7f-tzfrq -n istio-system -o json | grep kzscaler
```
