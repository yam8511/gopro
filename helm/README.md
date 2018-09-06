# Helm

## 安裝

```shell
curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh
```

1. [官方安裝網址](https://docs.helm.sh/using_helm/#installing-helm)
2. [部落格教學網址](https://jimmysong.io/posts/manage-kubernetes-native-app-with-helm/)

## 初始化

```shell
helm init

# 查看 pod 可以看到 tiller
kubectl get pods --namespace kube-system

# 設定權限
kubectl create -f rbac-config.yaml
```

## 建立 Chart

```shell
helm create mychart
```