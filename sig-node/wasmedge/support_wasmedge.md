


````markdown
---
title: KubeEdge and WasmEdge Integration
sidebar_position: 4
---

WasmEdge is a lightweight, high-performance, and extensible **WebAssembly runtime** for cloud native. It is compliant with OCI standard and is currently a **CNCF (Cloud Native Computing Foundation) Sandbox project**. It powers edge computing, serverless apps, microservices, service mesh, and IoT devices.

In this article, we will introduce how to run a WasmEdge simple demo app with **containerd** over **KubeEdge**. [KubeEdge+WasmEdge+CRI-O](https://wasmedge.org/book/en/use_cases/kubernetes/kubernetes/kubeedge.html) also provides the operation steps about running a WasmEdge demo app with CRI-O over KubeEdge.

We assume that you have completed the installation of the Kubernetes cluster. [Creating a cluster with kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/) will help you create cluster with kubeadm.

---

# Install Containerd And WasmEdge on Edge Node

WasmEdge provides the one-click installation scripts. You can directly execute:

```bash
wget -qO- [https://raw.githubusercontent.com/second-state/wasmedge-containers-examples/main/containerd/install.sh](https://raw.githubusercontent.com/second-state/wasmedge-containers-examples/main/containerd/install.sh) | bash
````

to complete the installation of containerd, WasmEdge, and crun. See [installation scripts](https://github.com/second-state/wasmedge-containers-examples/blob/main/containerd/install.sh) to learn more details.

Next we will introduce the process step by step.

## 1\. Install containerd

  * See [quick start with containerd](https://github.com/containerd/containerd/blob/main/docs/getting-started.md), download the official binary and install it.

  * Configure containerd:

  ` bash   sudo mkdir -p /etc/containerd   sudo bash -c "containerd config default > /etc/containerd/config.toml"    `

## 2\. Install WasmEdge

Use the [install script](https://github.com/WasmEdge/WasmEdge/blob/master/utils/install.sh) to install WasmEdge on your edge node.

```bash
curl -sSf [https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh](https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh) | bash
```

## 3\. Install crun and modify containerd config

The **crun** project has WasmEdge support baked in. Next we need to get the crun binary and install it.

  * Configure the Dependencies Required

<!-- end list -->

```bash
sudo apt update
sudo apt install -y make git gcc build-essential pkgconf libtool \
    libsystemd-dev libprotobuf-c-dev libcap-dev libseccomp-dev libyajl-dev \
    go-md2man libtool autoconf python3 automake
```

  * Configure, build, and install a crun binary with WasmEdge support.

<!-- end list -->

```bash
git clone [https://github.com/containers/crun](https://github.com/containers/crun)
cd crun
./autogen.sh
./configure --with-wasmedge
make
sudo make install
```

  * Modify containerd config and start containerd.

It mainly includes two modifications:
(1) Configure containerd to use **`crun`** as the low-level OCI runtime;
(2) Set pod\_annotations: `pod_annotations = ["module.wasm.image/ variable .*"]`.

```bash
wget [https://raw.githubusercontent.com/second-state/wasmedge-containers-examples/main/containerd/containerd_config.diff](https://raw.githubusercontent.com/second-state/wasmedge-containers-examples/main/containerd/containerd_config.diff)
sudo patch -d/ -p0 < containerd_config.diff
sudo systemctl start containerd
```

---

# Install KubeEdge

After completing the above work, you can install KubeEdge according to the official website [Installation Guide](https://release-1-21.docs.kubeedge.io/docs/setup/install-with-keadm), which mainly includes three steps.

For full compatibility with the latest KubeEdge features and the container runtime configuration, we will use a more recent stable version of KubeEdge (e.g., **v1.21.1**) for the installation examples.

## 1\. Install Keadm

```bash
# Get the latest stable release link from the KubeEdge GitHub releases page
KUBEEDGE_VERSION=v1.21.1
# FIX: Corrected variable interpolation in the wget command line.
wget [https://github.com/kubeedge/kubeedge/releases/download/$](https://github.com/kubeedge/kubeedge/releases/download/$){KUBEEDGE_VERSION}/keadm-${KUBEEDGE_VERSION}-linux-amd64.tar.gz
tar -zxvf keadm-${KUBEEDGE_VERSION}-linux-amd64.tar.gz
cp keadm-${KUBEEDGE_VERSION}-linux-amd64/keadm /usr/local/bin/keadm
```

## 2\. Setup cloud side (CloudCore)

````bash
# Use a compatible KubeEdge version and specify the profile
KUBEEDGE_VERSION=v1.21.1
keadm init --advertise-address="THE-EXPOSED-IP" --profile version=${KUBEEDGE_VERSION} --kube-config=/root/.kube/config
# Note: --profile expects a version string (e.g., version=v1.21.1)
``` 

After installing CloudCore, run `keadm gettoken` in the cloud side to retrieve the token, which will be used when joining edge nodes.

## 3. Install EdgeCore and join edge node

KubeEdge **v1.21.1** (or later versions) uses containerd as the default container runtime. You don't need to specify the runtime type when installing EdgeCore.

```bash
KUBEEDGE_VERSION=v1.21.1
# Replace $token with the actual token obtained from 'keadm gettoken'
keadm join --cloudcore-ipport="THE-EXPOSED-IP":10000 --token=${token} --kubeedge-version=${KUBEEDGE_VERSION}
````

---

# Enable Kubelet Logs Feature

Since the WebAssembly application used in the demo needs to enable the Kubelet logs feature, please see [Advanced Debugging](https://kubeedge.io/docs/advanced/debug) to learn how to turn on it.

---

# Run a Simple WebAssembly App

We can run the WebAssembly-based image in the Kubernetes cluster. The annotations ensure that the WasmEdge runtime is correctly used.

```bash
kubectl run -it --restart=Never wasi-demo --image=hydai/wasm-wasi-example:with-wasm-annotation --annotations="module.wasm.image/variant=compat-smart" /wasi_example_main.wasm 50000000
```

The output in the terminal demonstrates the successful execution of the WebAssembly application:

```
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

According to the pod log, the WebAssembly app of the pod successfully deployed to the edge node.

```bash
kubectl describe pod wasi-demo
```

The output confirms the pod's status, indicating it ran on the specified node and completed successfully:

```
Name:          wasi-demo
Namespace:     default
Priority:      0
Node:          kubeedge-dev-linux-baoyue/192.168.1.123
Start Time:    Mon, 06 Mar 2023 19:41:42 +0800
Labels:        run=wasi-demo
Annotations:   module.wasm.image/variant: compat-smart
Status:        Succeeded
IP:            10.88.0.13
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
    State:           Terminated
      Reason:        Completed
      Exit Code:     0
      Started:       Mon, 06 Mar 2023 19:41:50 +0800
      Finished:      Mon, 06 Mar 2023 19:41:50 +0800
    Ready:           False
    Restart Count:   0
    Environment:     <none>
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

It can also be seen on the edge node that the container successfully completed and exited:

```bash
crictl ps -a
```

```
CONTAINER             IMAGE                          CREATED               STATE                 NAME                   ATTEMPT               POD ID
83aafc8745f86         0423b8eb71e31                  7 seconds ago         Exited                wasi-demo              0                     b23754a9ee905
```

```
```
