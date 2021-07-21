# SIG Robotics Charter

This charter adheres to the conventions described in [KubeEdge Open Governance](https://github.com/kubeedge/community/blob/master/GOVERNANCE.md) and uses the Roles and Organization Management outlined in the governance doc.

## Scope

The SIG Robotics focuses on technical discussion, API definition, reference architecture, implementation in robots, to empower existing robotics ecosystems with KubeEdge, based on the edge-cloud synergy architecture, to improve robot intelligence and development efficiency.

### In scope
#### Areas of Focus
1. Based on KubeEdge, to implement cloud robot platform, to integrate cloud computing technologies(cloud native, AI, storage etc) with robots:
    - Support heterogeneous robots, including mobile robots, robotics arm, and many other robot platforms.
    - Integrate popular open source robotics technologies with KubeEdge to make the robots development, debug, simulation, deployment and management easier.
    - Cooperate with other open source communities, e.g.,  [ROS](https://www.ros.org/), [Gazebo](http://gazebosim.org/), also with other SIGs like [SIG AI](https://github.com/kubeedge/community/tree/master/sig-ai) and [SIG IoT](https://github.com/kubeedge/community/tree/master/sig-device-iot)
1. Synergy mechanisms research for robotics workloads, **including but not limited to**:
    - Cloud environment for robot programming, compilation, simulation, packaging and deployment
    - Integration with cloud services like AI training, data storage, monitoring
    - Edge-cloud synergy system for robots
        - Robot application deployment to robot base on edge
        - Fleet ops
        - Tele ops
        - Cloud service integration like AI(training and inference), data storage, monitoring
    - Additional protocol integration like [Micro-XRCE-DDS](https://github.com/eProsima/Micro-XRCE-DDS)
- Provide basic container environment, examples.
- Cooperate with partners to make use cases to promote cloud robot technology.
### Out of scope
- Re-invent existing robotics framework or tooling, e.g., ROS, Gazebo, DDS
- Offer domain/application-specific algorithms, e.g., robot hardware, path planning, SLAM


## Roles and Organization Management

This SIG follows and adheres to the Roles and Organization Management outlined in KubeEdge Open Governance and opts-in to updates and modifications to KubeEdge Open Governance.

### Additional responsibilities of Chairs

- Manage and curate the project boards associated with all sub-projects ahead of every SIG meeting so they may be discussed
- Ensure the agenda is populated 24 hours in advance of the meeting, or the meeting is then cancelled
- Report the SIG status at events and community meetings wherever possible
- Actively promote diversity and inclusion in the SIG
- Uphold the KubeEdge Code of Conduct especially in terms of personal behavior and responsibility
