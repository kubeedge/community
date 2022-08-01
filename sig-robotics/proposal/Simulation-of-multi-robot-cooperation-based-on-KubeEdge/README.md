## 1 Motivation

The multi-robot system has become the development trend of the robotics industry. The multi-robot system has better robustness than a single robot and has a broader application space. The correctness and effectiveness of the multi-robot coordination and control algorithm can be verified by simulating the process of multi-robot system completing tasks in the environment through computer simulation.

### 1.1 Goals

- Add a park model to Gazebo, and add a variety of robots to cooperate to complete the inspection tasks of the entire park. And based on KubeEdge for multi-robot scheduling, planning and collaborative simulation.
- Robot : quadrotor UAV, wheeled robot;
- The cloud is responsible for the task distribution of each robot, as well as dynamic multi-robot scheduling. The edge(robots) feeds back its own status information in real time.

## 2 Proposal

### 2.1 Design Architecture

![image-20220801145913088](images/image-20220801145913088.png)

### Cloud

- The cloud receives the position, speed, and other running status information of each robot at the edge in real time, which is used as the basis for task issuance and dynamic task adjustment.
- According to the scheduling strategy, the cloud sends tasks to the edge end: the location of the target point and the planned inspection path. And according to the status information reported by each robot in real time, it is adjusted in real time.
- Scheduling strategy: like a multi-trip dealer problem (maintain a dynamic task list, update the scheduled task of each robot according to the working status of the robot)

### Edge

- Based on the target point information and path information sent by the cloud, the edge robot moves in the Gazebo environment by simulating the motion controller.
- In the process of moving, The robot is equipped with a lidar sensor to collect information about the surrounding environment, and combined with their own pose, speed to conduct  local path planning.
- The edge robot reports its own status information in real time during the process of moving.

