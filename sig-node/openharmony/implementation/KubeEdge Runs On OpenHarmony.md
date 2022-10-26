# KubeEdge runs on RK3568 OpenHarmony



This tutorial is divided into two key steps: the steps of running Dcoker on OpenHarmony and the steps of running KubeEdge on OpenHarmony. It also shows the specific process of running Kubeedge on OpenHarmony devices, using HiHope DAYU200 as examples. Here, the runtime used for KubeEdge is Docker.



## Key steps for OpenHarmony to run Docker

##### 1. Prepare the OpenHarmony kernel that supports Docker 

- Modify kernel configuration:

  cgroup and namespace related features, main modified files is kernel/Linux/config/Linux - 5.10 / arch/arm64 / configs/rk3568_standard_defconfig

- Modify the network kernel configuration:

  The network is mainly in bridge mode, so you need to enable the corresponding configuration and support the network packet forward function.

- Modify the overlay filesystem:

  The overlay filesystem used by Docker recommends the unencrypted f2fs as the backing filesystem.

##### 2. Compile and install OpenHarmony on your device

##### 3. Install the necessary tools

Some of the necessary tools for Docker are iproutes, iptables, and busybox. Copy the static binaries of these tools to the /bin/ directory on the OpenHarmony device and add execution permissions.

##### 4. Install the Docker container engine components

Copy the Docker static binary file to the /bin/ directory on the Android device and add the execution permission.

##### 5. Mount OpenHarmony system resources

Create OpenHarmony's /etc/cgroups.json file to mount all cgroup subsystems.

##### 6. Set the OpenHarmony environment

- Create directories needed for Docker to run

```
# Open the access of access
mount -o rw,remount   -t  auto /

# Create related directories
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

- Directory needed to mount for Docker

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

- Add some route rules

```
ip rule add pref 1 from all lookup main
ip rule add pref 2 from all lookup default
```

- Close the selinux

```
setenforce 0
```

- Create file /etc/dock/daemon. json and write

```
{"registry-mirrors":["https://docker.mirrors.ustc.edu.cn"],"experimental":true}
```

- Run Docker

```
dockerd -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock & 
```

- Verify docker working status

```
# Ensure the OpenHarmony is networked then execute hello-world
docker run hello-world

# You can see
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
93288797bd35: Pull complete 
Digest: sha256:2498fce14358aa50ead0cc6c19990fc6ff866ce72aeb5546e1d59caac3d0d60f
Status: Downloaded newer image for hello-world:latest

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
```



## Key steps for OpenHarmony to run KubeEdge



Use KubeEdge to bridge the cloud and the OpenHarmony edge to achieve cloud-side collaboration.



##### 1. Install KubeEdge cloudcore on the cloud and get token

##### 2. KubeEdge edgecore source static compilation

Compile KubeEdge source code on Arm server to obtain edgecore static binary file.

```
# 1. Match OpenHamrony's memory.stat format,

This is actually a third-party library that fails to include all possible memory.stat formats, which I have submitted to issues and pr to fix.
You can also modify it yourself /vendor/github.com/opencontainers/runc/libcontainer/cgroups/fscommon/utils.go


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


# 2. Since openharmony uses overlay2 to install docker, you need to change the DefaultRootDir of edged.go to the same path as docker overlay2

# 3. Compilation edgecore

cd kubeedge-1.9.1
docker build -t kubeedge/edgecore:v1.9.1 -f build/edge/Dockerfile .
docker cp $(docker create --rm kubeedge/edgecore:v1.9.1):/usr/local/bin/edgecore ./edgecore.1.9.1

# 4. In the kubeedge directory there is edgecore.tag executable file copied to the openharmony board /bin/ 
```



##### 3. Install edgecore on the OpenHarmony edge

```
# The cpuset.mems file needs to add "0" initial amount
echo "0" > /dev/cpuset/background/cpuset.mems

# Add the path to localhost
echo "127.0.0.1 localhost localhost" > /etc/hosts

# Add edgecore.yaml
mkdir -p /etc/kubeedge/config
cd /etc/kubeedge/config
edgecore --minconfig > edgecore.yaml

