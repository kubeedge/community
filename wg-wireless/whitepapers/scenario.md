KubeEdge wireless scenario

#### Scenario：

##### 	Vehicle Network/Vessel network/UAV network:

<img src="https://snz04pap002files.storage.live.com/y4m9TFrXL6j2ajyTvHbGHvTWUo5w6kT7SMZ2DZq1t1ivm9Trt4Lo7wt4-azpzI5O746g8fXqBkMz0iweMjyNQdBS4WaXnv5HLtyxU_kHOFNGrEpgfrbocPHwC1Jb6TnskFPTfWwvQUzkhgsz-AoRc1VffTLSUX8klBth8EMHG_KmYO5zoKt_Vpk8XVRDn2UO_-J?width=640&height=397&cropmode=none" width="640" height="397" />

**Key differences** between wireless and wire:  

- Communication mode may need to change from TCP, IP mode to multi-cast or broadcast mode.
- Network for dynamics scenarios.  

With the development of computing, powerful computing equipment could be seen as multi-agent. Therefore, the  problem of wireless interaction between multiple agents is the main scenario discussed here. 

**Description**:  In the mobile scenario, multiple vehicles support similar service, and each vehicles wireless connected with each other as a KubeEdge node. (e.g. NIO ET7 with NVIDIA Orin * 4 , which has computing power of 1016 TOPS)

- Off-line autonomy and node management: Five cars form an mesh network, which can trans information processing between each other in off-line conditions.
- Leader Selection: Choose a Kube as cluster head, to do the overall management of other equipment.
- KubeEdge-wireless should monitor the whole or partly state of network, to aid in networking decisions.
- According to the Service Level Agreement, KubeEdge makes decisions on the networking mode and limits the networking scope.
- Inspired by 3GPP 36.885 standards.



##### **Air/Sea rescue Collaboration：**

<img src="https://snz04pap002files.storage.live.com/y4mJdEgovz7T_GDsALgDPEWrN4MhF1P7MudfdjqUTnhIxMSo5vkqrUdk8NpCFe6ypykfj-c0tXva_S67FrJP0G03ntfE6hptIxSc0d296PTI-WVFY-Sg8BYNN2JO0JuoqMVCSY-ytW43_0gp3FUGYZhaOMPInDR_1gc80DlY7s7vlbVV1dEC4YBnlYYCr3euiTW?width=657&height=362&cropmode=none" width="657" height="362" />

- The environmental complexity of the sea and sky makes communication more difficult
- The lack of base station support in this scenario necessitates a change in the communication mode.
- The search area is enlarged, but the individual energy is limited, and cluster head is needed for management.
- It is necessary to study the dynamic network structure topology to cope with the constantly changing environment of airspace and sea area.

 
