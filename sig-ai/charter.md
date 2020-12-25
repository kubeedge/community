# SIG AI Charter

This charter adheres to the conventions described in [KubeEdge Open Governance](https://github.com/kubeedge/community/blob/master/GOVERNANCE.md) and uses the Roles and Organization Management outlined in the governance doc.

## Scope

The SIG AI focuses on technical discussion, API definition, reference architecture, implementation in the Edge AI field, to enable AI applications better running on edge (including cost saving, performance improvement, and data protection).
    
### In scope
#### Areas of Focus
1. Empower KubeEdge with existing AI ecosystems, to support execution of Edge AI applications and services：
    - Support heterogeneous edge hardware, e.g., Ascend, Kunlun, Cambrian, and Rockchip
    - Integrate typical AI frameworks into KubeEdge, e.g., Tensorflow, Pytorch, PaddlePaddle and Mindspore etc.
    - Integrate KubeFlow and ONNX into KubeEdge, to enable interoperability of edge models with diverse formats
    - Cooperate with other open source communities, e.g., Akraino and LF AI
1. Synergy models research for AI workloads, **including but not limited to**:
    - Cloud training and edge inference
    - Edge-cloud-collaborative inference
        - Knowledge distillation for the cloud and edge model
    - Incremental learning
    - Federated learning
    - Edge model and dataset management
1. Edge AI benchmarking relevant work, to help users identify most important dimensions when developing, evaluating Edge AI applications-and-services system:
    - Provide Contextual Metrics
        - For typical Edge AI applications scenarios
    - Provide Standardized Evaluation Settings
        - Standardized datasets, architectures, and hardware
        - For each routine AI module, e.g., data collection, data preprocessing, train and inference 
        - For each architecture layer, i.e., cloud, edge, and end-device 
### Out of scope
- Re-invent existing AI framework, e.g., Tensorflow, Pytorch and Mindspore
- Offer domain/application-specific algorithms, e.g., facial recognition and text classification


## Roles and Organization Management

This SIG follows and adheres to the Roles and Organization Management outlined in KubeEdge Open Governance and opts-in to updates and modifications to KubeEdge Open Governance.

### Additional responsibilities of Chairs

- Manage and curate the project boards associated with all sub-projects ahead of every SIG meeting so they may be discussed
- Ensure the agenda is populated 24 hours in advance of the meeting, or the meeting is then cancelled
- Report the SIG status at events and community meetings wherever possible
- Actively promote diversity and inclusion in the SIG
- Uphold the KubeEdge Code of Conduct especially in terms of personal behavior and responsibility
