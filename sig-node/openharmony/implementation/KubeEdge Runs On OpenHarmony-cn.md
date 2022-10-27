# RK3568 OpenHarmony运行KubeEdge



本教程主要分为了两个关键步骤：OpenHarmony运行Docker的步骤和OpenHarmony运行KubeEdge的步骤。同时以润和DAYU200开发板为例，展示在OpenHarmony上运行KubeEdge两个关键步骤的操作详情和参考文档。这里，KubeEdge所用的运行时为Docker。



## OpenHarmony运行Docker的关键步骤

##### 1. 准备支持Docker容器的OpenHarmony内核
- 修改内核配置：
cgroup和namespace相关特性，主要修改的文件是kernel/linux/config/linux-5.10/arch/arm64/configs/rk3568_standard_defconfig 
- 修改network内核配置：
网络主要用的是bridge模式，所以要打开相应的配置以及支持网络包forward功能。
- 修改overlay filesystem：
Docker使用的overlay filesystem推荐使用未加密的f2fs作为backing filesystem。

##### 2. 编译并且在设备上安装运行Openharmony

##### 3. 安装必要的工具
Docker需要用到的一些必要工具有iproutes、iptables和busybox。将这些工具的静态二进制文件拷贝到OpenHarmony设备/bin/目录下，并且添加执行权限。

##### 4. 安装Docker容器引擎组件
将Docker静态二进制文件拷贝到OpenHarmony设备 /bin/ 目录下，并添加执行权限

##### 5. 挂载OpenHarmony系统资源
创建OpenHarmony的/etc/cgroups.json文件，挂载所有的cgroup子系统。

##### 6. OpenHarmony环境设置
- 创建docker运行所需要的目录
```
# 创建相关目录
mkdir /var
mkdir /run
mkdir /tmp
mkdir /opt
mkdir /usr
mkdir /data/var
mkdir /data/run
mkdir /data/tmp
mkdir /data/opt
mkdir /data/etc
mkdir /data/etc/docker
mkdir /data/usr
mkdir /data/bin
mkdir /data/root
```
- 挂载docker所需要的目录
```
mount tmpfs /sys/fs/cgroup -t tmpfs -o size=1G
mkdir /sys/fs/cgroup/blkio
mkdir /sys/fs/cgroup/cpu
mkdir /sys/fs/cgroup/cpuacct
mkdir /sys/fs/cgroup/cpuset
mkdir /sys/fs/cgroup/devices
mkdir /sys/fs/cgroup/freezer
mkdir /sys/fs/cgroup/hugetlb
mkdir /sys/fs/cgroup/memory
mkdir /sys/fs/cgroup/net_cls
mkdir /sys/fs/cgroup/net_prio
mkdir /sys/fs/cgroup/perf_event
mkdir /sys/fs/cgroup/pids
mkdir /sys/fs/cgroup/systemd
# mount --bind
mount --bind /data/etc/docker /etc/docker
mount --bind /data/var /var
mount --bind /data/run /run
mount --bind /data/tmp /tmp
mount --bind /data/opt /opt
mount --bind /data/usr /usr
mount --bind /data/bin /bin
mount --bind /data/root /root
#mount cgroup
mount -t cgroup -o none,name=systemd cgroup /sys/fs/cgroup/systemd
mount -t cgroup -o blkio,nodev,noexec,nosuid cgroup /sys/fs/cgroup/blkio
mount -t cgroup -o cpu,nodev,noexec,nosuid cgroup /sys/fs/cgroup/cpu
mount -t cgroup -o cpuacct,nodev,noexec,nosuid cgroup /sys/fs/cgroup/cpuacct
mount -t cgroup -o cpuset,nodev,noexec,nosuid cgroup /sys/fs/cgroup/cpuset
mount -t cgroup -o devices,nodev,noexec,nosuid cgroup /sys/fs/cgroup/devices
mount -t cgroup -o freezer,nodev,noexec,nosuid cgroup /sys/fs/cgroup/freezer
mount -t cgroup -o hugetlb,nodev,noexec,nosuid cgroup /sys/fs/cgroup/hugetlb
mount -t cgroup -o memory,nodev,noexec,nosuid cgroup /sys/fs/cgroup/memory
mount -t cgroup -o net_cls,nodev,noexec,nosuid cgroup /sys/fs/cgroup/net_cls
mount -t cgroup -o net_prio,nodev,noexec,nosuid cgroup /sys/fs/cgroup/net_prio
mount -t cgroup -o perf_event,nodev,noexec,nosuid cgroup /sys/fs/cgroup/perf_event
mount -t cgroup -o pids,nodev,noexec,nosuid cgroup /sys/fs/cgroup/pids
```
- 添加路由规则
```
ip rule add pref 1 from all lookup main
ip rule add pref 2 from all lookup default
```
- 关闭selinux
```
setenforce 0
```
- 创建文件 /etc/docker/daemon.json 并且写入
```
{"registry-mirrors":["https://docker.mirrors.ustc.edu.cn"],"experimental":true}
```

##### 7. 运行docker
```
dockerd -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock & 
```

##### 8. 验证docker工作状态
```
# 确保OpenHarmony已经联网，执行以下命令会显示“ Hello from Docker! ”
docker run hello-world
```





## OpenHarmony运行KubeEdge的关键步骤

利用KubeEdge打通云端和OpenHarmony边端，实现云边协同。

##### 1. 云端安装KubeEdge cloudcore并且获得token
##### 2. KubeEdge edgecore 源码静态编译

