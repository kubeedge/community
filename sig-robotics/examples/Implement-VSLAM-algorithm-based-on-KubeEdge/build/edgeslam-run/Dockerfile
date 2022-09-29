FROM jike5/orbslam2-dev:latest

LABEL maintainer="2719825969@qq.com"
USER root

# # Clone code
WORKDIR /workspace
RUN git clone https://github.com/jike5/Edgeslam

# Build
RUN cd /workspace/Edgeslam && chmod +x build.sh \
  && chmod +x build_ros.sh \
  && rm -rf ./build \
  && ./build.sh

RUN cd /workspace/Edgeslam/Examples/ROS/Edge_SLAM \
  && rm -rf build \
  && mkdir build

ENV ROS_PACKAGE_PATH "/opt/ros/melodic/share:/workspace/Edgeslam/Examples/ROS"

WORKDIR /workspace/Edgeslam/Examples/ROS/Edge_SLAM/build
RUN bash -c "source /opt/ros/melodic/setup.bash && export ROS_PACKAGE_PATH=/opt/ros/melodic/share:/workspace/Edgeslam/Examples/ROS && cmake .. -DROS_BUILD_TYPE=Release && make -j"

RUN echo "source /workspace/Edgeslam/Examples/ROS/Edge_SLAM/build/devel/setup.bash" >> ~/.bashrc

WORKDIR /workspace/Edgeslam/Examples
RUN chmod +x cloud_run.sh edge_run.sh 
