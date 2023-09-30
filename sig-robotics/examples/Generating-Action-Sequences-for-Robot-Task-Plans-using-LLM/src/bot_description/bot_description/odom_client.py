import rclpy
from rclpy.node import Node
from geometry_msgs.msg import Twist
from nav_msgs.msg import Odometry
from math import sqrt


class OdomUserClientNode(Node):
    def __init__(self):
        super().__init__("go_front")
        self.publisher_ = self.create_publisher(Twist, "cmd_vel", 10)
        self.subscription_ = self.create_subscription(
            Odometry, "odom", self.odometry_callback, 10
        )
        self.subscription_  # prevent unused variable warning
        self.distance_to_travel = 0
        self.distance_traveled = 0

    def run(self):
        while True:
            target_distance: float = float(str(input("Please input dir: ")))
            self.go_front(target_distance)

    def odometry_callback(self, msg):
        self.distance_traveled = sqrt(
            msg.pose.pose.position.x**2 + msg.pose.pose.position.y**2
        )

    def go_front(self, distance):
        self.distance_to_travel = distance
        self.distance_traveled = 0

        twist = Twist()
        twist.linear.x = 0.5  # Adjust the linear velocity as needed
        twist.angular.z = 0.0

        while self.distance_traveled < self.distance_to_travel:
            self.publisher_.publish(twist)
            self.get_logger().info(f"Distance traveled: {self.distance_traveled}")
            rclpy.spin_once(self)

        twist.linear.x = 0.0
        self.publisher_.publish(twist)
        self.get_logger().info("Finished moving forward")


def main(args=None):
    rclpy.init(args=args)
    node = OdomUserClientNode()
    node.run()
    rclpy.shutdown()
