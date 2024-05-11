WasmEdge是一个由CNCF托管、符合OCI标准的云原生WebAssembly运行时，拥有高性能、更轻量、更安全等优势，目前广泛应用于边缘计算、Serverless应用、微服务、服务网格和IoT设备等技术。

本篇文章将介绍如何在KubeEdge的边缘节点中，使用containerd运行一个WasmEdge的demo。[KubeEdge+WasmEdge+CRI-O](https://wasmedge.org/book/zh/kubernetes/kubernetes/kubeedge.html)同样提供了使用CRI-O运行WasmEdge的操作步骤。

本文假设您已经完成K8s集群的安装，可参阅[使用kubeadm创建集群](https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/)。

# 在边缘节点上安装containerd和WasmEdge

WasmEdge已经提供了一键安装脚本，直接执行：
```
wget -qO- https://raw.githubusercontent.com/second-state/wasmedge-containers-examples/main/containerd/install.sh | bash
````
即可完成containerd、WasmEdge和crun的一键安装，可以参阅[安装脚本](https://github.com/second-state/wasmedge-containers-examples/blob/main/containerd/install.sh)。

本文会逐步介绍安装过程。

## 1. 安装containerd

* 参阅[开始使用containerd](https://github.com/containerd/containerd/blob/main/docs/getting-started.md)，下载官方二进制文件，并执行安装

* 配置containerd：

  ```
  sudo mkdir -p /etc/containerd
  sudo bash -c "containerd config default > /etc/containerd/config.toml"
  ```

## 2. 安装WasmEdge

使用WasmEdge提供的[安装脚本](https://github.com/WasmEdge/WasmEdge/blob/master/utils/install.sh)安装WasmEdge。

```
curl -sSf https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash
```

## 3. 安装crun并修改containerd配置

crun项目已经内置了对WasmEdge的支持，接下来我们需要获取一个支持WasmEdge的crun二进制文件，并完成安装与配置。

* 配置crun需要的依赖

```
$ sudo apt update
$ sudo apt install -y make git gcc build-essential pkgconf libtool \
      libsystemd-dev libprotobuf-c-dev libcap-dev libseccomp-dev libyajl-dev \
      go-md2man libtool autoconf python3 automake
```

* 配置、构建并安装一个支持WasmEdge的crun

```
git clone https://github.com/containers/crun
cd crun
./autogen.sh
./configure --with-wasmedge
make
sudo make install
```

* 修改containerd的配置，并启动containerd服务。

主要包括两处修改：
（1）将containerd配置为使用`crun`作为底层OCI runtime；
（2）设置pod_annotations来传播到OCI规范：`pod_annotations = ["module.wasm.image/ variable .*"]`。

```
wget https://raw.githubusercontent.com/second-state/wasmedge-containers-examples/main/containerd/containerd_config.diff
sudo patch -d/ -p0 < containerd_config.diff
sudo systemctl start containerd
```

# 安装kubeedge

在完成上述工作之后，即可按照KubeEdge官网[安装指南](https://kubeedge.io/en/docs/setup/keadm/)开始安装，主要包括三个步骤。

## 1. 安装Keadm

```
wget https://github.com/kubeedge/kubeedge/releases/download/v1.12.1/keadm-v1.12.1-linux-amd64.tar.gz
tar -zxvf keadm-v1.12.1-linux-amd64.tar.gz
cp keadm-v1.12.1-linux-amd64/keadm/keadm /usr/local/bin/keadm
```

## 2. 在云端节点安装CloudCore

```
keadm init --advertise-address="THE-EXPOSED-IP" --profile version=v1.13.0 --kube-config=/root/.kube/config
``` 

CloudCore安装完成后，使用`keadm gettoken`获取token，用于安装EdgeCore时使用。

## 3. 在边缘节点安装EdgeCore

由于KubeEdge v1.13.0默认使用容器运行时为containerd，因此在安装EdgeCore时可以不指定容器运行时。如果安装的是v1.13.0之前的版本，需要在执行keadm join命令时，使用`--runtimetype`和`--remote-runtime-endpoint`指定边缘节点需要的容器运行时。

```
keadm join --cloudcore-ipport="THE-EXPOSED-IP":10000 --token=$token} --kubeedge-version=v1.13.0
```

# 启用Kubectl日志功能

由于本次演示所使用的WebAssembly应用需要启动Kubectl logs功能，因此需要打开KubeEdge的Logs/Exec开关，具体步骤参考[Enable Kubectl logs Feature](https://kubeedge.io/zh/docs/setup/keadm/#enable-kubectl-logs-feature)

# 运行一个简单的WebAssembly应用

我们可以直接在云端运行基于WebAssembly的镜像。

```
$ kubectl run -it --restart=Never wasi-demo --image=hydai/wasm-wasi-example:with-wasm-annotation --annotations="module.wasm.image/variant=compat-smart" /wasi_example_main.wasm 50000000

