# Roadmap

This document defines a high level roadmap for KubeEdge development.

The [milestones defined in GitHub](https://github.com/kubeedge/kubeedge/milestones) represent the most up-to-date plans.

The roadmap below outlines KubeEdge’s 2026 feature plan.

## SIG Node

- Continuous follow up Kubernetes release 
- Mac OS and RTOS support
- Edge Cluster and edge swarm cluster
- Support for serverless computing
- Support runtimeclass

## SIG Device-IOT

- Device discovery
- Batch devices management
- Establishment and maintenance of the multi-language mapper-framework repo

## SIG Security

- Spiffe research

## SIG Scalability

- Scalability and performance testing with EdgeMesh integrated
- Scalability and performance testing for IoT devices scenario

## SIG Networking

- ServichMesh
  - Combined with projects such as `istio` or `kmesh` to bring richer service mesh functions to edge scenarios.
- Performance optimization: Kernel-level traffic forwarding based on eBPF (extended Berkeley Packet Filter)
- Distributed messaging system

## SIG AI

- AI Conformance support
- Provide hands-on examples for deploying models at edge

## SIG Robotics

- Universal robot control system
- Provide examples and solutions for robotics scenarios

## SIG Testing

- AI generate UT and E2E test
- Integration testing optimization
- Perform testing on the hardware requisites required for KubeEdge, such as memory usage, bandwidth, and other metrics

## SIG Cluster-Lifecycle

- Remote maintenance enhancement
- CloudCore HA enhancement
- Stability improvement for Edge Kube-API interface and logs/exec feature

## UI

- Add new tutorial page on official website
- Container deployment for Dashboard
- Dashboard enhancement iteration

## Documentation

- Multi-language docs maintenance with AI tool
- Restructure and refactor the contents.

## Experience

- Example repository integrates with KubeEdge core in CI/CD pipeline