# Modify cloud-token and mqtt-ip in edgecore.yaml
# Start the edgecore program
```



------





## The reference is implemented on HiHope DAYU200



### Environmental information

- Development board：HiHope DAYU200

- Chip：RockChip RK3568（arm64）

- OpenHarmony version：3.1 release
- Linux kernel：5.10
- Docker：18.03.1
- KubeEdge：1.9.1

(Reminder: 1.9.1 is chosen here because k8s and KubeEdge versions are matched and some higher versions of k8s are not supported.)



### Operation procedure and reference documents

------

#### DockerOnOpenHarmony

##### 1. Export the original configuration for detection

Download the Test script

- [https://github.com/moby/moby/blob/master/contrib/check-config.sh](https://github.com/moby/moby/blob/master/contrib/check-config.sh)
- [https://gitee.com/mirrors_moby/moby/blob/master/contrib/check-config.sh#](https://gitee.com/mirrors_moby/moby/blob/master/contrib/check-config.sh#)
- Download the Test script

Test docker support of the original OpenHarmony3.1release kernel.

- After the source code is compiled, the kernel Img will be generated. The specific location is /out/kernel/src_tmp/linux-5.10/arch/boot/

- Go to the source tree /out/kernel/src_tmp/linux-5.10/

- Enter the scripts/extract-ikconfig arch/arm/boot/*Image > /home/.config command

- Of course, the output.config can go to another path

- ps: Also can the OS run  “cat /proc/config. gz | gzip - > d/sdcard/config”  command under /home generated config to the kernel configuration file similar to the following:

```
# The section of the configuration file：
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

- Run script

```
./check-config.sh    /home/.config
```

- The test results are shown in the figures

​                  ![](image/check-1.png)

​                  ![](image/check-2.png)

​                  ![](image/check-3.png)



##### 2. Modify the source kernel configuration

According to the above check results

- You need to enable both required and optional items

- Network Drivers: we usually use the bridge, do not need to open other.

- Storage Drivers: Open the overlay part.

Modifying Kernel Configuration

- Modify openharmony3.1/kernel/Linux/config/Linux-5.10/arch/arm64/configs/rk3568_standard_defconfig

  (Add namespace, control group, network, overlay filesystem, etc.)

Add the following configuration:

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

- Append some directories that docker needs to function properly on openharmony

  Modify  openharmony3.1/ build/ohos/images/build_image.py

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

​	Modify  openharmony3.1/base/security/selinux/sepolicy/base/system/file_contexts

```
# Add some directories for docker
/run                u:object_r:rootfs:s0
/var                u:object_r:rootfs:s0
/opt                u:object_r:rootfs:s0
/usr                u:object_r:rootfs:s0
/lib                u:object_r:rootfs:s0
```



##### 3. Compile and run the OpenHarmony system on your device

- Compile

```
# Pull the latest OpenHarmony compiled environment docker image. The image is large, so pulling time is long, please be patient.
docker pull swr.cn-south-1.myhuaweicloud.com/openharmony-docker/openharmony-docker:1.0.0

# Go to the OpenHarmony source code directory and start the image
cd openharmony
docker run --name ohos_build -it -v $(pwd):/home/openharmony swr.cn-south-1.myhuaweicloud.com/openharmony-docker/openharmony-docker:1.0.0

# Pre-compiled toolkit: download and compilation time is long, please be patient
./build/prebuilts_download.sh

# I use the HiHope DAYU200 development board is rk3568 motherboard, you fill in the name according to the motherboard you use.
# Execute the compilation script: If the first compilation is unsuccessful and it is not one of the following errors, consider running it again.
./build.sh --product-name rk3568 --ccache

# The compiled img is placed in the directory out/rk3568/packages/phone/images
```

- Burn
- The OH system is accessed via the hdc_std standard tool, which communicates with the PC through the blue debug line, also known as the burn line

```
# Check the availability of OpenHarmony equipment
D:>hdc_std.exe list targets 

# Enter the OpenHarmony system
D:>hdc_std.exe shell

# Modify the permissions in the root directory to allow file operations
mount -o rw,remount   -t  auto /
```



##### 4. Format the sd card as the f2fs file system

