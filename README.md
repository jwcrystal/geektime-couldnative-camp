# 總結筆記

參加雲原生第4期，5/13開課。

到現在不知不覺的結束了，沒想到這麼快。。。

一路上從孟老師和助教們身上學到很多知識和經驗

我自己是自學Golang的，時間起於開課前幾週，所以前2課的練習對我自己是有所幫助的，更加熟悉語言特性，再到後面的容器docker、容器編排平台Kubernetes，一路學習上去。

萬丈高樓平地起，最後來複習一下kubernets最基本元件做結尾，即使後來好用的組件，也是基於基本元件展開的。

## Go語言

- 藉由課後練習 1.1-1.2，熟悉`語言特性`，及`goroutine`的用法

```go
// 課後練習
基于 Channel 编写一个简单的单线程生产者消费者模型：
队列：
队列长度 10，队列元素类型为 int
生产者：
每 1 秒往队列中放入一个类型为 int 的元素，队列满时生产者可以阻塞
消费者：
每一秒从队列中获取一个元素并打印，队列为空时消费者阻塞
// ---
func main() {
	ch := make(chan string)
	go Producer(ch)
	Consumer(ch)
}

func Producer(c chan string) {
	for i := 0; i < 10; i++ {
		c <- strconv.Itoa(i)
		time.Sleep(time.Second)
		fmt.Printf("Producer[%d]\n", i)
	}
	close(c)
}

func Consumer(c chan string) {
	//fmt.Printf("Concumer[%s]\n", <-c)
	for v := range c {
		fmt.Printf("Concumer[%s]\n", v)
	}
}

```

- 並使用Golang原生的`http.handler`，撰寫一個http server
  
    ```go
    type HandlerFunc func(w http.ResponseWriter, res *http.Request)
    
    type Engine struct {
    	router map[string]HandlerFunc
    	//middlewares []HandlerFunc
    }
    
    func New() *Engine {
    	return &Engine{router: make(map[string]HandlerFunc)}
    }
    
    // implement the ServeHTTP interface from net/http
    func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    	key := r.Method + "-" + r.URL.Path
    	if handler, ok := e.router[key]; ok {
    		// using middleware is a good choice, but skip this part here
    		//e.middlewares = append(e.middlewares, handler)
    
    		// Logger
    		LogHandler(http.StatusOK, r)
    		handler(w, r)
    	} else {
    		w.WriteHeader(http.StatusNotFound)
    		fmt.Fprintf(w, "404 Not Found: %s\n", r.URL)
    	}
    }
    ```
    
    - 當時把這個封裝成簡易框架設計，以便使用

```go
// Module2 - a simple http server

1. build demo and launch it
$ go build -o demo
$ ./demo

/*
>>log
http server :8080  started
2022/06/01 18:47:12 [200] 127.0.0.1 - / - 3.417µs
*/

2. Provide 2 API
- /
  - return IP of client
- /healthz
  - return data of header and specific variable ("Version") in env
  - log request on the terminal of server

```

## Docker

學習docker使用及背後原理和技術：`cgroups`、`namespace`、`Union FS`

### cgroups

- 一種資源控制方案
- 對CPU、Memory、硬體I/O進行限制
- 可以用`Hierarchy`組織管理，sub-cgroup除了受到自身限制外，也受到parent-cgroup的資源限制

### namespace

- 一種 Linux Kernel提供資源隔離方案
- 為process提供不同的`namespace`進行區隔，互不干擾
- 每個`namespace`會包含以下這些

| Namespace Type | 隔離資源 |
| --- | --- |
| IPC | System V IPC 和 POSIX 消息對列 |
| Network | 網路設備、協議、端口等 |
| PID | 進程 process |
| Mount | 掛載點 |
| UTS | 主機名、域名 |
| USR | 用戶、用戶組 |
- 課後練習