在Arm server上编译KubeEdge源码获得edgecore静态二进制文件
```
# 1. 匹配OpenHamrony的memory.stat格式

# 这边其实是第三方库未能包含所有可能得memory.stat格式，这边我已提交issues：https://github.com/opencontainers/runc/issues/3643。

# 也可以修改/vendor/github.com/opencontainers/runc/libcontainer/cgroups/fscommon/utils.go

func ParseKeyValue(t string) (string, uint64, error) {

	tmp := strings.Replace(t, ":", "", -1)
	tmpo := strings.Replace(tmp, "\t", "", -1)
	tmpt := strings.Replace(tmpo, " kB", "", -1)
	count := strings.Count(tmpt, " ")
	tmps := strings.Replace(tmpt, " ", "", count-1)

	parts := strings.SplitN(tmps, " ", 3)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("line %q is not in key value format", t)
	}

	value, err := ParseUint(parts[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("unable to convert to uint64: %v", err)
	}

	return parts[0], value, nil
}


# 2. 由于openharmony安装docker用的是overlay2，所以需要修改kubeedge存放kubelet文件路径使其与docker overlay2的路径一样。

# 3. 编译edgecore

docker build -t kubeedge/edgecore:tag -f build/edge/Dockerfile .
docker cp $(docker create --rm kubeedge/edgecore:tag):/usr/local/bin/edgecore ./edgecore.tag

# 4. 在kubeedge目录下有edgecore.tag可执行文件拷贝到openharmony板子/bin/上
```

##### 3. OpenHarmony 边缘端安装edgecore

```
# 1. cpuset.mems文件添加“0”初始量
echo "0" > /dev/cpuset/background/cpuset.mems

# 2. 添加localhost的路径
echo "127.0.0.1 localhost localhost" > /etc/hosts

# 3. 添加edgecore.yaml
mkdir -p /etc/kubeedge/config
cd /etc/kubeedge/config
edgecore --minconfig > edgecore.yaml

# 4. 修改edgecore.yaml中的cloud-token
# 5. 启动edgecore程序即可
```




------





## 在润和DAYU200上实现参考




### 环境信息
- 开发板：润和DAYU200
- 芯片：RockChip RK3568（arm64）
- OpenHarmony版本：3.1 release
- Linux内核：5.10
- Docker：18.03.1
- KubeEdge：1.9.1 

（提醒：这里选择1.9.1，因为k8s与KubeEdge版本是有匹配的，有些高版本的k8s是不支持的。）



### 操作过程与参考文件

------------------------------------------------------------------
#### DockerOnOpenHarmony

##### 1. 原配置导出检测

 下载检测脚本
