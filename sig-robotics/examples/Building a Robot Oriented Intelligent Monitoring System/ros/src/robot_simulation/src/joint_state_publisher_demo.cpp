/**
 * 此代码用于测试机械臂仿真环境工作是否正常。
 * 打开rviz后，仿真环境会是空白的，因为没有角度输入
 * clion运行此代码，机械臂会开始运动，仿真环境正常
 * 
 */

#include <ros/ros.h>
#include <geometry_msgs/Twist.h>
#include <sensor_msgs/JointState.h>
#include <string>
#include <iostream>
#include <vector>
using namespace std;

int main(int argc, char **argv)
{
	// ROS节点初始化
	ros::init(argc, argv, "joint_state_publisher");

	// 创建节点句柄
	ros::NodeHandle m_nh;

	// 创建一个Publisher，发布名为/joint_states的topic，消息类型为sensor_msgs::JointState，队列长度5
	ros::Publisher m_jointStatePub = m_nh.advertise<sensor_msgs::JointState>("joint_states", 5);

	// 设置循环的频率
	ros::Rate loop_rate(10);

	while (ros::ok())
	{
		for (int i = 0; i < 500; i++)
		{
			// 初始化sensor_msgs::JointState类型的消息
			sensor_msgs::JointState m_jointStateMsg;

			m_jointStateMsg.name.resize(12);
			m_jointStateMsg.position.resize(12);

			m_jointStateMsg.header.stamp = ros::Time::now();
			m_jointStateMsg.name[0] = "left_joint_1";
			m_jointStateMsg.name[1] = "left_joint_2";
			m_jointStateMsg.name[2] = "left_joint_3";
			m_jointStateMsg.name[3] = "left_joint_4";
			m_jointStateMsg.name[4] = "left_joint_5";
			m_jointStateMsg.name[5] = "left_joint_6";
			m_jointStateMsg.name[6] = "right_joint_1";
			m_jointStateMsg.name[7] = "right_joint_2";
			m_jointStateMsg.name[8] = "right_joint_3";
			m_jointStateMsg.name[9] = "right_joint_4";
			m_jointStateMsg.name[10] = "right_joint_5";
			m_jointStateMsg.name[11] = "right_joint_6";

			m_jointStateMsg.position[0] = 0 + (i * 0.01);
			m_jointStateMsg.position[1] = 0;
			m_jointStateMsg.position[2] = 0;
			m_jointStateMsg.position[3] = 0;
			m_jointStateMsg.position[4] = 0;
			m_jointStateMsg.position[5] = 0;
			m_jointStateMsg.position[6] = 0;
			m_jointStateMsg.position[7] = 0;
			m_jointStateMsg.position[8] = 0;
			m_jointStateMsg.position[9] = 0;
			m_jointStateMsg.position[10] = 0;
			m_jointStateMsg.position[11] = 0;

			// 发布消息
			m_jointStatePub.publish(m_jointStateMsg);


			// 按照循环频率延时
			loop_rate.sleep();
		}


	}

	return 0;
}
