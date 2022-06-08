# Module3

### Cgroup 調用

- 在cgroup cpu子系統目錄創建目錄結構
```shell
cd /sys/fs/cgroup/cpu
mkdir cpudemo
```
- 執行 `busyloop`
- 使用 `top`查看`cpu`使用情況
- 通過`cgroup`限制`cpu`
  - 把進程加入到cpudemo進程配置
    ```shell
    cd /sys/fs/cgroup/cpu/cpudemo
    # 找到busyloop的PID，加入到cgroup.proc
    ps -ef | grep busyloop
    echo <pid> > cgourp.proc
    # 或是直接透過下面指令加入
    ps -ef | grep busyloop | grep -v grep | awk '{print $2}' > cgroup.proc
    ```
  - 設置`cpuquota` 
    ```shell
    #  cpu.cfs_period_us 用來配置CPU時間週期長度
    #  cpu.cfs_quota_us 當前cgroup在cpu.cfs_period_us中最多能使用的cpu時間數
    echo 10000 > cpu.cfs_quota_us    
    ```
  - 查看`cpu`使用量變為`10%`