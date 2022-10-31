#Building a Robot Oriented Intelligent Monitoring System

This folder contains the mechanical arm and mobile robot model, the control code and the UI.  

##Environment
###Cloud
- Ubuntu 20.04
- ROS neotic
###Edge
- Ubuntu 20.04
- ROS neotic

##Content

###`robot_ Simulation`is the simulation and communication system code of double arm robot 
- In urdf, it is the configuration file of mechanical arm model
- Launch the program for the robot arm simulation environment  
 - Through`roslaunch robot_simulation jaka_rviz.launch`start the virtual environment of the mechanical arm.
- The src contains the simulation system test code and communication system code  
 - `joint_state_publisher_demo` is a communication system test code. After running, the mechanical arm moves, indicating that the environment is built normally.Run the simulation environment before running `rosrun robot_simulation joint_state_publisher_demo`。  

###`chassis_simulation`is the chassis simulation and communication system code
- launchis the startup code of the chassis simulation environment  
 - gazebo: `roslaunch chassis_simulation chassis_simulation_gazebo.launch`  
 - rviz: `roslaunch chassis_simulation chassis_simulation_rviz.launch`  
- The src contains the simulation system test code and communication system code.  
 - `velocity_publisher` is a test code that can control the front and rear of the chassis through wasd. Run the simulation environment before running `rosrun chassis_simulation velocity_publisher`。  
###`jaka_zu7_v2`和`zv7`are the mechanical arm models (no need to change)

###grafana
The `Robot monitoring system-1,664,458,556,895. json` in it is the UI configuration file of grafana. Import grafana to see the Robot monitoring system interface