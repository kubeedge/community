#!/usr/bin/env python3
import rclpy
from rclpy.node import Node
from geometry_msgs.msg import Twist
from nav_msgs.msg import Odometry

import math
from typing import List

from promptulate.utils.color_print import print_text

from bot_description.core import (
    RobotAgent,
    RobotController,
    RobotObserver,
)
from bot_description.schema import Operator


def get_operators(node: Node):
    """Get all operators of the robot"""
    operators: List[Operator] = []

    def go_front_callback(distance: float, speed: float = 0.5):
        target_distance: float = node.total_distance + distance
        print_text(
            f"[odom] current distance {node.total_distance} target distance {target_distance}"
        )

        cmd_vel = Twist()
        cmd_vel.linear.x = speed

        while node.total_distance < target_distance:
            # broadcast to send message
            if not getattr(node, "publisher", None):
                node.publisher = node.create_publisher(Twist, "cmd_vel", 10)
            node.publisher.publish(cmd_vel)

            print_text(
                f"[odom] current distance {node.total_distance} target distance {target_distance}",
                "yellow",
            )
            rclpy.spin_once(node)

        print_text("command go front has finish", "green")

    operators.append(
        Operator(
            name="go_front", description="robot go front", callback=go_front_callback
        )
    )

    def go_back_callback(distance: float, speed: float = -0.5):
        target_distance: float = node.total_distance + distance
        print_text(
            f"[odom] current distance {node.total_distance} target distance {target_distance}",
            "yellow",
        )

        cmd_vel = Twist()
        cmd_vel.linear.x = speed

        while node.total_distance < target_distance:
            # broadcast to send message
            node.publisher.publish(cmd_vel)

            print_text(
                f"[odom] current distance {node.total_distance} target distance {target_distance}",
                "yellow",
            )
            rclpy.spin_once(node)

        print_text("command go back has finish", "green")

    operators.append(
        Operator(name="go_back", description="robot go back", callback=go_back_callback)
    )

    def turn_left_callback(angle: float):
        target_yaw: float = math.degrees(node.yaw) + angle
        print_text(
            f"[odom] current yaw {math.degrees(node.yaw)} target yaw {target_yaw}",
            "yellow",
        )

        # broadcast to send message
        cmd_vel = Twist()
        cmd_vel.angular.z = 0.5

        while target_yaw > math.degrees(node.yaw):
            node.publisher.publish(cmd_vel)

            print_text(
                f"[odom] current yaw {math.degrees(node.yaw)} target yaw {target_yaw}",
                "yellow",
            )
            rclpy.spin_once(node)

        print_text("command turn left has finish", "green")

    operators.append(
        Operator(
            name="turn_left",
            description="Turn left in place, No displacement",
            callback=turn_left_callback,
        )
    )

    def turn_right_callback(angle: float):
        target_yaw: float = math.degrees(node.yaw) - angle
        print_text(
            f"[odom] current yaw {math.degrees(node.yaw)} target yaw {target_yaw}",
            "yellow",
        )

        # broadcast to send message
        cmd_vel = Twist()
        cmd_vel.angular.z = -0.5

        while target_yaw < math.degrees(node.yaw):
            node.publisher.publish(cmd_vel)

            print_text(
                f"[odom] current yaw {math.degrees(node.yaw)} target yaw {target_yaw}",
                "yellow",
            )
            rclpy.spin_once(node)

        print_text("command turn right has finish", "green")

    operators.append(
        Operator(
            name="turn_right",
            description="Turn right in place, No displacement",
            callback=turn_right_callback,
        )
    )

    def stop_callback():
        # broadcast to send message
        cmd_vel = Twist()

        if not getattr(node, "publisher", None):
            node.publisher = node.create_publisher(Twist, "cmd_vel", 10)
        node.publisher.publish(cmd_vel)

    operators.append(
        Operator(name="stop", description="stop the robot", callback=stop_callback)
    )

    return operators


class UserClientNode(Node):
    def __init__(self):
        super().__init__("user_client")
        print_text("user_client node startup", "green")

        controller = RobotController(get_operators(self))
        observer = RobotObserver([])
        self.robot_agent = RobotAgent(controller, observer)

        self.total_distance = 0.0
        self.previous_x, self.previous_y = 0.0, 0.0
        self.yaw, self.roll, self.pitch = 0.0, 0.0, 0.0

        self.publisher = self.create_publisher(Twist, "cmd_vel", 10)
        self.subscription = self.create_subscription(
            Odometry, "odom", self.odom_callback, 10
        )
        self.subscription

    def odom_callback(self, msg):
        x = msg.pose.pose.position.x
        y = msg.pose.pose.position.y
        delta_distance: float = math.sqrt(
            (x - self.previous_x) ** 2 + (y - self.previous_y) ** 2
        )
        self.total_distance += delta_distance
        self.previous_x = x
        self.previous_y = y

        orientation = msg.pose.pose.orientation
        self.roll, self.pitch, self.yaw = self.quaternion_to_euler(
            orientation.x, orientation.y, orientation.z, orientation.w
        )

    def quaternion_to_euler(self, x, y, z, w):
        t0 = +2.0 * (w * x + y * z)
        t1 = +1.0 - 2.0 * (x * x + y * y)
        roll = math.atan2(t0, t1)

        t2 = +2.0 * (w * y - z * x)
        t2 = +1.0 if t2 > +1.0 else t2
        t2 = -1.0 if t2 < -1.0 else t2
        pitch = math.asin(t2)

        t3 = +2.0 * (w * z + x * y)
        t4 = +1.0 - 2.0 * (y * y + z * z)
        yaw = math.atan2(t3, t4)

        return roll, pitch, yaw

    def run(self):
        """Startup RobotAgent"""
        while True:
            user_input: str = input("Please input your demand: ")
            self.robot_agent.run(user_input)


def main(args=None):
    rclpy.init(args=args)
    node = UserClientNode()
    node.run()
    rclpy.shutdown()
