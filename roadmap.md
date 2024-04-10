# Roadmap

This document defines a high level roadmap for KubeEdge development.

The [milestones defined in GitHub](https://github.com/kubeedge/kubeedge/milestones) represent the most up-to-date plans.

The roadmap below outlines KubeEdge’s 2024 feature plan.

## SIG Node

- Continuous follow up Kubernetes release
- Support edge nodes running on mac OS and RTOS 
- Support in-cluster config for edge pods to access Kube-APIServer
- Enhancements to device plugin on edge nodes, such as support for multiple virtual GPUs
- Support event report to cloud 
- Support for serverless computing
- Remote maintenance
- Add edge nodes in batches

## SIG Device-IOT

- Multi-language Mappers support
- Devices discovery 
- Integration with time-series databases and other databases
- Video stream 
- Enhance device management capabilities, such as device writing and device status monitoring

## SIG Security

- SLSA / CodeQL (There is still some provenance work remaining to reach SLSA L4)
- Spiffe research
- Support for certificates with multiple encryption algorithms, and provide interface capabilities
- Add admission for edge-cloud Messaging Channel 

## SIG Scalability

- Scalability and performance testing with EdgeMesh integrated
- Scalability and performance testing for IoT devices scenario

## Stability

- Stability maintenance of CloudCore, including stability testing and issue resolution
- EdgeMesh stability
- Enhanced reliability of cloud-edge collaboration, such as stability improvement of - Edge Kube-API interface and logs/exec feature

## SIG Networking

- ServichMesh
  - Combined with projects such as istio or kmesh to bring richer service mesh functions to edge scenarios.
- Large-scale optimization
  - In large-scale deployments, there is a high load on the edge kube-apiserver. Consider using IPVS (IP Virtual Server) technology to handle the requests efficiently
  - Having a large number of services significantly increases the number of iptables rules on the nodes
- Performance optimization: Kernel-level traffic forwarding based on eBPF (extended Berkeley Packet Filter)
- Distributed messaging system

## SIG AI

- Distributed deployment of the LLM model 
  - Deploy a large language model (LLM) on multiple edge nodes using KubeEdge. The LLM can be used for various natural language processing tasks, such as code implementations, text generation, machine translation, summarization, etc
  - The distributed deployment can reduce the computation consumption of the LLM, as well as improve its scalability and fault tolerance
- Edge-Cloud benchmarking of the LLM model 
  - Compare the performance and resource consumption of the LLM model running on the edge nodes versus the cloud servers using KubeEdge. The LLM can be evaluated on different metrics, such as accuracy, speed, memory, CPU, etc 
  - The benchmarking can help optimize the LLM model for different scenarios and environments, as well as identify the trade-offs and challenges of edge-cloud collaboration
- Integration of different types of LLM models
  - Integrate different types of LLM models, such as large language/ visual/ multi-modal models, with KubeEdge. The LLM models can be combined to achieve more complex and diverse language generation and understanding tasks, such as question answering, dialogue, image captioning, etc
  - The integration can leverage the advantages of each LLM model and enhance the overall functionality and capability of the edge-cloud system

## SIG Robotics

- Universal robot control system
  - The standard protocol for robot control systems has been open sourced(https://github.com/kubeedge/robolink), and a universal robot control system will be implemented based on this standard in the future

## SIG Testing

- Increase unit test coverage Improve.
- Improve e2e test case coverage 
- Integration testing.
- Conformance test improve
- Perform testing on the hardware requisites required for KubeEdge, such as memory usage, bandwidth, and other metrics

## SIG Cluster-Lifecycle

- Router High Availability (HA) support
- Enhancement for Keink tool, Keadm tool
- Edgecore config can be used for Keadm join
- Enhance the installation tool Keadm
- Optimize installation(keadm join) process
- Enhancement for image prepull 
- Support OTA mode 

## UI

- Dashboard release iteration

## Experience

- Example library enhancement
- Go online to Killer-Coda