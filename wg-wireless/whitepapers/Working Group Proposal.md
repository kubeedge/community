KubeEdge-Wireless Working Group Proposal

Author: Jason Ning

### Target: 

- Our goal is to make an open platform to enable **Wireless**-Edge computing, extending native containerized application orchestration capabilities to hosts at **Wireless**-Edge, which built upon kubernetes and provides fundamental infrastructure support for **Wireless**-network, app deployment and meta-data synchronization between cloud and **Wireless**-Edge.
- **Extend** KubeEdge to support **wireless mobile scenarios**.

### Overview:

- With the development of computing and communication technologies, powerful computing equipments constitutes multi-agent systems. The agents have both computing and communication capacities and interact with each others. Therefore, the problem of wireless agents interaction between multiple agents and system management is the main scenario discussed here. 
- The communication relationship between different Kubeedge nodes can not be easily seen as connect or disconnect,  it should be more fine-grain communication management.

This proposal is about to carry out the wireless network issues under kubeedge management. 



### In scope

- Cloud or edge service quality problems caused by user movement
- Random off-line issues present new challenges to kubeEdge management in wireless communication. 
- Optimize the quality of communication under the scenario of side collaboration
- Proactive changes in communication intensity require decentralized cloud management capabilities.

### Out of scope

- Do not invent communication protocol.

### Scenario：

**Scenario_1: Vehicle Network/Vessel Network/UAV Network:**

<img src="https://snz04pap002files.storage.live.com/y4m9TFrXL6j2ajyTvHbGHvTWUo5w6kT7SMZ2DZq1t1ivm9Trt4Lo7wt4-azpzI5O746g8fXqBkMz0iweMjyNQdBS4WaXnv5HLtyxU_kHOFNGrEpgfrbocPHwC1Jb6TnskFPTfWwvQUzkhgsz-AoRc1VffTLSUX8klBth8EMHG_KmYO5zoKt_Vpk8XVRDn2UO_-J?width=640&height=397&cropmode=none" width="640" height="397" style="zoom: 67%;" />

**Key differences** between wireless and wire:  

- Communication mode may need to change from TCP, IP mode to multi-cast or broadcast mode.
- Network for dynamics scenarios.  

**Description**:  In the mobile scenario, multiple vehicles support similar services, and vehicles wireless connected with each other as KubeEdge nodes.
- Off-line autonomy and node management: Five cars form an mesh network, which can transfer information processing between each other in off-line conditions.
- Leader selection: Choose a node as a cluster head, to do the overall management of other equipments.
- KubeEdge-wireless should monitor the whole or partly network status, to aid in networking decisions.
- According to the Service Level Agreement, KubeEdge can make decisions on the networking modes and limit the network scope.
- Inspired by 3GPP 36.885 standards.

**Scenario_2: Air/Sea rescue Collaboration：**

<img src="https://snz04pap002files.storage.live.com/y4mJdEgovz7T_GDsALgDPEWrN4MhF1P7MudfdjqUTnhIxMSo5vkqrUdk8NpCFe6ypykfj-c0tXva_S67FrJP0G03ntfE6hptIxSc0d296PTI-WVFY-Sg8BYNN2JO0JuoqMVCSY-ytW43_0gp3FUGYZhaOMPInDR_1gc80DlY7s7vlbVV1dEC4YBnlYYCr3euiTW?width=657&height=362&cropmode=none" width="657" height="362" style="zoom:67%;" />


- The environmental complexity of the sea and sky makes communication more difficult.
- The lack of base station support in this scenario necessitates a change in the communication mode.
- The search target area is enlarged, but the individual energy is limited, and cluster head is needed for management.
- It is necessary to study the dynamic network topology to cope with the constantly changing environment of airspace and sea area.

### Proposal:

- KubeEdge has the ability to manager the complex wireless networking status without sensing by service.
  - Allow users to be compatible with different network connection mode, under plug-in management
  - With the change of networking modes, the traditional communication modes of KubeEdge need to be extended, such as including broadcast and multi-cast.
- Optimize the quality of communication under the scenario of side collaboration
  - Topologies will constantly changing, and the scheduling of communication capabilities will become a critical issue.
  - Consider the situation of poor network channel quality.
  - Actively consider the network instability caused by mobility, especial for users joint and exit suddenly.
  - Active feedback to adjust the network, active networking.
- Edge nodes have partial cloud capabilities
  - Enable the edge devices to self-organize network and partial autonomy when off-line or discontinuous network. 
  - In most cases, multiple edge devices may spontaneously form a small group, and the outside nodes can join or quit the group.

### Areas of Focus

   Wireless communication tech, reenforcement learning, multi-access protocol, Invoke the schema on the communication computing topology





 