```bash
### Namespace
# 在new network namespace 執行 sleep指令
unshare -fn sleep 60
# 查看進程
ps -ef | grep sleep
root 32882 4935 0 10:00 pts/0 00:00:00 unshare -fn sleep 60
root 32883 32882 0 10:00 pts/0 00:00:00 sleep 60
# 查看 Network namespace
lsns -t net
4026532508 net 2 32882 root unassigned unshare
# 進入進程，查看ip
nsenter -t 32882 -n ip a
1: lo <loopback> mtu 65536 ...
		link/loopback 00:00:00 ...

### Cgroups
# 新建一個test dir：cpudemo
cd /sys/fs/cgroup/cpu/cpudemo
# 找到busyloop的PID，加入到cgroup.proc
ps -ef | grep busyloop
echo <pid> > cgourp.proc
# 或是直接透過下面指令加入
ps -ef | grep busyloop | grep -v grep | awk '{print $2}' > cgroup.proc
#  cpu.cfs_period_us 用來配置CPU時間週期長度
#  cpu.cfs_quota_us 當前cgroup在cpu.cfs_period_us中最多能使用的cpu時間數
echo 10000 > cpu.cfs_quota_us
# 執行 top 查看CPU使用情況 是否變為 10%
```

- union FS (file system)
    - 將不同目錄掛載到同一個虛擬文件系統下
    - 另一個更常用的是將一個`readonly`的branch和一個`writeable`的branch聯合一起

### docker

- 啟動將`rootfs`以`readonly`方式讀取並檢查，利用`union mount`將一個`readwrite`文件掛載在`readonly`的`rootfs`上
- 一組`readonly`和`writeable`結構為一個`container`運作方式，為一個`FS`層。並可以將下層的`FS`設為`readonly`向上疊加
- 鏡像具有共享特性，所以對容器可寫層操作需要依靠存儲驅動提供的機制來支持，提供對存儲和資源的使用率
    - 寫時複製 COW
        - copy-on-write
        - 從鏡像文件系統複製到容器的可寫層文件系統進行修改，跟原本的文件，相互獨立
    - 用時分配
        - 被創建後才分配空間
- 課後練習

```bash
# OverlayFS practice
mkdir upper lower merged work
echo "from lower" > lower/in_lower.txt
echo "from upper" > lower/in_upper.txt
echo "from lower" > lower/in_both.txt
echo "from upper" > lower/in_both.txt

sudo mount -t overlay overlay -o lowerdir=`pwd`/lower, upperdir=`pwd`/upper, workdir=`pwd`/work `pwd`/merged
cat merged/in_both.txt
```