- The docker overlay filesystem recommends that the backing filesystem is unencrypted f2fs, while the RK3568 data partition is encrypted ext4, which can be solved by extending the RK3568 storage with a micro sd card and formatting the sd card as f2fs.
- Prepare an sd card and insert it into the DAYU200 board.

```
# View system file system format
blkid

# View system file system and mounts
df -h

# Determine the name of the sd card and format the sd card as f2fs
blkid
mkfs.f2fs     /dev/sdc1 
```



##### 5. Install iproute, iptables and busybox  components

- Network connection via wired or wireless

```
# View the running network
ifconfig

# View all network interfaces 
ifconfig  -a

# Open ip forward
echo "1" > /proc/sys/net/ipv4/ip_forward
```

- Since the HiHope DAYU200 is an arm64 type motherboard, it is necessary to install the arm64 version of the third-party component and it is a static binary file. 
- There are three ways to get static binaries: directly on the official website, source compiled on the arm host to get static binaries and cross-compiled on the x86_64 host to get static binary components.

```
# Cross-compiling environment construction
# Check what versions of the compiler tools are available for installation
apt-cache search aarch64

# Select a version to install
apt-get install gcc-10-aarch64-linux-gnu
# If output error "E: Unmet dependencies",run command
apt --fix-broken install

# You can check the version information inside, what is the host and what is the target
aarch64-linux-gnu-gcc -v
```

- Install iproute2

```
# Cross-compile to get static binaries

# Download iproute2-4.9.0 source code
wget https://mirrors.edge.kernel.org/pub/linux/utils/net/iproute2/iproute2-4.9.0.tar.gz 

# Unzip and go to the source code file and modify the configure file,where CC and ar are changed to aarch64-linux-gnu- path
./configure

# Modify Makefile,CC to be the path of cross-compiler, modify SUBDIRS=lib ip

# Static compilation
make LDFLAGS=-static

# The process can be referred to https://www.cxybb.com/article/hylaking/95336108

# If you have an arm host, you can also statically compile the source code directly on the arm host (optional)
./configure
make LDFLAGS=-static


# The compiled executable binaries ip can be found in the source code file directory.
# Transfer the ip program to the openhramony system using the hdc_std command.
hdc_std.exe file send d:\ip  /data/tmp/

# If sending fails, change the permissions of the home directory.
hdc_std.exe shell
mount -o rw,remount   -t  auto /

# In the OpenHarmony give the ip permission  and add it to the environment variable, then you can use the iproute2 command
chmod +x ip
export PATH=$PATH:/data/tmp/
```

- Install iptables

```
# The installation of iptables and iproute2 is similar

# Download and unzip the source code
wget http://www.netfilter.org/projects/iptables/files/iptables-1.8.7.tar.bz2

# Cross-compile or statically compile directly on the arm host to obtain static binaries
# The cross-compilation process can be found in https://www.cnblogs.com/eleclsc/p/11686287.html

# Uploading files to the openharmony system via the hdc_std tool
hdc_std.exe file send d:\iptables-1.8.7  /data/tmp/

# In the OpenHarmony give the iptables permission  and add it to the environment variable, then you can use the iptables command.

# What is special is that although the DAYU200 board is arm64-bit, iptables is cross-compiled to 32-bit openharmony to work properly.
```

- Install busybox

```
# The busybox website provides the arm64 static binaries directly
wget https://busybox.net/downloads/binaries/1.28.1-defconfig-multiarch/busybox-armv8l

# Upload to the openharmony system via hdc_std and perform the following actions to use all the tools of the busybox.
mkdir /system/xbin/
hdc_std.exe file send d:\busybox  /system/xbin/
cd /system/xbin
chmod +x busybox
export PATH=$PATH:/system/xbin/
busybox --install .

# Specific steps
https://www.cnblogs.com/biang/p/6703238.html
```

- Put in os-release component

```
# Download os-release and the official website link is https://www.linux.org/docs/man5/os-release.html

# Path to storage os-release is /usr/lib/os-release
```



##### 6. Installing docker container engine components

- Create cgroups.json

```
# Create cgroups.json under /etc/
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

- Installing docker static binaries

```
# Download docker static binaries,for 32-bit select armhf version.
wget https://download.docker.com/linux/static/stable/aarch64/docker-18.09.2.tgz 

