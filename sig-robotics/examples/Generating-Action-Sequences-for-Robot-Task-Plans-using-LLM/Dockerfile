FROM ros:humble

WORKDIR /ros2_ws

COPY . /ros2_ws/src

RUN apt-get update && apt-get install -y \
    python3-pip \
    python3-rosdep \
    python3-vcstool \
    && rm -rf /var/lib/apt/lists/*


RUN rosdep init && rosdep update
RUN rosdep install --from-paths src -y

RUN colcon build

ENV ROS_DOMAIN_ID=0

CMD ["bash"]
