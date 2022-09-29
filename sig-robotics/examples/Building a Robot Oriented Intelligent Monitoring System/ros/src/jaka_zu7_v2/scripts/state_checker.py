#! /usr/bin/env python3

import rospy
import moveit_msgs.srv
import moveit_commander
import geometry_msgs.msg


class StateChecker:
    def __init__(self):
        self.state_valid_srv = rospy.ServiceProxy('/check_state_validity', moveit_msgs.srv.GetStateValidity)
        self.state_valid_srv.wait_for_service()

        self.state_valid_req = moveit_msgs.srv.GetStateValidityRequest()
        self.state_valid_req.robot_state.joint_state.name = []
        for i in range(1, 7):
            self.state_valid_req.robot_state.joint_state.name.append(f"joint_{i:d}")
        self.state_valid_req.group_name = 'arm'

        self.single = None
        self.scene = moveit_commander.PlanningSceneInterface()
        self.obj_name = "obj"
        rospy.sleep(1.0)

    def update_obj_position(self, pos: list):
        obj_pose = geometry_msgs.msg.PoseStamped()
        obj_pose.header.frame_id = "base_link"
        obj_pose.pose.position.x = pos[0]
        obj_pose.pose.position.y = pos[1]
        obj_pose.pose.position.z = pos[2]
        obj_pose.pose.orientation.w = 1
        self.scene.add_sphere(self.obj_name, obj_pose, radius=0.02)
        if not self._wait_for_scene_update(self.obj_name, 5.0):
            rospy.logwarn("cannot update obj")

    def _wait_for_scene_update(self, obj_name: str, timeout=2.0):
        start = rospy.get_time()
        seconds = rospy.get_time()
        while (seconds - start < timeout) and not rospy.is_shutdown():
            is_known = obj_name in self.scene.get_known_object_names()
            if is_known:
                return True
            # Sleep so that we give other threads time on the processor
            rospy.sleep(0.1)
            seconds = rospy.get_time()
        return False

    def is_valid(self, joint_vals: list) -> bool:
        self.state_valid_req.robot_state.joint_state.position = joint_vals
        resp: moveit_msgs.srv.GetStateValidityResponse = self.state_valid_srv(self.state_valid_req)
        return resp.valid

    def is_current_valid(self) -> bool:
        if not self.single:
            self.single = moveit_commander.MoveGroupCommander('arm')
        state = self.single.get_current_state()
        self.state_valid_req.robot_state = state
        resp: moveit_msgs.srv.GetStateValidityResponse = self.state_valid_srv(self.state_valid_req)
        return resp.valid


if __name__ == "__main__":
    rospy.init_node("checker_node")

    checker = StateChecker()

    checker.update_obj_position([0.3, 0.4, 0.2])
    rospy.sleep(2.0)
    checker.update_obj_position([0.3, 0.4, 0.5])

    print("current:", checker.is_current_valid())
    print(checker.is_valid([1.25, -0.27, 2.16, 0.76, 0.10, 0.00]))
    print(checker.is_valid([1.31, -0.27, 2.10, 0.76, 0.10, 0.00]))
