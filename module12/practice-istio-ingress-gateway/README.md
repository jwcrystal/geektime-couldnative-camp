# module12 - Istio Ingress Gateway

> note: 此OS架構為 linux/arm64

## Install Istio for arm64

- download istioctl

```shell
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.14.3 TARGET_ARCH=arm64 sh -
```

- use istio operator 

```shell
# 部署 istio operator
istioctl operator init --hub=ghcr.io/resf/istio
# install with istio operator
kubectl create ns istio-system
kubectl apply -f - <<EOF
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  namespace: istio-system
  name: example-istiocontrolplane
spec:
  profile: demo
EOF
# ------------
kubectl get pods -n istio-system
# NAME                                    READY   STATUS    RESTARTS   AGE
# istio-egressgateway-6d4d975594-khkrw    1/1     Running   0          32s
# istio-ingressgateway-69984d8d6d-rb8jp   1/1     Running   0          32s
# istiod-54cfc7dbf8-bwq4p                 1/1     Running   0          44s
```

## Usage

- 建立一個test namespac為`httpmesh`

```shell
kubectl create ns httpmesh
```

- 注入Istio sidecar到kubernets pods

```shell
kubectl label ns httpmesh istio-injection=enabled
```

- Check ingress ip

```shell
k get svc -nistio-system
# ----
# istio-ingressgateway   LoadBalancer   $INGRESS_IP
```

- 手動簽發cert 或是 使用**letsencrypt**申請cert

```shell
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cloudnative Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt  

kubectl create -n istio-system secret tls cncamp-tls --key=cncamp.io.key --cert=cncamp.io.crt
```

- Access httpserver via ingress

```shell
curl --resolve httpsserver.cncamp.io:443:$INGRESS_IP https://httpsserver.cncamp.io/healthz -v -k
```

## Demo

![image-20220814232624788](/Users/jwang/Documents/Workspace/GoPractice/Projects/cncamp/module12/practice-istio-ingress-gateway/assets/image-20220814232624788.png)

![image-20220814233036679](/Users/jwang/Documents/Workspace/GoPractice/Projects/cncamp/module12/practice-istio-ingress-gateway/assets/image-20220814233036679.png)

![image-20220814233214754](/Users/jwang/Documents/Workspace/GoPractice/Projects/cncamp/module12/practice-istio-ingress-gateway/assets/image-20220814233214754.png)

- 這邊顯示的IP為節點IP

![image-20220814233255301](/Users/jwang/Documents/Workspace/GoPractice/Projects/cncamp/module12/practice-istio-ingress-gateway/assets/image-20220814233255301.png)
