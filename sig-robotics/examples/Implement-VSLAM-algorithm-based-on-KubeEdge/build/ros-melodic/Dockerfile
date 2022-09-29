# Refer: https://github.com/turlucode/ros-docker-gui

FROM ubuntu:18.04

LABEL Description="ROS-Melodic-Desktop (Ubuntu 18.04)"

# Install packages without prompting the user to answer any questions
ENV DEBIAN_FRONTEND noninteractive 

# Install packages
RUN apt-get update && apt-get install -y \
locales \
lsb-release \
mesa-utils \
git \
subversion \
nano \
terminator \
xterm \
wget \
curl \
htop \
libssl-dev \
build-essential \
dbus-x11 \
software-properties-common \
gdb valgrind && \
apt-get clean && rm -rf /var/lib/apt/lists/*

# Install cmake 3.22.1
RUN git clone https://gitlab.kitware.com/cmake/cmake.git && \
cd cmake && git checkout tags/v3.22.1 && ./bootstrap --parallel=8 && make -j8 && make install && \
cd .. && rm -rf cmake

# Install tmux 3.2
RUN apt-get update && apt-get install -y automake autoconf pkg-config libevent-dev libncurses5-dev bison && \
apt-get clean && rm -rf /var/lib/apt/lists/*
RUN git clone https://github.com/tmux/tmux.git && \
cd tmux && git checkout tags/3.2 && ls -la && sh autogen.sh && ./configure && make -j8 && make install

# Install new paramiko (solves ssh issues)
RUN apt-add-repository universe
RUN apt-get update && apt-get install -y python-pip python build-essential && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN /usr/bin/yes | pip install --upgrade "pip < 21.0"
RUN /usr/bin/yes | pip install --upgrade virtualenv
RUN /usr/bin/yes | pip install --upgrade paramiko
RUN /usr/bin/yes | pip install --ignore-installed --upgrade numpy protobuf

# Locale
RUN locale-gen en_US.UTF-8  
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en  
ENV LC_ALL en_US.UTF-8
ENV ROS_DISTRO melodic

# Terminator Config
RUN mkdir -p /root/.config/terminator/
COPY docker/ros-melodic/terminator_config /root/.config/terminator/config 
COPY docker/ros-melodic/terminator_background.png /root/.config/terminator/background.png

# Install ROS
RUN sh -c 'echo "deb http://packages.ros.org/ros/ubuntu $(lsb_release -sc) main" > /etc/apt/sources.list.d/ros-latest.list'
RUN apt-key adv --keyserver 'hkp://keyserver.ubuntu.com:80' --recv-key C1CF6E31E6BADE8868B172B4F42ED6FBAB17C654

RUN apt-get update && apt-get install -y --allow-downgrades --allow-remove-essential --allow-change-held-packages \
libpcap-dev \
libopenblas-dev \
gstreamer1.0-tools libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-good1.0-dev \
ros-melodic-desktop-full python-rosinstall python-rosinstall-generator python-wstool build-essential python-rosdep \
ros-melodic-socketcan-bridge \
ros-melodic-geodesy && \
apt-get clean && rm -rf /var/lib/apt/lists/*

# Configure ROS
RUN rosdep init && rosdep update
RUN echo "source /opt/ros/melodic/setup.bash" >> /root/.bashrc
RUN echo "export ROSLAUNCH_SSH_UNKNOWN=1" >> /root/.bashrc

# Entry script - This will also run terminator
COPY docker/ros-melodic/entrypoint_setup.sh /
ENTRYPOINT ["/entrypoint_setup.sh"]

# ---
CMD ["terminator"]