Random number: 626588879
Random bytes: [175, 254, 19, 202, 67, 26, 244, 82, 225, 201, 104, 99, 152, 44, 222, 233, 182, 185, 95, 166, 130, 74, 36, 88, 88, 69, 141, 106, 155, 79, 80, 7, 91, 239, 112, 27, 182, 103, 49, 215, 171, 109, 80, 51, 190, 237, 166, 167, 87, 10, 235, 81, 159, 75, 22, 161, 94, 12, 97, 157, 216, 223, 41, 80, 5, 137, 124, 89, 158, 246, 1, 109, 20, 90, 125, 29, 236, 239, 238, 7, 195, 1, 244, 241, 226, 145, 118, 44, 235, 250, 225, 155, 210, 235, 137, 9, 194, 118, 72, 251, 113, 255, 164, 110, 94, 212, 150, 59, 228, 220, 164, 243, 68, 64, 77, 115, 124, 70, 201, 111, 73, 171, 27, 0, 225, 130, 80, 66]
Printed from wasi: This is from a main function
This is from a main function
The env vars are as follows.
PATH: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
TERM: xterm
HOSTNAME: wasi-demo
HOME: /
The args are as follows.
/wasi_example_main.wasm
50000000
File content is This is in a file
```

由上述pod日志可以看出WebAssembly应用的pod成功部署到边缘节点，容器成功执行完成并退出。
```
$ kubectl describe pod wasi-demo

Name:         wasi-demo
Namespace:    default
Priority:     0
Node:         kubeedge-dev-linux-baoyue/192.168.1.123
Start Time:   Mon, 06 Mar 2023 19:41:42 +0800
Labels:       run=wasi-demo
Annotations:  module.wasm.image/variant: compat-smart
Status:       Succeeded
IP:           10.88.0.13
IPs:
  IP:  10.88.0.13
  IP:  2001:4860:4860::d
Containers:
  wasi-demo:
    Container ID:  containerd://1e1df3ece5d3d67ead3375e82df039d617b9698e421be2a085eadbfe273b2a06
    Image:         hydai/wasm-wasi-example:with-wasm-annotation
    Image ID:      docker.io/hydai/wasm-wasi-example@sha256:525aab8d6ae8a317fd3e83cdac14b7883b92321c7bec72a545edf276bb2100d6
    Port:          <none>
    Host Port:     <none>
    Args:
      /wasi_example_main.wasm
      50000000
    State:          Terminated
      Reason:       Completed
      Exit Code:    0
      Started:      Mon, 06 Mar 2023 19:41:50 +0800
      Finished:     Mon, 06 Mar 2023 19:41:50 +0800
    Ready:          False
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-lznkt (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             False 
  ContainersReady   False 
  PodScheduled      True 
Volumes:
  kube-api-access-lznkt:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:                      <none>
```

在边缘节点也可以查看对应的容器状态。

```
crictl ps -a 

CONTAINER           IMAGE                      CREATED             STATE               NAME                ATTEMPT             POD ID
83aafc8745f86       0423b8eb71e31              7 seconds ago       Exited              wasi-demo           0                   b23754a9ee905
```