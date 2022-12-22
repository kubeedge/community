# SIG Node Charter

This charter adheres to the conventions described in the [KubeEdge Open Governance](https://github.com/kubeedge/community/blob/master/GOVERNANCE.md) and used the Roles and Organization Management outlined in the governance doc.

## Scope 

The SIG Node is responsible for the edge agent running on various edge device. we controls the interactions between pods and host resources at the edge. We manage the lifecycle of nodes and pods scheduled to the node. We focus on enabling a broad set of workload types, including workloads with hardware specific or performance sensitive requirements. We support an extremely lightweight edge agent(EdgeCore) to run on resource constrained edge. Due to the complex scenario at the edge, we manage the local metadata to support edge autonomy. We also aim to continuously improve node performance and reliability.

### In scope

#### Areas of Focus

- Node lifecycle management.
- Support EdgeCore running on various hardware and os.
- Enable a broad set of workload types at the edge.
- Support Kubelet integrated into Edged and manage its feature.
- Extremely lightweight with Edged(lite-Kubelet).
- Node level performance and scalability.
- Node reliability.
- Node-level resource management and monitoring information collection for hosts at the edge, such as power supply, battery power, etc.
- Nodegroup and edgeapplication management.
- Communication across nodes, such as node dynamic trunk. (with [EdgeMesh](https://github.com/kubeedge/edgemesh))
- Hardware integration test. (with [sig-testing])

### Out of scope

- Device management in edge ([sig-device-iot](https://github.com/kubeedge/community/blob/master/sig-device-iot))
- Node level security ([sig-security](https://github.com/kubeedge/community/tree/master/sig-security))

## Roles and Organization Management

This SIG follows and adheres to the Roles and Organization Management outlined in KubeEdge Open Governance and opts-in to updates and modifications to KubeEdge Open Governance.

### Additional responsibilities of Chairs

- Manage and curate the project boards associated with all sub-projects ahead of every SIG meeting so they may be discussed
- Ensure the agenda is populated 24 hours in advance of the meeting, or the meeting is then cancelled
- Report the SIG status at events and community meetings wherever possible
- Actively promote diversity and inclusion in the SIG
- Uphold the KubeEdge Code of Conduct especially in terms of personal behavior and responsibility
