# Generating Action Sequences for Robot Task Plans using LLM

## Environment

- Ubuntu 22.04
- ROS2 humble
- gpt-3.5-turbo
- gazebo

## Setup

**Clone the community/sig-robotics repository**

```shell
git clone https://github.com/kubeedge/community
```

**Install relational packages and build the project**

```shell
cd ./sig-robotics/examples/Generating-Action-Sequences-for-Robot-Task-Plans-using-LLM/
rosdep install --from-paths src -y

cd /src/bot_description/
pip3 install -r requirements.txt

cd ./sig-robotics/examples/Generating-Action-Sequences-for-Robot-Task-Plans-using-LLM/
colcon build
```

**Simulation in gazebo**

```shell
source install/setup.bash
ros2 launch bot_description gazebo.launch.py
ros2 run bot_description user_client
```
