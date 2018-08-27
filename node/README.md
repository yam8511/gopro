# Node 操作

- 顯示Node (順便顯示標籤)

kubectl get nodes --show-labels

- 顯示Node的詳細資訊

kubectl describe nodes <my_node>

e.g.

```shell
kubectl describe nodes zuolar-mint
```

- 新增Node的標籤

kubectl label nodes <my_node> <label_key>=<label_value>

e.g.

```shell
kubectl label nodes zuolar-mint app=test
```

- 移除Node的標籤

kubectl label nodes <my_node> <label_key>-

e.g.

```shell
kubectl label nodes zuolar-mint app-
```
