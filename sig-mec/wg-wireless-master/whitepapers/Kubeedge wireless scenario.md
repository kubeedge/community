Kubeedge wireless scenario

#### Scenario：

##### 	Vehicle Network/Vessel network/UAV network:

<img src="../JPG/image-20210203172851933.png" alt="image-20210203172851933" style="zoom:50%;" />

**Key differences** between wireless and wire:  

- Communication mode may need to change from TCP, IP mode to multi-cast or broadcast mode.
- Network for dynamics scenarios.  

With the development of computing, powerful computing equipment could be seen as multi-agent. Therefore, the  problem of wireless interaction between multiple agents is the main scenario discussed here. 

**Description**:  In the mobile scenario, multiple vehicles support similar service, and each vehicles wireless connected with each other as a KubeEdge node. (e.g. NIO ET7 with NVIDIA Orin * 4 , which has computing power of 1016 TOPS)

- Off-line autonomy and node management: Five cars form an mesh network, which can trans information processing between each other in off-line conditions.
- Leader Selection: Choose a Kube as cluster head, to do the overall management of other equipments.
- KubeEdge-wireless should monitor the whole or partly state of network, to aid in networking decisions.
- According to the Service Level Agreement, KubeEdge makes decisions on the networking mode and limits the networking scope.
- Inspired by 3GPP 36.885 standards.



##### **Air/Sea rescue Collaboration：**

<img src="../JPG/Kube-Wireless流程图-海洋搜救.png" alt="Kube-Wireless流程图-海洋搜救" style="zoom: 67%;" />

- The environmental complexity of the sea and sky makes communication more difficult
- The lack of base station support in this scenario necessitates a change in the communication mode.
- The search area is enlarged, but the individual energy is limited, and cluster head is needed for management.
- It is necessary to study the dynamic network structure topology to cope with the constantly changing environment of airspace and sea area.

 
