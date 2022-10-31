FROM jike5/ros-melodic:cpu

RUN apt-get update

ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get install -y gnupg2 curl lsb-core vim wget python-pip libpng16-16 libjpeg-turbo8 libtiff5

RUN apt-get install -y \
        # Base tools
        cmake \
        build-essential \
        git \
        unzip \
        pkg-config \
        python-dev \
        # OpenCV dependencies
        python-numpy \
        # Pangolin dependencies
        libgl1-mesa-dev \
        libglew-dev \
        libpython2.7-dev \
        # orbslam2 must use Eigen3.3.2, we will install another later
        # libeigen3-dev \
        apt-transport-https \
        ca-certificates\
        software-properties-common

RUN apt-get update && apt-get install -y \
	# websocketpp dependencies
        libarmadillo-dev\
        libasio-dev\
        # network tools
        iputils-ping net-tools iproute2 telnet

# # build Eigen3
RUN wget https://github.com/eigenteam/eigen-git-mirror/archive/refs/tags/3.3.2.zip \
  && unzip 3.3.2 \
  && cd eigen-git-mirror-3.3.2 \
  && mkdir build \
  && cd build \
  && cmake .. \
  && make install \
  && cp -r /usr/local/include/eigen3/Eigen /usr/local/include

# # Build Pangolin
RUN git clone https://github.com/stevenlovegrove/Pangolin.git \
  && cd Pangolin \
  && git checkout 7987c9b \
  && mkdir build \
  && cd build \
  && cmake .. \
  && make \
  && make install

COPY docker/orbslam2-dev/ros_entrypoint.sh /ros_entrypoint.sh
RUN chmod +x  /ros_entrypoint.sh
ENV ROS_DISTRO melodic
ENV LANG en_US.UTF-8

ENTRYPOINT ["/ros_entrypoint.sh"]

USER $USERNAME
# terminal colors with xterm
ENV TERM xterm
CMD ["bash"]
