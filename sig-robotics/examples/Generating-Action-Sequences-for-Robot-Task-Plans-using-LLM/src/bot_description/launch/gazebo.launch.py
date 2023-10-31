# Copyright 2021 The KubeEdge Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
from launch import LaunchDescription
from launch.actions import ExecuteProcess
from launch_ros.actions import Node
from launch_ros.substitutions import FindPackageShare


def generate_launch_description():
    package_name = "bot_description"
    urdf_name = "bot_gazebo.urdf"
    robot_name_of_model = "bot"

    ld = LaunchDescription()
    pkg_path = FindPackageShare(package=package_name).find(package_name)
    urdf_model_path = os.path.join(pkg_path, f"urdf/{urdf_name}")
    gazebo_world_path = os.path.join(pkg_path, "world/neighborhood.world")

    # startup gazebo simulation
    start_gazebo_cmd = ExecuteProcess(
        cmd=[
            "gazebo",
            "--verbose",
            "-s",
            "libgazebo_ros_init.so",
            "-s",
            "libgazebo_ros_factory.so",
            gazebo_world_path,
        ],
        output="screen",        
    )
    ld.add_action(start_gazebo_cmd)

    # create the robot
    spawn_entity_cmd = Node(
        package="gazebo_ros",
        executable="spawn_entity.py",
        arguments=["-entity", robot_name_of_model, "-file", urdf_model_path],
        output="screen",
        description="Spawn robot entity",        
    )
    ld.add_action(spawn_entity_cmd)

    # startup robot state publisher
    start_robot_state_publisher_cmd = Node(
        package="robot_state_publisher",
        executable="robot_state_publisher",
        arguments=[urdf_model_path],
        description="robot state publisher",
    )
    ld.add_action(start_robot_state_publisher_cmd)

    return ld