# Transfer it to OpenHarmony and unzip it, add environment variables
tar zxvf /system/bin/
export PATH=$PATH:/system/bin/
export PATH=$PATH:/system/bin/docker/
```



##### 7. Docker environment preparation

Run the script (if the script does not work you can also follow the script step by step)

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

# This corresponds to the sd card device that has just been formatted with f2fs, and the specific name can be viewed with blkid.
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
# setup dns nameserver and docker images registry
echo "{\"registry-mirrors\":[\"https://docker.mirrors.ustc.edu.cn\"],\"experimental\":false,\"storage-driver\": \"overlay2\",\"data-root\": \"/mnt/f2fs\"}" > /etc/docker/daemon.json

# open br_netfilter module
#modprobe br_netfilter
setenforce 0
```



##### 8. Run Docker

```
# run dockerd
dockerd -D -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock &
```



##### 9. Verify the working status of docker

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


# View information about docker
docker version
docker info
docker ps -a
docker images
```

![](image/dockerhello.jpg)



#### Install KubeEdge cloudcore in the cloud and get token

```
# Take centos7 as an example and deployment k8s v1.21.0


# 1. Turn off the firewall and selinux
systemctl stop firewalld
systemctl disable firewalld
sed -i 's/enforcing/disabled/' /etc/selinux/config

# The final /etc/selinux/config is as follows:
# This file controls the state of SELinux on the system.
# SELINUX= can take one of these three values:
#     enforcing - SELinux security policy is enforced.
#     permissive - SELinux prints warnings instead of enforcing.
#     disabled - No SELinux policy is loaded.
#SELINUX=enforcing
SELINUX=disabled
# SELINUXTYPE= can take one of three two values:
#     targeted - Targeted processes are protected,
#     minimum - Modification of targeted policy. Only selected processes are protected. 
#     mls - Multi Level Security protection.
SELINUXTYPE=targeted

setenforce 0

# 2. Permanently shut down swap partition
swapoff -a

# 3. The official repository is not available, it is recommended to use the Aliyuan repository, execute the following command to add the kubernetes.repo repository

cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes 
baseurl=http://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64 
enabled=1
gpgcheck=0
repo_gpgcheck=0
gpgkey=http://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg
       http://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
 

# 4. Installation kubectl kubeamd kubelet
yum install -y kubelet-1.21.0 kubeadm-1.21.0 kubectl-1.21.0
# k8s How to downgrade version can be referred to  https://blog.csdn.net/u012069313/article/details/125561711


# 5. Then start the kubelet on the master server
kubeadm init   --apiserver-advertise-address=10.0.1.176 --image-repository registry.aliyuncs.com/google_containers   --kubernetes-version v1.21.0   --service-cidr=10.140.0.0/16 --pod-network-cidr=10.240.0.0/16

# Setting the profile 
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown $(id -u):$(id -g) $HOME/.kube/config


# 6. Install the flannel plugin (optional)
# Start flannel on master node 
cd /etc/kubeedge
kubectl create -f kube-flannel.yml

# 7. Finally check the master node to make sure the k8s master part is deployed.
kubectl get nodes

# 8. Deploy KubeEdge's cloudcore (KubeEdge is 1.9.1 and installed with keadm)

# Download keadm v1.9.1
https://github.com/kubeedge/kubeedge/releases

# Unzip
wget https://github.com/kubeedge/kubeedge/releases/download/v1.9.1/keadm-v1.9.1-linux-amd64.tar.gz

# Run
cd keadm-v1.9.1-linux-amd64
cd keadm/
./keadm init --advertise-address=106.14.255.17  --kubeedge-version=1.9.1

# 9.Get token
keadm gettoken 
```



#### KubeEdgeOnOpenHarmony

- Compiling edgecore on Arm server

```
# Match openhamrony's memory.stat format, Modify kubeedge-1.9.1/vendor/github.com/opencontainers/runc/libcontainer/cgroups/fscommon/utils.go


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


# Since OpenHarmony installs docker with overlay2，so need to Modify kubeedge-1.9.1/edge/pkg/edged/edged.go

Origin：
	DefaultRootDir = "/var/lib/edged"
	// ContainerLogsDir is the location of container logs.
	ContainerLogsDir                 = "/var/log/containers"

