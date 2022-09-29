# Building a Robot-Oriented Intelligent Monitoring System

## Background

With the continuous development of robotics and cloud computing technology, researchers hope to combine these two technologies, so in 2010, Professor James Kuffer first proposed the concept of cloud robot. The cloud robot can be regarded as a combination of the cloud and the robot side, which can realize real-time monitoring and scheduling management of the robot, and can greatly improve the performance of the robot.

During the operation of the robot, the running state of the robot changes all the time, monitoring the running state of the robot, on the one hand, it is convenient for the operator to schedule the robot, and the performance of different robots can be maximized by distributing different tasks. optimal. On the other hand, when the robot runs abnormally, the operator can repair and maintain it in time. Therefore, it is necessary to develop a robot-oriented intelligent monitoring system.

Based on KubeEdge, KubeEdge SIG Robotics focuses on the field of cloud robotics, integrating cloud computing technology with robotics, especially the robot management system based on the cloud-edge collaborative architecture, empowering the existing robot ecosystem. Based on the edge-cloud collaborative architecture, the robot's intelligence and development efficiency are improved. Specifically, KubeEdge SIG Robotics supports heterogeneous robots, including mobile robots, robotic arms and many other robotic platforms. It is possible to cooperate with other open source communities such as ROS, Gazebo and also with other SIGs such as SIG AI and SIG IoT. Makes robot development, debugging, simulation, deployment and management easier.

However, at present, KubeEdge SIG Robotics lacks the visual monitoring function for robots, so this project will develop a visual intelligent monitoring system for robots in the robot virtual machine simulation environment provided by KubeEdge SIG Robotics based on KubeEdge, prometheus, grafana and other open source components. Provide a reference case for robot developers and users.

## Goals

For the developers and users of the robot intelligent monitoring system, the goal of this project is to provide developers and users with a practical case for reference in the process of robot development based on KubeEdge SIG Robotics, including:

- End-to-end test cases (monitoring indicators completed 20+, beautiful interface, monitoring delay less than 3s)
 
- System Deployment and User Manual

## Proposal

### Project scope

Expand the usage scenarios of KubeEdge SIG Robotics, realize the visual monitoring function for robots, and provide end-to-end test cases to improve the development and deployment efficiency of robots.

- Robots in a virtual simulation environment for acquiring monitoring data, such as data collection of key monitoring indicators for robotic platforms such as robotic arms and mobile robots.

- An intelligent visual monitoring system interface for displaying monitoring data, displaying collected robot monitoring data, such as CPU, memory, laser ranging, robot moving speed, acceleration, etc.

- The deployment and use manual of the system is convenient for developers to carry out customized secondary development based on this project.

### Targeting users

- Developer: Building a complete solution for a robot intelligent monitoring system based on KubeEdge SIG Robotics

- End user: Implement the function of monitoring the real-time status of the robot.

## Design Details

### Project Architecture

For this project, the system is divided into three modules: **data acquisition, data monitoring, and data display**. The architecture diagram of this project is shown in the figure below.

![Image text](https://github.com/ycr-sjtu/community/blob/master/sig-robotics/propoasal/Building%20a%20Robot%20Oriented%20Intelligent%20Monitoring%20System/images/System%20Architecture.png)

- **Data acquisition module**Based on ROS development, it is divided into real robots and simulated robots in gazebo. The robot monitoring data is obtained in real time through the corresponding ROS node and transmitted to the cloud.

- **Data monitoring module**Developed based on Prometheus. Through the exporter, the monitoring data of the robot in the virtual environment can be collected, and then sent to the prometheus server. Retrieval is responsible for grabbing monitoring indicator data on the active target host. Storage storage mainly stores the collected data to the disk. PromQL is a query language module provided by Prometheus. After receiving the alarm information from the Prometheus server, it will deduplicate, group, and route to the corresponding receiver to issue an alarm. Common receiving methods are: email, WeChat, DingTalk, etc.

- **Data display module**Developed based on grafana. Visualize monitoring data through dashboards.

After running the project instance, when the robot state is changed in the simulation environment, the robot index changes can be observed on the monitoring interface. When the specified conditions are triggered (such as the distance is too close, the speed is too fast), the system will automatically send alarm information.

### Robot Monitoring Metrics（22）


**1.Onboard computer（9）**：Current robot host IP, current robot name, memory, memory usage, number of CPU cores, CPU speed, CPU usage, GPU memory, GPU usage

**2.AGV（8）**：Laser sensor (ranging), IMU sensor (acceleration), IMU sensor (angular velocity), running time, position coordinates, moving speed, battery power, ambient temperature

**3.robotic arm（5）**：Manipulator joint pose, manipulator joint torque, manipulator joint speed, manipulator end coordinates, running time

## Roadmap

- July 2022, **project preparation stage**. Learning of open source components such as KubeEdge, prometheus, and grafana.
  - Learn the basic operations of open source components such as KubeEdge, prometheus, and grafana
  - Reproduce existing open source intelligent monitoring system projects
- August 2022, **project development stage**. Build a robot virtual simulation environment to obtain monitoring indicator data. Build a visual interface to realize real-time monitoring of 20+ indicators and alarm functions for abnormal data.
  - Build a robot virtual simulation environment based on ROS and gazebo to obtain robot monitoring index data.
  - Build a robot intelligent monitoring system based on prometheus, capture and store time series data, and set robot abnormal alarm information.
  - Develop the visual interface of intelligent monitoring system based on grafana, configure the system data interface, and improve the interface layout.
- September 2022, **project improvement stage**. Organize project materials and prepare for project review.
  - Includes deployment and usage manuals for open source components, necessary code and configuration files
  - The project PR is merged into the open source community