- [https://github.com/moby/moby/blob/master/contrib/check-config.sh](https://github.com/moby/moby/blob/master/contrib/check-config.sh)
- [https://gitee.com/mirrors_moby/moby/blob/master/contrib/check-config.sh#](https://gitee.com/mirrors_moby/moby/blob/master/contrib/check-config.sh#)
- 检测脚本就是check-config.sh

检测原始OpenHarmony3.1release内核对docker的支持
```
# 方式一、通过内核镜像直接导出config文件
# 找到源码编译后会产生内核img镜像并且记住路径
cd /out/kernel/src_tmp/linux-5.10/arch/arm/boot/

# 进入源码目录下
cd /out/kernel/src_tmp/linux-5.10/

# 执行extract-ikconfig脚本
scripts/extract-ikconfig   arch/arm/boot/*Image >  /home/.config

# 方式二、在OpenHarmony运行态导出config
cd /home
cat /proc/config.gz | gzip -d > /sdcard/.config
```

生成的.config就是内核的配置文件类似于下面：

```
# 配置文件的部分：
# Automatically generated file; DO NOT EDIT.
# Linux/arm64 5.10.79 Kernel Configuration
#
CONFIG_CC_VERSION_TEXT="aarch64-linux-gnu-gcc (Linaro GCC 7.5-2019.12) 7.5.0"
CONFIG_CC_IS_GCC=y
CONFIG_GCC_VERSION=70500
CONFIG_LD_VERSION=228020000
CONFIG_CLANG_VERSION=0
CONFIG_LLD_VERSION=0
CONFIG_CC_CAN_LINK=y
CONFIG_CC_CAN_LINK_STATIC=y
CONFIG_CC_HAS_ASM_GOTO=y
CONFIG_CC_HAS_ASM_INLINE=y
CONFIG_IRQ_WORK=y
CONFIG_BUILDTIME_TABLE_SORT=y
CONFIG_THREAD_INFO_IN_TASK=y

#
# General setup
#
CONFIG_INIT_ENV_ARG_LIMIT=32
# CONFIG_COMPILE_TEST is not set
CONFIG_LOCALVERSION=""
CONFIG_LOCALVERSION_AUTO=y
CONFIG_BUILD_SALT=""
CONFIG_DEFAULT_INIT=""
CONFIG_DEFAULT_HOSTNAME="localhost"
CONFIG_SWAP=y
CONFIG_SYSVIPC=y
CONFIG_SYSVIPC_SYSCTL=y
# CONFIG_POSIX_MQUEUE is not set
# CONFIG_WATCH_QUEUE is not set
CONFIG_CROSS_MEMORY_ATTACH=y
# CONFIG_USELIB is not set
CONFIG_AUDIT=y
CONFIG_HAVE_ARCH_AUDITSYSCALL=y
CONFIG_AUDITSYSCALL=y

#
# IRQ subsystem
#
CONFIG_GENERIC_IRQ_PROBE=y
CONFIG_GENERIC_IRQ_SHOW=y
CONFIG_GENERIC_IRQ_SHOW_LEVEL=y
CONFIG_GENERIC_IRQ_EFFECTIVE_AFF_MASK=y
CONFIG_GENERIC_IRQ_MIGRATION=y
CONFIG_HARDIRQS_SW_RESEND=y
CONFIG_GENERIC_IRQ_CHIP=y
CONFIG_IRQ_DOMAIN=y
CONFIG_IRQ_DOMAIN_HIERARCHY=y
CONFIG_GENERIC_IRQ_IPI=y
CONFIG_GENERIC_MSI_IRQ=y
CONFIG_GENERIC_MSI_IRQ_DOMAIN=y
CONFIG_IRQ_MSI_IOMMU=y
CONFIG_HANDLE_DOMAIN_IRQ=y
CONFIG_IRQ_FORCED_THREADING=y
CONFIG_SPARSE_IRQ=y
# CONFIG_GENERIC_IRQ_DEBUGFS is not set
# end of IRQ subsystem

CONFIG_GENERIC_IRQ_MULTI_HANDLER=y
CONFIG_GENERIC_TIME_VSYSCALL=y
CONFIG_GENERIC_CLOCKEVENTS=y
CONFIG_ARCH_HAS_TICK_BROADCAST=y
CONFIG_GENERIC_CLOCKEVENTS_BROADCAST=y
```
- 执行脚本
  ./check-config.sh    /home/.config

- 检测结果如图

![](image/check-1.png)

​                  ![](image/check-2.png) 

![](image/check-3.png)



##### 2. 修改源码内核配置

依据上述check的结果
- 需要将必选/可选的都打开
- Network Drivers部分：我们一般用bridge就可以了，暂不用再开其它了。
- Storage Drivers部分: 打开overlay部分就可以了。

修改内核配置
- 修改openharmony3.1/kernel/linux/config/linux-5.10/arch/arm64/configs/rk3568_standard_defconfig 
（主要集中在namespace, control group, network, overlay filesystem等方面支持的追加）

增加如下配置：
```
--- configbk	2022-09-19 14:54:36.873667819 +0800
+++ /home/cpss/bj/project/openharmony31/kernel/linux/config/linux-5.10/arch/arm64/configs/rk3568_standard_defconfig	2022-09-03 15:31:29.170980918 +0800
@@ -30,6 +30,7 @@

+CONFIG_POSIX_MQUEUE=y
+CONFIG_SCHED_WALT=y
+CONFIG_PSI=y
+CONFIG_PAGE_COUNTER=y
+CONFIG_CGROUP_BPF=y
+CONFIG_MEMCG_KMEM=y
+CONFIG_MEMCG_SWAP_ENABLED=y
+CONFIG_BLK_CGROUP=y
+CONFIG_BLK_DEV_THROTTLING=y
+CONFIG_RT_GROUP_SCHED=y
+CONFIG_CGROUP_PIDS=y
+CONFIG_CGROUP_HUGETLB=y
+CONFIG_CGROUP_PERF=y
+CONFIG_NET_CLS_CGROUP=y
+CONFIG_BPF_SYSCALL=y

+CONFIG_BINFMT_MISC=y
+CONFIG_TLS=y
+CONFIG_IP_MULTIPLE_TABLES=y
+CONFIG_IP_MROUTE_MULTIPLE_TABLES=y
+CONFIG_INET_ESP=y
+CONFIG_IPV6_MIP6=y
+CONFIG_IPV6_MULTIPLE_TABLES=y
+CONFIG_IPV6_MROUTE=y
+CONFIG_IPV6_MROUTE_MULTIPLE_TABLES=y
+CONFIG_NF_CONNTRACK=y
+CONFIG_NETFILTER_XT_MARK=y
+CONFIG_NETFILTER_XT_SET=y
+CONFIG_NETFILTER_XT_TARGET_CHECKSUM=y
+CONFIG_NETFILTER_XT_MATCH_ADDRTYPE=y
+CONFIG_NETFILTER_XT_MATCH_CONNTRACK=y
+CONFIG_NETFILTER_XT_MATCH_IPVS=y
+CONFIG_NETFILTER_XT_MATCH_CGROUP=y
+CONFIG_IP_SET=y
+CONFIG_IP_SET_HASH_IP=y
+CONFIG_IP_SET_HASH_NET=y
+CONFIG_IP_VS=y
+CONFIG_IP_VS_NFCT=y
+CONFIG_IP_VS_PROTO_TCP=y
+CONFIG_IP_VS_PROTO_UDP=y
+CONFIG_IP_VS_RR=y
+CONFIG_IP_VS_WRR=y
+CONFIG_IP_VS_SH=y

+CONFIG_IP_NF_MATCH_STATE=y
+CONFIG_IP_NF_MATCH_LIMIT=y
+CONFIG_IP_NF_TARGET_LOG=y
+CONFIG_NF_NAT=y
+CONFIG_IP_NF_FTP=y
+CONFIG_IP_NF_TARGET_REDIRECT=y
+CONFIG_IP_NF_CONNTRACK=y
+CONFIG_IP_NF_IRC=y
+CONFIG_IP_NF_NAT=y
+CONFIG_IP_NF_FILTER=y
+CONFIG_IP_NF_TARGET_MASQUERADE=y
+CONFIG_BRIDGE=y
+CONFIG_BRIDGE_NETFILTER=y
+CONFIG_CGROUP_NET_PRIO=y
+CONFIG_STREAM_PARSER=y
+CONFIG_DRIVERS_HDF_LIGHT=y
+CONFIG_HYPERHOLD=y
+CONFIG_HYPERHOLD_DEBUG=y
+CONFIG_HYPERHOLD_ZSWAPD=y
+CONFIG_HYPERHOLD_FILE_LRU=y
+CONFIG_HYPERHOLD_MEMCG=y
+CONFIG_ZRAM_GROUP=y
+CONFIG_ZRAM_GROUP_DEBUG=y
+CONFIG_ZLIST_DEBUG=y
+CONFIG_ZRAM_GROUP_WRITEBACK=y
+CONFIG_REGMAP_SPI=y
+CONFIG_MACVLAN=y
+CONFIG_VXLAN=y
+CONFIG_AUFS_FS=y
+CONFIG_VETH=y
+CONFIG_DRM_DW_HDMI_I2S_AUDIO=y
+CONFIG_SND_TIMER=y
+CONFIG_SND_PCM=y
+CONFIG_SND_PCM_ELD=y
+CONFIG_SND_PCM_IEC958=y
+CONFIG_SND_DMAENGINE_PCM=y
+CONFIG_SND_HWDEP=y
+CONFIG_SND_SEQ_DEVICE=y
+CONFIG_SND_RAWMIDI=y
+CONFIG_SND_JACK=y
+CONFIG_SND_JACK_INPUT_DEV=y
+CONFIG_SND_PCM_TIMER=y
+CONFIG_SND_HRTIMER=y
+CONFIG_SND_DYNAMIC_MINORS=y
+CONFIG_SND_MAX_CARDS=32
+CONFIG_SND_PROC_FS=y
+CONFIG_SND_VERBOSE_PROCFS=y
+CONFIG_SND_SEQUENCER=y
+CONFIG_SND_SEQ_DUMMY=y
+CONFIG_SND_SEQ_HRTIMER_DEFAULT=y
+CONFIG_SND_SEQ_MIDI_EVENT=y
+CONFIG_SND_SEQ_MIDI=y
+CONFIG_SND_DRIVERS=y
+CONFIG_SND_HDA_PREALLOC_SIZE=64
+CONFIG_SND_USB=y
+CONFIG_SND_USB_AUDIO=y
+CONFIG_SND_USB_AUDIO_USE_MEDIA_CONTROLLER=y
+CONFIG_SND_SOC=y
+CONFIG_SND_SOC_GENERIC_DMAENGINE_PCM=y
+CONFIG_SND_SOC_ROCKCHIP=y
+CONFIG_SND_SOC_ROCKCHIP_I2S=y
+CONFIG_SND_SOC_ROCKCHIP_I2S_TDM=y
+CONFIG_SND_SOC_ROCKCHIP_PDM=y
+CONFIG_SND_SOC_ROCKCHIP_SPDIF=y
+CONFIG_SND_SOC_ROCKCHIP_SPDIFRX=y
+CONFIG_SND_SOC_ROCKCHIP_MAX98090=y
+CONFIG_SND_SOC_ROCKCHIP_MULTICODECS=y
+CONFIG_SND_SOC_ROCKCHIP_RT5645=y
+CONFIG_SND_SOC_ROCKCHIP_HDMI=y
+CONFIG_SND_SOC_DUMMY_CODEC=y
+CONFIG_SND_SOC_HDMI_CODEC=y
+CONFIG_SND_SOC_ES7202=y
+CONFIG_SND_SOC_ES7243E=y
+CONFIG_SND_SOC_ES8311=y
+CONFIG_SND_SOC_ES8316=y
+CONFIG_SND_SOC_MAX98090=y
+CONFIG_SND_SOC_RK3308=y
+CONFIG_SND_SOC_RK3328=y
+CONFIG_SND_SOC_RK817=y
+CONFIG_SND_SOC_RK_CODEC_DIGITAL=y
+CONFIG_SND_SOC_RL6231=y
+CONFIG_SND_SOC_RT5616=y
+CONFIG_SND_SOC_RT5640=y
+CONFIG_SND_SOC_RT5645=y
+CONFIG_SND_SOC_RT5651=y
+CONFIG_SND_SOC_SPDIF=y
+CONFIG_SND_SOC_TS3A227E=y
+CONFIG_SND_SIMPLE_CARD_UTILS=y
+CONFIG_SND_SIMPLE_CARD=y
+CONFIG_ANDROID_PARANOID_NETWORK=y
+CONFIG_ACCESS_TOKENID=y
+CONFIG_F2FS_GRADING_SSR=y
+CONFIG_OVERLAY_FS=y
+CONFIG_HUGETLBFS=y
+CONFIG_HUGETLB_PAGE=y
+CONFIG_CRYPTO_SEQIV=y
```
- 追加docker在openharmony上正常运行所需的一些目录
修改openharmony3.1/build/ohos/images/build_image.py
```
def _prepare_root(system_path, target_cpu):
    root_dir = os.path.join(os.path.dirname(system_path), 'root')
    if os.path.exists(root_dir):
        shutil.rmtree(root_dir)
    os.makedirs(root_dir, exist_ok=True)

    _dir_list = [
        'config', 'dev', 'proc', 'sys', 'updater', 'system', 'vendor', 'data',
        'chipset', 'storage', 'mnt', 'tmp', 'sys_prod', 'chip_prod',
        'run', 'var', 'opt', 'usr'
    ]

    for _dir_name in _dir_list:
        os.makedirs(os.path.join(root_dir, _dir_name), exist_ok=True)
    os.symlink('/system/bin', os.path.join(root_dir, 'bin'))
    os.symlink('/system/bin/init', os.path.join(root_dir, 'init'))
    os.symlink('/system/etc', os.path.join(root_dir, 'etc'))

    if target_cpu == 'arm64':
        os.symlink('/system/lib64', os.path.join(root_dir, 'lib'))
    else:
        os.symlink('/system/lib', os.path.join(root_dir, 'lib'))
```


修改openharmony3.1/base/security/selinux/sepolicy/base/system/file_contexts
```
# 增加for docker
/run                u:object_r:rootfs:s0
/var                u:object_r:rootfs:s0
/opt                u:object_r:rootfs:s0
/usr                u:object_r:rootfs:s0
/lib                u:object_r:rootfs:s0
```



##### 3. 编译并且在设备运行OpenHarmony系统

- 编译

```
# 拉取最新openharmony编译环境docker镜像，镜像较大，拉取时间较长，请耐心等待
docker pull swr.cn-south-1.myhuaweicloud.com/openharmony-docker/openharmony-docker:1.0.0

# 进入openharmony源码目录下启动镜像
cd openharmony
docker run --name ohos_build -it -v $(pwd):/home/openharmony swr.cn-south-1.myhuaweicloud.com/openharmony-docker/openharmony-docker:1.0.0

# 预编译工具包：下载和编译时间较长，请耐心等待
./build/prebuilts_download.sh

# 这边使用的润和DAYU200开发板是rk3568主板，大家根据自己使用的主板进行name的填写
# 执行编译脚本
./build.sh --product-name rk3568 --ccache

# 编译后的img放在路径 out/rk3568/packages/phone/images 目录里
```

- 烧录
- 并且通过hdc_std标准工具（通过蓝色的debug线也就是烧录线与PC进行通信）进入oh系统

```
# 查看是否有oh设备
D:>hdc_std.exe list targets 
 
# 进入oh系统
D:>hdc_std.exe shell

# 修改root目录下的权限使其可以进行文件操作
mount -o rw,remount   -t  auto /
```



##### 4. 格式化sd卡为f2fs文件系统

- docker overlay filesystem推荐backing filesystem是未加密的f2fs。而RK3568的data分区是加密的ext4，可以通过micro sd card卡扩展RK3568的存储将sd card格式化为f2fs解决此问题。
- 准备一个sd卡，插入到DAYU200板子上
```
# 查看系统文件系统格式
blkid

# 查看系统文件系统和挂载情况
df -h

# 确定sd卡的名称并且格式化sd卡为f2fs
blkid
mkfs.f2fs     /dev/sdc1 
```



##### 5. 安装ip、iptables和busybox必备组件

- 通过有线或者无线连接网络
```
# 查看正在运行的网络信息
ifconfig

# 查看所有的网络接口 
ifconfig  -a

# 开启ip forward
echo "1" > /proc/sys/net/ipv4/ip_forward
```

- 由于润和DAYU200是arm64类型的主板，所以需要安装第三方组件的arm64版本并且是静态二进制文件。
- 有三种方式：在官网直接获得静态二进制文件，在arm主机源码编译获得静态二进制文件以及在x86_64主机交叉编译获得静态二进制组件。

```
# 交叉编译环境的搭建
# 查看有哪些编译工具的版本可以安装
apt-cache search aarch64

# 选择某个版本的进行安装
apt-get install gcc-10-aarch64-linux-gnu
# 如果出现缺失E: Unmet dependencies 错误，再运行下面一条命令即可
apt --fix-broken install

# 可以查看里边的版本信息，里边有host是什么，target是什么
aarch64-linux-gnu-gcc -v
```

- 安装iproute2
```
# 交叉编译的方式获得静态二进制文件

# 下载iproute2-4.9.0源码
wget https://mirrors.edge.kernel.org/pub/linux/utils/net/iproute2/iproute2-4.9.0.tar.gz 

# 生成 Makefile
./configure

# 修改 Makefile 
CC:=/toolchain/gcc-linux-gnueabi/bin/arm-linux-gnueabi-gcc
AR:=/toolchain/gcc-linux-gnueabi/bin/arm-linux-gnueabi-ar
SUBDIRS=lib ip

# 静态编译
make LDFLAGS=-static

# 过程可参考 https://www.cxybb.com/article/hylaking/95336108

# 如果有arm主机，也可以直接在arm主机上直接静态编译源码（可选）
./configure
make LDFLAGS=-static

# 在源码文件目录下可以找到编译好的可执行二进制文件ip
# 将ip文件上传传到openhramony系统内
hdc_std.exe file send d:\ip  /data/tmp/

# 如果发送失败，修改一下主目录的权限
hdc_std.exe shell
mount -o rw,remount   -t  auto /

# 在openharmony内给ip权限并且将其加入环境变量,便可使用iproute2命令了
chmod +x ip
export PATH=$PATH:/data/tmp/
```

- 安装iptables
```
# iptables与iproute2的安装类似

# 源码下载并且解压
wget http://www.netfilter.org/projects/iptables/files/iptables-1.8.7.tar.bz2

# 交叉编译或者在arm主机直接静态编译获得静态二进制文件
# 交叉编译过程可以参考 https://www.cnblogs.com/eleclsc/p/11686287.html

# 通过hdc_std上传到openharmony系统内
hdc_std.exe file send d:\iptables-1.8.7  /data/tmp/

# 在openharmony内给iptables权限并且将其加入环境变量,便可使用iptables命令了

# 比较特别的是虽然DAYU200板子是arm64位的，iptables交叉编译成32位openharmony才能正常使用
```

- 安装busybox组件
```
# busybox官网直接提供了arm64静态二进制文件
wget https://busybox.net/downloads/binaries/1.28.1-defconfig-multiarch/busybox-armv8l

# 通过hdc_std上传到openharmony系统内,如下操作即可使用busybox的所有工具了
mkdir /system/xbin/
hdc_std.exe file send d:\busybox  /system/xbin/
cd /system/xbin
chmod +x busybox
export PATH=$PATH:/system/xbin/
busybox --install .

# 具体可以参考
https://www.cnblogs.com/biang/p/6703238.html
```

- 放入os-release组件
```
# 下载os-release 官网链接为 https://www.linux.org/docs/man5/os-release.html

# 存放的路径
/usr/lib/os-release
```



##### 6. 安装docker容器引擎组件

- 创建cgroups.json
```
# 在/etc/下创建cgroups.json，用busybox包含的工具创建即可
touch /etc/cgroups.json
cd etc
vi cgroups.json

# cgroups.json内容填写一下内容

{
  "Cgroups": [
    {
      "UID": "system",
      "GID": "system",
      "Mode": "0755",
      "Controller": "blkio",
      "Path": "/dev/blkio"
    },
    {
      "UID": "system",
      "GID": "system",
      "Mode": "0755",
      "Controller": "cpu",
      "Path": "/dev/cpu"
    },
    {
      "Mode": "0555",
      "Path": "/dev/cpuacct",
      "Controller": "cpuacct"
    },
    {
      "UID": "system",
      "GID": "system",
      "Mode": "0755",
      "Controller": "cpuset",
      "Path": "/dev/cpuset"
    },
    {
      "UID": "system",
      "GID": "system",
      "Mode": "0755",
      "Controller": "memory",
      "Path": "/dev/memcg"
    },
    {
      "UID": "system",
      "GID": "system",
      "Mode": "0755",
      "Controller": "schedtune",
      "Path": "/dev/stune"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "devices",
      "Path": "/dev/devices"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "freezer",
      "Path": "/dev/freezer"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "hugetlb",
      "Path": "/dev/hugetlb"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "net_cls",
      "Path": "/dev/net_cls"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "net_prio",
      "Path": "/dev/net_prio"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "perf_event",
      "Path": "/dev/perf_event"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "pids",
      "Path": "/dev/pids"
    },
    {
      "GID": "system",
      "UID": "system",
      "Mode": "0755",
      "Controller": "rdma",
      "Path": "/dev/rdma"
    }
  ],
  "Cgroups2": {
    "UID": "root",
    "GID": "root",
    "Mode": "0600",
    "Path": "/dev/cg2_bpf"
  }
}

```

- 安装docker静态二进制文件
```
# 下载docker static binaries,若为32位选择armhf版。
wget https://download.docker.com/linux/static/stable/aarch64/docker-18.09.2.tgz 

# 将它传到OpenHarmony上并且解压，加入环境变量
tar zxvf 到/system/bin下
export PATH=$PATH:/system/bin/
export PATH=$PATH:/system/bin/docker/
```



##### 7. docker环境准备

执行下述脚本run.sh（如果脚本无法运行也可以按脚本一步一步来）

```
# mkdir on /system
if [ ! -d "/system/etc/docker" ]; then
        mkdir /system/etc/docker
fi

# mkdir on /data
if [ ! -d "/data/var" ]; then
        mkdir /data/var
else
        rm -rf /data/var/run
fi
if [ ! -d "/data/run" ]; then
        mkdir /data/run
fi
if [ ! -d "/data/tmp" ]; then
        mkdir /data/tmp
fi
if [ ! -d "/data/opt" ]; then
        mkdir /data/opt
fi
if [ ! -d "/data/etc" ]; then
        mkdir /data/etc
        mkdir /data/etc/docker
fi
if [ ! -d "/data/usr" ]; then
        mkdir /data/usr
fi
if [ ! -d "/mnt/f2fs" ]; then
        mkdir /mnt/f2fs
fi

# 这里对应刚刚已经被f2fs格式化的sd卡设备，可以用blkid查看到具体的名称
mount /dev/block/mmcblk0p1 /mnt/f2fs/

mount tmpfs /sys/fs/cgroup -t tmpfs -o size=1G
if [ ! -d "/sys/fs/cgroup/blkio" ]; then
        mkdir /sys/fs/cgroup/blkio
        mkdir /sys/fs/cgroup/cpu
        mkdir /sys/fs/cgroup/cpuacct
        mkdir /sys/fs/cgroup/cpuset
        mkdir /sys/fs/cgroup/devices
        mkdir /sys/fs/cgroup/freezer
        mkdir /sys/fs/cgroup/hugetlb
        mkdir /sys/fs/cgroup/memory
        mkdir /sys/fs/cgroup/net_cls
        mkdir /sys/fs/cgroup/net_prio
        mkdir /sys/fs/cgroup/perf_event
        mkdir /sys/fs/cgroup/pids
        mkdir /sys/fs/cgroup/rdma
        mkdir /sys/fs/cgroup/schedtune
        mkdir /sys/fs/cgroup/systemd
fi

# mount --bind
mount --bind /data/etc/docker /etc/docker
mount --bind /data/var /var
mount --bind /data/run /run
mount --bind /data/tmp /tmp
mount --bind /data/opt /opt
mount --bind /data/usr /usr

mount -t cgroup -o none,name=systemd cgroup /sys/fs/cgroup/systemd
mount -t cgroup -o blkio,nodev,noexec,nosuid cgroup /sys/fs/cgroup/blkio
mount -t cgroup -o cpu,nodev,noexec,nosuid cgroup /sys/fs/cgroup/cpu
mount -t cgroup -o cpuacct,nodev,noexec,nosuid cgroup /sys/fs/cgroup/cpuacct
mount -t cgroup -o cpuset,nodev,noexec,nosuid cgroup /sys/fs/cgroup/cpuset
mount -t cgroup -o devices,nodev,noexec,nosuid cgroup /sys/fs/cgroup/devices
mount -t cgroup -o freezer,nodev,noexec,nosuid cgroup /sys/fs/cgroup/freezer
mount -t cgroup -o hugetlb,nodev,noexec,nosuid cgroup /sys/fs/cgroup/hugetlb
mount -t cgroup -o memory,nodev,noexec,nosuid cgroup /sys/fs/cgroup/memory
mount -t cgroup -o net_cls,nodev,noexec,nosuid cgroup /sys/fs/cgroup/net_cls
mount -t cgroup -o net_prio,nodev,noexec,nosuid cgroup /sys/fs/cgroup/net_prio
mount -t cgroup -o perf_event,nodev,noexec,nosuid cgroup /sys/fs/cgroup/perf_event
mount -t cgroup -o pids,nodev,noexec,nosuid cgroup /sys/fs/cgroup/pids
mount -t cgroup -o rdma,nodev,noexec,nosuid cgroup /sys/fs/cgroup/rdma
mount -t cgroup -o schedtune,nodev,noexec,nosuid cgroup /sys/fs/cgroup/schedtune

# ip route
ip rule add pref 1 from all lookup main
ip rule add pref 2 from all lookup default

###
# 设置DNS名称服务器和docker image注册表
echo "{\"registry-mirrors\":[\"https://docker.mirrors.ustc.edu.cn\"],\"experimental\":false,\"storage-driver\": \"overlay2\",\"data-root\": \"/mnt/f2fs\"}" > /etc/docker/daemon.json

# 打开br_netfilter模块
# modprobe br_netfilter
setenforce 0
```



##### 8. 运行docker

```
# 运行 dockerd
dockerd -D -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock &
```



##### 9. 验证docker工作状态

```
# docker run hello-world

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (arm64v8)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/


# 查看docker的信息
docker version
docker info
docker ps -a
docker images
```

![](image/dockerhello.jpg)



#### 云端安装KubeEdge cloudcore并且获得token

```
# 以centos7为例  部署k8s v1.21.0


# 1. 关闭防火墙、selinux和swap分区
systemctl stop firewalld
systemctl disable firewalld
sed -i 's/enforcing/disabled/' /etc/selinux/config
setenforce 0
swapoff -a
 
# 2. 安装 k8s

# 3. 最终查看master节点来确定k8s master部分部署完成
kubectl get nodes

# 4. 部署KubeEdge的cloudcore(KubeEdge选择的是1.9.1并且使用keadm安装)

# 下载keadm v1.9.1
wget https://github.com/kubeedge/kubeedge/releases/download/v1.9.1/keadm-v1.9.1-linux-amd64.tar.gz

# 运行cloudcore
cd keadm-v1.9.1-linux-amd64
cd keadm/
./keadm init --advertise-address=106.14.255.17  --kubeedge-version=1.9.1

# 5. 获取云端token
keadm gettoken 
```



#### KubeEdgeOnOpenHarmony

- Arm server上编译edgecore
```
# 1. 匹配OpenHamrony的memory.stat格式

# 这边其实是第三方库未能包含所有可能得memory.stat格式，这边我已提交issues：https://github.com/opencontainers/runc/issues/3643。

# 也可以修改/vendor/github.com/opencontainers/runc/libcontainer/cgroups/fscommon/utils.go

func ParseKeyValue(t string) (string, uint64, error) {

	tmp := strings.Replace(t, ":", "", -1)
	tmpo := strings.Replace(tmp, "\t", "", -1)
	tmpt := strings.Replace(tmpo, " kB", "", -1)
	count := strings.Count(tmpt, " ")
	tmps := strings.Replace(tmpt, " ", "", count-1)

	parts := strings.SplitN(tmps, " ", 3)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("line %q is not in key value format", t)
	}

	value, err := ParseUint(parts[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("unable to convert to uint64: %v", err)
	}

	return parts[0], value, nil
}


# 2. 由于openharmony安装docker用的是overlay2，所以需要修改edge存放kubelet文件的路径。

# 1.12版本之后的edged可以在egdecore.yaml加入一个字段 “rootDirectory”来修改路径。 
# 1.12之前的版本则需要修改/edge/pkg/edged/edged.go
原：
	DefaultRootDir = "/var/lib/edged"
	// ContainerLogsDir is the location of container logs.
	ContainerLogsDir                 = "/var/log/containers"

修改：(与docker overlay2的路径一样)
	DefaultRootDir = "/mnt/f2fs/lib/edged"
	// ContainerLogsDir is the location of container logs.
	ContainerLogsDir                 = "/mnt/f2fs/log/containers"


# 3. 编译edgecore

cd kubeedge-1.9.1
docker build -t kubeedge/edgecore:v1.9.1 -f build/edge/Dockerfile .
docker cp $(docker create --rm kubeedge/edgecore:v1.9.1):/usr/local/bin/edgecore ./edgecore.1.9.1

# 4. 在kubeedge-1.9.1目录下有edgecore.1.9.1可执行文件，拷贝到openharmony板子/system/bin上
```

- 修改openharmony运行时配置和启动edgecore
```
# 1. cpuset.mems文件添加“0”初始量
echo "0" > /dev/cpuset/background/cpuset.mems

# 2. 添加localhost的路径，
echo "0" > 127.0.0.1 localhost localhost  /etc/hosts

# 3. 添加edgecore.yaml
mkdir -p /etc/kubeedge/config
cd /etc/kubeedge/config
edgecore --minconfig > edgecore.yaml

# 4. 修改edgecore.yaml，例如cloud-token
# 5. 启动edgecore
edgecore
```

edgecore.yaml参考
```
# With --minconfig , you can easily used this configurations as reference.
# It's useful to users who are new to KubeEdge, and you can modify/create your own configs accordingly. 
# This configuration is suitable for beginners.
apiVersion: edgecore.config.kubeedge.io/v1alpha1
database:
  dataSource: /var/lib/kubeedge/edgecore.db
kind: EdgeCore
modules:
  edgeHub:
    enable: true
    heartbeat: 15
    httpServer: https://8.142.141.36:10002
    tlsCaFile: /etc/kubeedge/ca/rootCA.crt
    tlsCertFile: /etc/kubeedge/certs/server.crt
    tlsPrivateKeyFile: /etc/kubeedge/certs/server.key
    token: e8ceef8db0cd06fad0cc13c966527eb06b4c09a5f85818635b2a777ea915945d.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTY5NDQ5MzV9.ff38T0LnXpJEdndfSvOuUUofG-5umKQY0IZ36_MFAOs
    websocket:
      enable: true
      handshakeTimeout: 30
      readDeadline: 15
      server: 8.142.141.36:10000
      writeDeadline: 15
  edged:
    cgroupDriver: cgroupfs
    cgroupRoot: ""
    cgroupsPerQOS: true
    clusterDNS: 169.254.96.16
    clusterDomain: cluster.local
    devicePluginEnabled: false
    dockerAddress: unix:///var/run/docker.sock
    enable: true
    gpuPluginEnabled: false
    hostnameOverride: gandroid
    customInterfaceName: wlan0
    podSandboxImage: kubeedge/pause:3.1
    remoteImageEndpoint: unix:///var/run/dockershim.sock
    remoteRuntimeEndpoint: unix:///var/run/dockershim.sock
    runtimeType: docker
  eventBus:
    enable: false
    mqttMode: 2
    mqttQOS: 0
    mqttRetain: false
    mqttServerExternal: tcp://127.0.0.1:1883
    mqttServerInternal: tcp://127.0.0.1:1884
  metaManager:
    metaServer:
      enable: true    
```

最终效果是：

![](image/edgecore.jpg)

![](image/cloudcore.jpg)





# 项目中可能遇到的问题及解决方案



- 问题一：二进制部署Docker流程

二进制部署Docker虽然官网有相应的静态二进制，但是要通过静态二进制运行起来Docker并不是一件容易的事。社区的老师真的提供了很多帮助，其中一位老师提供了他在Android运行Docker容器关键步骤教程让我来做参考。我分享给大家 https://github.com/ThunderSoft001/kubeedgeOnAndroid

- 问题二：OpenHarmony许多开发板不适配无法运行

OpenHarmony还处于一个研发和推广阶段，所以支持OpenHarmony开发板并不是很多。在购入两块开发板不太适配OpenHarmony后，我联系了华为运营通过参与润和DAYU200体验适用活动申请到OpenHarmony官方用的开发板。DAYU200适配OpenHarmony各个标准版本，所以节省了研究开发板的烧录和运行时间。

- 问题三：OpenHarmony内核部分配置即使设置也无效

OpenHarmony 3.1 release使用的是5.10内核，有些配置是默认关闭的需要在deconfig去配置以下才能打开。然后有些配置即使在deconfig设置了，也无法打开。经过搜索和寻求帮助，发现有些配置是有依赖配置的，需要将依赖配置也打开才可以。查看具体配置的熟悉可以通过 https://cateee.net/lkddb/web-lkddb/CGROUPS.html 查看。

- 问题四：iptables静态交叉编译后无法在OpenHarmony上运行

运行docker需要一些基本软件工具，而OpenHarmony很多组件是需要自己去编译安装的。我所使用的DAYU200开发板是rk3568主板即arm64得。经过静态交叉编译iptables源码得到了arm64位的二进制文件，拷贝到OpenHarmony上一直无法运行。寻求社区老师帮助后，交叉编译iptables源码得到arm32位的静态二进制文件就可以在OpenHarmony上运行。解决方法是从64位的降到32位，这种情况平时安装软件的时候也有出现过。

- 问题五：OpenHarmony与KubeEdge的memory.stat格式不匹配

KubeEdge的edgecore在初始化的过程中会去访问系统memory.stat文件获取系统内存相关的信息。然而OpenHarmony在填充memory.stat时做了格式转换，导致edgecore无法读取memory.stat的最后四行。经过代码阅读和考虑，鉴于无法估量OpenHarmony对memory.stat数据的依赖程度，暂时修改了KubeEdge所使用的opencontainers第三方库的代码来解决这个问题。

- 问题六：OpenHarmony作为边缘端成功加入到云端后，云端无法下发任务。

在成功运行起来edgecore后，需要验证KubeEdge是否部署成功。需要从云端下发hello-world到边缘端并且运行起来，然而pod一直处于cannotrun状态。在请教过社区的老师后，发现是docker使用的overlay2模式下运行pod的路径与KubeEdge启动pod路径不同，需要将edgecore的defaultdir改成docker的pod路径即可。
