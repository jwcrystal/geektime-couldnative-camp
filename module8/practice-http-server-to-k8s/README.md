# module8 - http server 部署至 k8s

> note: 此OS架構為 linux/arm64

此項目採用之前module3建立的`httpserver image`
把環境參數或是設定透過`ConfigMap`進行控制。

### Usage

把`pod`和相對應的`configmap`加入集群即可訪問。

```shell
kubectl create -f deploy.yaml
kubectl create -f config.yaml // config map
```

### 項目部署實現

- 優雅啟動
- 優雅關閉
- 配置資源需求跟Qos
- 實現 liveness、readiness探活


### Demo

![demo_printenv.png](assets/demo_printenv.png)
![demo_healthz.png](assets/demo_healthz.png)