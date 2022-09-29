/**
 * 该例程将发布velocity_publisher/cmd_vel话题，消息类型geometry_msgs::Twist
 */
#include <ros/ros.h>
#include <geometry_msgs/Twist.h>

int main(int argc, char **argv)
{
	// ROS节点初始化
	ros::init(argc, argv, "velocity_publisher");

	// 创建节点句柄
	ros::NodeHandle n;

	// 创建一个Publisher，发布名为velocity_publisher/cmd_vel的topic，消息类型为geometry_msgs::Twist，队列长度10
	ros::Publisher turtle_vel_pub = n.advertise<geometry_msgs::Twist>("/cmd_vel", 10);

	// 设置循环的频率
	ros::Rate loop_rate(10);

	int count = 0;
	while (ros::ok())
	{
	    // 初始化geometry_msgs::Twist类型的消息
        char c1 = getchar();
        char c2 = getchar();
        printf("%d\n", c1);

		geometry_msgs::Twist vel_msg;
        switch(c1){
            case 119:
                //w
                vel_msg.linear.x = 1;
                // 发布消息
                turtle_vel_pub.publish(vel_msg);
                break;
            case 97:
                //a
                vel_msg.angular.z = 1;
                // 发布消息
                turtle_vel_pub.publish(vel_msg);
                sleep(1);
                break;
            case 115:
                //s
                vel_msg.linear.x = -1;
                // 发布消息
                turtle_vel_pub.publish(vel_msg);
                sleep(1);
                break;
            case 100:
                //d
                vel_msg.angular.z = -1;
                // 发布消息
                turtle_vel_pub.publish(vel_msg);
                sleep(1);
                break;
            default:
                printf("错误");
                c1 = 0;
                c2 = 0;
                vel_msg.linear.x = 0;
                vel_msg.linear.y = 1;
                vel_msg.linear.z = 0;
                vel_msg.angular.x = 0;
                vel_msg.angular.y = 0;
                vel_msg.angular.z = 0;
                turtle_vel_pub.publish(vel_msg);
        }
        sleep(1);

        vel_msg.linear.x = 0;
        vel_msg.linear.y = 0;
        vel_msg.linear.z = 0;
        vel_msg.angular.x = 0;
        vel_msg.angular.y = 0;
        vel_msg.angular.z = 0;
        turtle_vel_pub.publish(vel_msg);
	    // 按照循环频率延时
	    loop_rate.sleep();

	}

	return 0;
}
