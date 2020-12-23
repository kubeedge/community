# SIG AI Charter

This charter adheres to the conventions described in [KubeEdge Open Governance](https://github.com/kubeedge/community/blob/master/GOVERNANCE.md) and uses the Roles and Organization Management outlined in the governance doc.

## Scope

SIG AI is responsible to provide general platform capabilities based on KubeEdge so  that AI applications running at the edge can benefit from cost reduction, model performance improvement, and data privacy protection.

### In scope

#### Areas of Focus

- Collaborate with AI framework open source communities, to help better running on KubeEdge, e.g. Tensorflow, Pytorch, PaddlePaddle, onnx etc. Works including but not limited to:
    - support heterogeneous hardware, such as Ascend, Kunglun, Cambrian, and Rockchip.
    - support ONNX models, such as tutorials.
- provide an **edge-cloud collaborative** AI framework based on KubeEdge capabilities, such as model and dataset.
- provide an end-to-end edge-cloud collaborative AI benchmarking framework, for typical AI applications scenarios on KubeEdge.
    - Contextual Metrics
    - End-to-end Testbed
        - Consisting of data collection, data preprocess, train, inference and other parts 
        - Covering cloud, edge, and end-device layers 
- research interest **include but not limited as follow**:
    - cloud training and edge inference， integrationg with kubeflow
    - incremental learning
    - joint inference
    - federated learning 
    - Knowledge distillation training model, which can be combined with join inference as a means of training large and small models.
    - KubeEdgeFlow, extend KubeFlow capabilities to edge


### Out of scope
- to re-invent existing ML framework, i.e., tensorflow, pytorch, mindspore, etc.
- to re-invent KubeEdge platform capability.
- to offer domain/application-specific algorithms, i.e., facial recognition, text classification, etc.


## Roles and Organization Management

This SIG follows and adheres to the Roles and Organization Management outlined in KubeEdge Open Governance and opts-in to updates and modifications to KubeEdge Open Governance.

### Additional responsibilities of Chairs

- Manage and curate the project boards associated with all sub-projects ahead of every SIG meeting so they may be discussed
- Ensure the agenda is populated 24 hours in advance of the meeting, or the meeting is then cancelled
- Report the SIG status at events and community meetings wherever possible
- Actively promote diversity and inclusion in the SIG
- Uphold the KubeEdge Code of Conduct especially in terms of personal behavior and responsibility
