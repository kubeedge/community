# Roadmap

该文档描述了KubeEdge开发的路线图

[在GitHub中定义的里程碑](https://github.com/kubeedge/kubeedge/milestones)代表了最新的计划。

下面的路线图概述了KubeEdge 2024年度功能特性规划。

## SIG Node

- 持续跟进Kubernetes版本
- 支持在mac、RTOS系统上运行边缘节点
- 支持in-cluster config，以便边缘Pod访问 Kube-APIServer
- 对节点上的设备插件进行改进，如支持多个虚拟GPU
- 节点侧event上报到云
- 支持边缘Serverless
- 远程运维
- 批量纳管边缘节点

## SIG Device-IOT

- 多语言Mapper支持
- 设备发现
- 时序数据库等数据库的集成
- 流数据
- 增强设备管理能力，如设备写入、设备状态监控

## SIG Security

- SLSA / CodeQL（要达到SLSA L4仍有一些工作要做）
- Spiffe调研
- 提供通用接口，支持多种加密算法证书
- 云边消息通道增加认证鉴权

## SIG Scalability

- 集成EdgeMesh规模和性能测试
- 针对IoT设备场景的规模和性能测试

## Stability

- CloudCore的稳定性维护，包括稳定性测试和问题修复
- EdgeMesh稳定性
- 提高云边协同可靠性，如改进边缘 Kube-API接口和logs/exec等稳定性

## SIG Networking

- 服务网格
  - 与istio或kmesh等项目结合，为边缘场景带来更丰富的Service Mesh功能。
- 大规模优化
  - 在大规模部署场景中，边缘kube apiserver的负载很高，考虑使用IPVS（IP虚拟服务器）技术来有效处理请求
  - 具有大量服务的情况下，会大大增加节点上的iptables规则数
- 性能优化：基于eBPF（扩展伯克利包过滤器）的内核级流量转发
- 分布式消息系统

## SIG AI

- LLM模型的分布式部署
  - 使用 KubeEdge 在多个边缘节点上部署大型语言模型 (LLM)。 LLM可用于各种自然语言处理任务，例如代码实现、文本生成、机器翻译、摘要等
  - 分布式部署可以减少LLM的计算消耗，并提高其可扩展性和容错能力
- LLM模型的云边基准测试
  - 比较在边缘节点上运行的 LLM 模型与使用 KubeEdge 的云服务器上运行的 LLM 模型的性能和资源消耗。 LLM 可以根据不同的指标进行评估，例如准确性、速度、内存、CPU 等
  - 基准测试可以帮助针对不同场景和环境优化LLM模型，并确定云边协同的权衡和挑战
- 不同类型LLM模型的整合
  - 将不同类型的LLM模型（例如大语言/视觉/多模态模型）与KubeEdge集成。 LLM模型可以组合以实现更复杂和多样化的语言生成和理解任务，例如问答、对话、图像字幕等
  - 集成可以充分利用各个LLM模型的优势，增强边云系统的整体功能和能力

## SIG Robotics

- 通用机器人控制系统
  - 机器人控制系统的标准协议已开源（https://github.com/kubeedge/robolink），未来将基于该标准实现通用的机器人控制系统

## SIG Testing

- 单元测试覆盖率提升
- 基于场景的e2e测试用例覆盖率提升
- 集成测试
- 一致性测试改进
- 对KubeEdge所需的硬件要求进行测试，例如内存使用情况、带宽及其他指标

## SIG Cluster-Lifecycle

- 消息路由高可用性（HA）支持
- Keink、Keadm工具的增强
- Edgecore配置可用于Keadm join
- 安装工具Keadm增强
- 优化安装（keadm join）流程
- 镜像预加载能力增强
- 支持OTA模式

## UI

- Dashboard版本迭代

## Experience

- Example库增强
- 上线到Killer-Coda