- 課後練習
    - 因OS結構為 `linux/arm64`，課後環境都需要而外建立無法用正常的`amd64`鏡像
        - 此在縮小鏡像體積時候，遇到以下錯誤
            - solved： 在`dockerfile`中加入`cgo=0`的環境變數
                - 因為會有不同平台會有動態連結的問題
    - 鏡像可以透過多端構建，只留build好的部分，排除過程產物，減少鏡像體積
      
        ```docker
        # 前者路徑為建置好的檔案路徑
        COPY --from=builder /build/demo/http-server /
        ```
        
    - `golang:scratch`和`golang:alpine`鏡像
        - 前者雖然體積小，但是沒包含基本除錯工具
    
    ```markdown
    # Resize Dockerfile
    ```dockerfile
    FROM golang:1.18-alpine AS builder
    
    ENV CGO_ENABLED=0 
    #   GO111MODULE=off  \
    #	GOOS=linux    \
    #	GOARCH=amd64
    
    WORKDIR /build
    COPY . .
    RUN echo "Install dependent modules" && \
        go mod download && go mod verify && \
        cd demo/ && \
        go build -o http-server .
    
    FROM busybox
    COPY --from=builder /build/demo/http-server /
    EXPOSE 8080
    CMD ["/http-server"]
    #ENTRYPOINT ["/http-server"]
    ```
    ```
    
    - 多進程容器鏡像
        - 需要捕捉`SIGTERM`完成優雅退出
        - 清理退出子進程避免殭屍進程
        - 可以從代碼優雅退出
        - 或是透過[tini](https://github.com/krallin/tini)監控

## Kubernetes

- 了解前身Google Borg的來歷
- 基於容器的應用部署、維護、滾動升級

### 聲明式API，核心對象

- `Node`
    - 節點的抽象，描述計算節點的狀況
    - `Node`是`Pod`真正運行的主機
- `Namespace`
    - 資源隔離的基本單位
    - 一組資源和對象的抽象集合
- `Pod`
    - 用來描述應用實例，kubernets最核心對象，也是調度的基本單位
    - 同一個`Pod`中的不同容器課共享資源
        - 共享`Network Namespace`
        - 可通過掛載`存儲卷`共享存儲
            - 存儲卷：從外部存儲掛載到`Pod`內部使用
            - 分為兩部分：`Volume`和`VolumeMounts`
                - Volume：定義Pod可以使用的存儲來源
                - VolumeMounts：定義如何掛載到容器內部
        - 共享`Security Context`
    - 單機限制`110`個Pod
- `Service`：將應用發布成服務，本質是負載均衡和域名服務的聲明
- 每個API對象都有四大屬性

### TypeMeta

- 通過此引GKV（`Group`, `Kind`, `Version`）模型定義一個對象類型
    - `Group`：將對象依據其功能分組
    - `Kind`：定義對象的基本類型
        - e.g. Node、Pod、Deployment等
    - `Version`：每季度會推出Kubernetes版本
        - e.g. v1alpha1、v1alpha2、v1（生產版本）

### MetaData

- 兩個重要屬性：`Namespace`和`Name`，分別定義對象的namespace歸屬，這兩屬性唯一定義了對象實例
    - Label
        - `KV-pairs`，kubernetes api支持以`label`作為過濾條件
        - label selector支持以下方式
            - 等式，如`app=nginx`或是`env≠production`
            - 集合，如 `env in (production, qa)`
    - Annotation
        - `KV-paris`，此為屬性擴展，更多面向開發及管理人員，所以需要像其他屬性合理歸類
        - 用來記錄一些附加訊息，如`deployment`使用`annotation`來記錄`rolling update狀態`
    - Finalizer
        - 本質為一個資源鎖
        - kubernetes在接受對象的刪除請求時，會檢查`Finalizer`是否為空，不為空只做邏輯刪除，即只會更新對象中的`metadata.deletionTimestamp`字段
    - ResourceVersion
        - 類似一個樂觀鎖
        - Kubernetes對象被客戶端讀取後`ResouceVersion`訊息也同時被讀取。此機制確保了分布式系統中任意多線程能夠無鎖併發訪問對象

### Spec & Status

- Spec：用戶期望方式，由用戶自行定義
    - 健康檢查
        - 探針類型
            - LivenessProbe
                - 檢查應用是否健康，若否則刪除並重新創建容器
            - ReadinessProbe
                - 檢查應用是否就緒且為正常服務狀態，若否則不會接受來時`Kubernets Services`的流量
            - StartupProbe
                - 檢查應用是否啟動成功，如果在`failureThreshold*periodSeconds`週期內為就緒，則應用會被重啟
        - 探活方式
            - Exec
            - Tcp socket
            - Http
- Status：對象實際狀態，由對應`Controller`收集狀態並更新
- 跟通用屬性不同，`Spec`和`Status`是每個對象獨有的

### Deployment instance yaml

```yaml
apiVersion: v1
kind: Deployment
metadata:
	labels:
		app: nginx
	name: nginx
spec:
	replica: 3 # 3 副本啟動
	selector:
		matchLabels:
			app: nginx
	template:
		metadata:
			labels:
				app: nginx
		spec:
			containers:
			- images: nginx:latest
				livenessProbe: # 探活
					httpGet:
						path: /
						port: 80
					initialDelaySeconds: 15
					timeoutSeconds: 1
				readinessProbe:
					httpGet:
						path: /ping
						port: 80
					initialDelaySeconds: 5
					timeoutSeconds: 1
				resources: # 資源限制
					limits:
					cpu: "500m"
					memory: "500Mi"
				volumeMounts:
				- name: data
				mountPath: /data
			volumes:
			- name: data
				emptyDit:{} # temp dir

# 資源限制也可以透過以下指令
$ kubectl set resources deployment hello-nginx -c=nginx --limis=cpu=500m, memory=128Mi
```

### ConfigMap

- 將`非機密數據`保存到kv-pairs中
- Pods可以將此用作`環境變數`、`命令參數`，或是`存儲卷中的配置文件`
- 此目的為將環境訊息跟鏡像解耦，便於應用配置更改

### Secret

- 保存和傳遞`key`、`憑證`等敏感訊息的對象
- 避免把敏感訊息直接明文寫在配置文件
- kubernetes集群中配置和使用服務難免要用到敏感訊息進行登陸、認證等功能，在配置文件中透過`secret`對象使用這些敏感訊息，來避免重覆，以及減少暴露機會

### User Account & Service Account

- user account：提供帳戶標示
- service account：為Pod提供帳戶標示
- 兩者區別為作用範圍
    - `user account`對應的是人的身分，與服務的Namespace無關，所以是跨Namespace的
    - `service account`對應的是一個運行中程序的身分，與特定Namespace相關

### Service

- 服務應用的抽象，通過`labels`提供`服務發現`和`負載均衡`
- 將與`labels`對應的`Pod IP`和`端口`組成`endpoints`，透過`Kube-proxy` 將服務負載均衡到這些`endpoints`上
- 每個`service`，會`自動分配一個只能在集群內部訪問的虛擬位址`（cluster ip）和 DNS，其他容器則透過該地址訪問服務

### Replica Set

- 提供高可用應用，構建多個同樣Pod的副本，為同一服務
- 每一個pod為一個無狀態模式進行管理
- 副本掛掉，`Controller`會自動重新創建一個新副本
- 為負載發生變更時，方便調整擴縮容策略

### Deployment

- 為集群中的一次更新操作
- 應用模式廣泛：創建新服、更新服務，甚至滾動升級服務
    - `滾動升級`服務，本質上是創建一個新的RS，然後將新RS中副本數加到想要的狀態，並將舊RS縮減到0的一個符合操作
- 目前管理應用方式，皆為採用此來管理

### StatefulSet

- 用來管理有狀態應用的工作負載API對象
- statefulset中的每個Pod都掛載自己獨立的存儲
- `如果一個Pod有問題，從其他節點啟動一個同樣名稱的Pod，還需要掛載原Pod的存儲，繼續使用原本的狀態提供服務`
- 適合statefulset的業務，為有狀態之服務
    - 數據庫服務：MySQL、PostgreSQL
    - 集群化管理服務：ZooKeeper、etcd
- statefulset提供的是把`特定Pod和特定存儲關聯起來`，保證狀態的延續性，所以Pod仍可以通過飄移到不同節點提供高可用，而存儲也可以通過外掛存儲提供高可靠性
- 跟`Deployment`差異
    - 身分標示
        - statefulset controller為每個Pod提供編號，編號從0開始
    - 數據存儲
        - statefulset 可以讓用戶定義`PVC`的 `volumeClaimTemplates`
        - Pod被創建時，kubernetes會以`volumeClaimTemplates`中定義的模板創建存儲卷，並掛載給Pod
        
        ```yaml
        apiVersion: apps/v1
        kind: StatefulSet
        metadata:
          name: web
        spec:
          selector:
            matchLabels:
              app: nginx # 需匹配 .spec.template.metadata.labels
          serviceName: "nginx"
          replicas: 3 # default：1
          minReadySeconds: 10  #default：0
          template:
            metadata:
              labels:
                app: nginx # 需匹配 .spec.selector.matchLabels
            spec:
              terminationGracePeriodSeconds: 10
              containers:
              - name: nginx
                image: registry.k8s.io/nginx-slim:0.8
                ports:
                - containerPort: 80
                  name: web
                volumeMounts:
                - name: www
                  mountPath: /usr/share/nginx/html
          volumeClaimTemplates:
          - metadata:
              name: www
            spec:
              accessModes: [ "ReadWriteOnce" ]
              storageClassName: "my-storage-class"
              resources:
                requests:
                  storage: 1Gi
        ```
        
    - Statefulset升級策略不同
        - onDelete
        - 滾動升級
        - 分片升級

### Job

- 用來控制批次處理任務的API對象
- 任務完成即自動退出
- 成功完成的標誌根據不同的`.spec.completions`策略而不同
    - 單Pod型任務有一個Pod成功就標示完成
    - 定數成功型任務保證有N個任務全部成功
    - 工作隊列型任務根據應用確認的全部成功才標示成功

### DaemonSet

- 守護進程，保證每個節點都有一個此類Pod運行
- 節點可能是集群節點，或是通過`nodeSelector`選的特定節點
- 經典用法
    - 節點上運行集群守護進程
    - 節點上運行日誌收集守護進程
    - 節點上運行監控守護進程

### PV & PVC

- PersistentVolume （PV）
    - 集群中的一個存儲卷
    - 可以管理員手動建立，或是用戶建立`PVC`時根據`StorageClass`動態設置
- PersistentVolumeClaim （PVC）
    - 為用戶對存儲的請求
    - 每個PVC對象都有spec和status，分別對應申請部分和狀態
    - `CSI external-provisioner`是一個監控Kubernetes PVC對象的`Sidecar`容器
        - 當用戶創建PVC後，Kubernetes會檢查PVC對應的`StorageClass`，如果`SC`中的`provioner`與某個插件匹配，該容器通過`CSI Endpoint`（通常是unix socket）調用`CreateVolume`方法，調用成功則創建`PV`對象
- StorageClass
    - 常見不同情境，用戶需要具有不同屬性（如：`性能`、`訪問模式`）的PV
    - 集群管理員需要提供不同性質的PV，且這些PV卷之間的差別不僅限於`卷大小`和`訪問模式`，同時又無法將實現細節暴露給用戶，所以又誕生了StorageClass

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: myclaim
spec:
  accessModes:
    - ReadWriteOnce
  volumeMode: Filesystem
  resources:
    requests:
      storage: 8Gi
  storageClassName: slow
  selector:
    matchLabels:
      release: "stable"
    matchExpressions:
      - {key: environment, operator: In, values: [dev]}
# ----
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: standard
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
reclaimPolicy: Retain
allowVolumeExpansion: true
mountOptions:
  - debug
volumeBindingMode: Immediate
```

### CRD

> ref： 
- [https://kubernetes.io/zh-cn/docs/concepts/extend-kubernetes/api-extension/custom-resources/](https://kubernetes.io/zh-cn/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [https://github.com/kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)
> 
- CustomResourceDefinition （CRD）
    - 允許用戶自訂義`Schema`
- 用戶可以基於CRD定義一切需要的模型，來滿足不同業務需求
- 基於CRD還推出了 `Operator Mode` 和 `Operator SDK`，以極低的開發成本定義新對象，並構建新對象的`Controller`
- 眾多主流的擴展應用都是基於CRD構建的
    - e.g. Istio

```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: foos.samplecontroller.k8s.io
  # for more information on the below annotation, please see
  # https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/2337-k8s.io-group-protection/README.md
  annotations:
    "api-approved.kubernetes.io": "unapproved, experimental-only; please get an approval from Kubernetes API reviewers if you're trying to develop a CRD in the *.k8s.io or *.kubernetes.io groups"
spec:
  group: samplecontroller.k8s.io
  versions:
    - name: v1alpha1
      served: true
      storage: true
      schema:
        # schema used for validation
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                deploymentName:
                  type: string
                replicas:
                  type: integer
                  minimum: 1
                  maximum: 10
            status:
              type: object
              properties:
                availableReplicas:
                  type: integer
      # subresources for the custom resource
      subresources:
        # enables the status subresource
        status: {}
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
# -----
# create a custom resource of type Foo
apiVersion: samplecontroller.k8s.io/v1alpha1
kind: Foo
metadata:
  name: example-foo
spec:
  deploymentName: example-foo
  replicas: 1
```

## Kube-apiserver

- Kubernetes最重要的核心組件之一
- 提供集群管理的API接口
- 提供認證、授權、准入
- 提供其他模塊之間的數據交互和通訊的橋樑
    - 其他模塊通過`API Server`查詢或修改數據
    - 只有`API Server`能直接操作`etcd`
- 每個請求都會經過多個階段訪問才會被接受

![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled.png)

### 認證

- 支援多種認證機制，並支持同時開啟多認證插件，只要有一個認證通過即可
    - 認證插件
        - x509
            - 使用x509客戶端生成證書，API Server啟動時配置 `—client-ca-file`。證書認證時，其中`CN為用戶名，機構為group名`
              
                ```go
                csr, cert/key, CA
                - create csr 證書請求
                - 生成csr同時，會有private key生成
                - 用csr去CA機構申請證書
                - CA機構給你一個cert，格式crt、pem
                ```
            
        - 靜態token
            - API Server啟動時配置 `—token-auth-file`，文件格式為`csv`
                - row： token, user, uid, “group1,group2,group3”
        - 靜態密碼文件
            - API Server啟動時配置 `—basic-auth-file`，文件格式為`csv`
                - row： password, user, uid, “group1,group2,group3”
        - ServiceAccount
            - Kubernetes自動生成的，並自動掛載到容器的`/run/secrets/kubernetes.io/serviceaccount`目錄中
        - OpenID
            - OAuth 2.0認證機制
        - Webhook令牌身分認證
            - `—authentication-token-webhook-config-file`
                - 指定配置文件，描述如何訪問webhook服務
            - `—authentication-token-webhook-cache-ttl`
                - 設定身分認證決定的緩存時間，默認2分鐘
        - 匿名請求
            - 如果使用`AlwaysAllowy`以外的認證模式，預設開啟，`—anonymous-auth=false` 則禁用匿名請求
- 課後練習
    - 採用`Gitlab api` + `k8s-v1 api`建立k8s webhook授權後台
    
    ![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%201.png)
    
    ```go
    //轉發認證請求
    decoder := json.NewDecoder(r.Body)
    		var tr authentication.TokenReview
    		err := decoder.Decode(&tr)
    		if err != nil {
    			log.Println("[Error]", err.Error())
    			w.WriteHeader(http.StatusBadRequest)
    			err := json.NewEncoder(w).Encode(map[string]interface{}{
    				"apiVersion": "authentication.k8s.io/v1beta1",
    				"kind":       "TokenReview",
    				"status": authentication.TokenReviewStatus{
    					Authenticated: false,
    				},
    			})
    			if err != nil {
    				return
    			}
    			return
    		}
    // 結果返回
    w.WriteHeader(http.StatusOK)
    		trs := authentication.TokenReviewStatus{
    			Authenticated: true,
    			User: authentication.UserInfo{
    				Username: resp.Username,
    				UID:      resp.UID,
    			},
    		}
    		json.NewEncoder(w).Encode(map[string]interface{}{
    			"apiVersion": "authentication.k8s.io/v1beta1",
    			"kind":       "TokenReview",
    			"status":     trs,
    		})
    // Demo
    $ curl http://localhost:3000/authenticate
    {"apiVersion":"authentication.k8s.io/v1beta1","kind":"TokenReview","status":{"authenticated":true,"user":{"username":"JWang10","uid":<uuid>}}}
    ```
    

### 授權

- 識別用戶是否有相對應的操作權限
- 通過組合屬性（用戶屬性、資源屬性、實體）策略向用戶授予訪問權限
- 授權方式
    - ABAC
        - 需要對Master節點的SSH的文件系統權限，使得授權變更成功需要重啟API Server，較於難管理
    - RBAC
        - 典型權限管理模型
        - Role / ClusterRole ： 權限集合
        - Rolebinding / ClusterRolebinding：將角色中定義的權限賦予用戶
    - Webhook
    - Node
- 高版本kubernetes默認授權方式為RBAC、Node

```go
root@kubemaster:~# cat role-namespace-admin.yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: namespace-admin
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
// ---
root@kubemaster:~# cat role-namespace-user.yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: namespace-user
rules:
- apiGroups: [""]
  resources: ["pods", "pods/log"]
  verbs: ["get", "list"]
- apiGroups: ["", "apps"]
  resources: ["deployments", "replicasets", "pods"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
// ----
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: namespace-admin-rolebinding
  namespace: default
subjects:
- kind: ServiceAccount
  name: test
roleRef:
  kind: Role
  name: namespace-admin // namespace-user
  apiGroup: rbac.authorization.k8s.io

```

### 准入

- 因資源有限，限制用戶使用資源
- 用戶可以根據自身需求定義准入插件來管理集群
    - mutating ：更改型
        - 向Pod中注入sidecar，[https://istio.io/v1.1/help/ops/setup/injection/](https://istio.io/v1.1/help/ops/setup/injection/)
        - 修改Pod spec中的部分字段，在label或annotation強制`加入創建pod的來源IP和請求者`
        - 修改pod spec，強制限制`cpu/ memory request`
    - validating：驗證型
        - 檢查鏡像url不是對應的倉庫地址，拒絕
        - 不允許CPU request超過10%的pod spec
    
    ![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%202.png)
    
- Alwayspullimages
    - 多租戶集群是有用的，強制鏡像需被拉去，而拉去鏡像就需要憑證
- Webhook
    - 如果需要大量定製`admission`，採用webhook方式，不需改源碼，減少維護成本
- NodeRestriction
    - 限制每一個kubelet只能修改的`Node`和`Pod`對象
    - 集群部署沒有設定`NodeRestriction`，用戶可以通過攻克一個子節點獲得kubelet的證書和密鑰，進而控制整個集群資源

### 限流

- MaxInFlightLimit，集群級別限流
    - `--max-requests-inflight` 和`--max-mutating-requests-inflight`
- Client限流
    - 代碼處理
    - [https://github.com/kubernetes/client-go/blob/master/util/workqueue/default_rate_limiters.go#L43](https://github.com/kubernetes/client-go/blob/master/util/workqueue/default_rate_limiters.go#L43)
- EventRateLimit
    - 只限制event
- APF，細粒度限制配置
    - 請求通過`FlowSchema`照其屬性分類，並分配優先級
    - [https://kubernetes.io/zh/docs/concepts/cluster-administration/flow-control/](https://kubernetes.io/zh/docs/concepts/cluster-administration/flow-control/)

### 高可用API Server

- Etcd集群高可用性
- API Server物理可用性
    - 多個server
    - 利用LB做負載均衡
    - F5做負載均衡
    - KeepLive+LVS
- 防止server掛掉
    - 外部controller
    - 限流
    - 對controller審計
    - `master`節點一定`禁止非kube-xxx`的pod

## Kube Scheduler

> ref:
[https://kubernetes.io/zh/docs/concepts/scheduling-eviction/scheduling-framework/](https://kubernetes.io/zh/docs/concepts/scheduling-eviction/scheduling-framework/)
> 
- 監聽`kube-apiserver`，檢查還沒分配`Node`的`Pod`，根據調度策略為這些Pod分配最合適節點，也就是`更新Pod的NodeName字段`
- 需要考慮因素
    - 公平調度
    - 資源利用
    - QoS
    - 親和性（affinity）和反親和性（anti-affinity）
    - 數據本地化
    - 內部負載干擾
    - deadlines
- 調度分為兩個階段：predicate和priority
    - predicate：過濾不符合條件的節點
    - priority：優先級排序
- 調度框架
    - 現在調度器已把所有算法整合到框架中
        - [https://github.com/kubernetes/enhancements/tree/master/keps/sig-scheduling/624-scheduling-framework](https://github.com/kubernetes/enhancements/tree/master/keps/sig-scheduling/624-scheduling-framework)
    - 提供更多自定位位點和可擴展性
    - 簡化調度器核心代碼，把部分實現轉移的`plugin`中
    - 提供一種高效機制，確認`plugins`的結果或使用`plugins`的結果
    - 支持`out-of-tree`擴展等

## Kube Controller Manager

- 種類及管理
    - [https://github.com/kubernetes/kubernetes/blob/master/cmd/kube-controller-manager/app/controllermanager.go#L408](https://github.com/kubernetes/kubernetes/blob/master/cmd/kube-controller-manager/app/controllermanager.go#L408)
- Leader Election
    - 提供Leader選舉機制，確保多個`Controller`實例同時運行，且只有`Leader`實例提供真正的服務，其他則處於就緒狀態，防止`Leader`出現故障，還能保證`Pod`能被即時調度，犧牲更多資源提升`Controller`可用性。
    - 本質是利用kubernetes中的`configmap`、`endpoint`或是`lease`資源實現一個分佈式鎖，拿到鎖的節點為leader，且定期`renew`。**當leader掛掉後，租約到期，其他節點就能成為新leader**
    
    ![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%203.png)
    
- Controller
    - `kube-controller-manager`將會啟動多個controller服務
    - `deployment-controller-manager`會啟動多個`informer`進行資源監聽
    - 當`api-server`將`deployment數據`存入到`etcd`後，`controller-manager`通過`reflector`對數據進行監聽，監聽到事件後將數據存入`DeltaFIFO`中，也會存入到自己的緩存中。 `informer`通過消費`DeltaFIFO`，將資源數據存入`indexer`中，同時將事件進行通知，由`controller`接受到通知後，將該事件發送到`workerqueue`中。而`workerqueue`中的數據如何進行處理，則是由`controller-manager`來控制
    - 代碼利用工廠模式創建各類`informer`，簡化了s實例創建
    
    ![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%204.png)
    
    - 常見內部Controller
        - [https://github.com/kubernetes/kubernetes/tree/master/pkg/controller](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller)

## Kubelet

![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%205.png)

- 每個節點上都運行一個`kubelet`服務進程，默認`10250`端口
- CRI與底層容器交互
- 管理Pod及其中的容器
- 通過`cAdvisor`監控節點和容器資源
- Managers
    - 存儲
    - 設備
    - 文件
    - 狀態
- GC
    - images GC
    - contrainer GC
- 監控接口
- syncloop
    - 檢測實際環境中container出現的變化，其每秒鐘列出列出當前節點所有pod和所有container，與自己緩存中podRecord中對比，生成每個container的event，送到`event channel`，kubelet主循環`syncLoop`負責處理`event channel`中的事件
    
    ![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%206.png)
    

### Pod lifecycle

- 啟動流程

![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%207.png)

![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%208.png)

![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%209.png)

### CRI

- Container Runtime Interface，CR運行於每個節點中，負責容器的整個生命週期
- kubernetes定義的一組gRPC服務，基於gRPC框架，通過Socket和CR通訊
- 包含`鏡像服務`（Image Service）及`運行時服務`（Runtime Service）
    - Image Service：提供鏡像調用服務
    - Runtime Service：管理容器生命週期和容器交互的調用
- Docker內部容器運行時功能的核心組件是`containerd`，後來`containerd`可直接跟kubelet通過CRI對接，獨立於Kubernetes中使用
  
    ![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%2010.png)
    
    ![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%2011.png)
    

### CNI

- Container Network Interface，用來設置和刪除容器的網路連通
- IP地址是以Pod為單位分配的，每個Pod有獨立的地址
- Pod裏面全容器都可以透過`localhost:port`來連結
- 由kubelet查找CHI插件來為容器設置網路
    - cni-bin-dir：執行文件目錄，默認`/opt/cni/bin`
    - cni-conf-dir：配置未見目錄，默認`/etc/cni/net.d`
- 常見CNI插件
    - flannel
    - calico
    - cilium

### CSI

- Container Storage Interface
- 目前`docker`和`containerd`都默認以`OverlayFS`作為運行時存儲驅動
- Kubernetes支持以插件形式來實現對不同存儲的擴展
- 分兩種類型的支持：
    - in-tree：在kubernetes內部代碼上支持
    - out-of-tree：通過接口支持的
- Kubernets存儲
    - 非持久化存儲主要是  `emptydir`，用於緩存、臨時儲存
    - 非 `emptydir` 的基本都是持久存儲
        - HostPath
        - StorageClas
        - PV
        - PVC

![Untitled](%E7%B8%BD%E7%B5%90%E7%AD%86%E8%A8%98%20212735cf19f54e50856542c1d0a1387f/Untitled%2012.png)