Modify: (same path as docker overlay2)
	DefaultRootDir = "/mnt/f2fs/lib/edged"
	// ContainerLogsDir is the location of container logs.
	ContainerLogsDir                 = "/mnt/f2fs/log/containers"


# Compile edgecore

cd kubeedge-1.9.1
docker build -t kubeedge/edgecore:v1.9.1 -f build/edge/Dockerfile .
docker cp $(docker create --rm kubeedge/edgecore:v1.9.1):/usr/local/bin/edgecore ./edgecore.1.9.1

# In the kubeedge-1.9.1 directory there is an executable edgecore.1.9.1, copy it to the openharmony board /system/bin
```

- Modify the OpenHarmony runtime configuration and start edgecore

```
# Cpuset.mems file adds "0" initial amount
echo "0" > /dev/cpuset/background/cpuset.mems

# Add the localhost path
echo "0" > 127.0.0.1 localhost localhost  /etc/hosts

# Add edgecore.yaml
mkdir -p /etc/kubeedge/config
cd /etc/kubeedge/config
edgecore --minconfig > edgecore.yaml

# Modify edgecore.yaml, e.g. cloud-token and mqtt-ip

# Start edgecore
edgecore
```

The reference of edgecore.yaml：

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

The end result：

![](image/edgecore.jpg)





# Problems and solutions that may be encountered in the project

- Problem 1: Binary deployment for Docker process.

Binary deployment Docker: Although there is a corresponding static binary on the official website, it is not easy to run up Docker via static binary. The teachers in the community really helped a lot, and one of them provided his tutorial on the key steps to run Docker containers on Android for me to refer to.  I share to you  https://github.com/ThunderSoft001/kubeedgeOnAndroid

- Problem 2: Many development boards are not adapted to run OpenHarmony.

OpenHarmony is still in the research development and promotion stage, so there are not many development boards supporting OpenHarmony. After purchasing two development boards that were not quite compatible with OpenHarmony, I contacted Huawei to apply for an official OpenHarmony development board by participating in the HiHope DAYU200 experience.

- Problem 3: Some configurations of OpenHarmony kernel are invalid even if set.

OpenHarmony 3.1 release is using 5.10 kernel, some configurations are turned off by default and need to be configured in deconfig to turn them on. Then some configurations could not be turned on even if they were set in deconfig. After searching and seeking help, I found that some configurations have dependent configurations and need to turn on the dependent configurations as well. Check the familiarity of the specific configuration can be viewed at https://cateee.net/lkddb/web-lkddb/CGROUPS.html.

- Problem 4: Iptables does not work on OpenHarmony after static cross-compilation.

Running docker requires some basic software tools, and many components of OpenHarmony need to be compiled and installed by yourself. The DAYU200 development board I am using is an rk3568 motherboard, i.e. arm64. After cross-compiling the iptables source code statically, I got the arm64 bit binary file, and copied it to OpenHarmony, but it never worked. After seeking help from community teachers, I was able to cross-compile the iptables source code to get the static binary file of arm32 bit and it will run on OpenHarmony. The solution was to downgrade from 64-bit to 32-bit, a situation that usually occurs when installing software.

- Problem 5: The file  format which is memory.stat  of  OpenHarmony and KubeEdge  don't match

KubeEdge's edgecore will access the system memory.stat file to get the system memory related information during the initialization process. However, OpenHarmony did a format conversion when populating memory.stat, resulting in edgecore not being able to read the last four lines of memory.stat. After reading and considering the code, I temporarily modified the code of the opencontainers third-party library used by KubeEdge to solve this problem, since I could not estimate the dependency of OpenHarmony on the memory.stat data.

- Problem 6: After OpenHarmony is successfully joined to the cloud as an edge, the cloud is unable to issue tasks.

After successfully running edgecore, I need to verify that KubeEdge is deployed successfully. I need to send hello-world from the cloud to the edge and run it, however, the pod keeps being unable to run. After consulting the teachers in the community, I found that the path to run the pod in overlay2 mode used by docker is different from the path to start the pod in KubeEdge, so I need to change the defaultdir of edgecore to the path to the pod in docker.

































