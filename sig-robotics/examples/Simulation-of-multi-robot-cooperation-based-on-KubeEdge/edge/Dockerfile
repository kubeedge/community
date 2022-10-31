# enviroment core


ARG OWNER=326737089@qq.com
ARG ROOT_CONTAINER=$OWNER/${ROS_DISTRO}-enviroment-core

FROM ubuntu as intermediate

# install git
RUN apt-get update
RUN apt-get install -y git

# add credentials on build
ARG SSH_PRIVATE_KEY
RUN mkdir /root/.ssh/
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

# Add git providers to known_hosts
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

RUN git clone --branch master  git@github.com:tuw-robotics/rosaria.git             p3dx/src/rosaria
RUN git clone --branch master  git@github.com:tuw-robotics/tuw_p3dx.git            p3dx/src/tuw_p3dx

FROM $ROOT_CONTAINER
USER root

ENV PRJ_DIR=${MY_HOME}/projects/p3dx

RUN apt-get update && apt-get -y upgrade

# Install locale and set
RUN apt-get install -y iputils-ping bash-completion

COPY --chown=${MY_USER}:${MY_USER} --from=intermediate /p3dx ${PRJ_DIR}
COPY --chown=${MY_USER}:${MY_USER} --from=intermediate /root/.ssh ${MY_HOME}/.ssh

RUN /bin/bash -c '. /opt/ros/${ROS_DISTRO}/setup.bash;  cd ${PRJ_DIR}/;  rosdep update;  rosdep install -i --from-path src --rosdistro noetic -y'
# ros-${ROS_DISTRO}-hokuyo-node
USER ${MY_USER}


    
RUN /bin/bash -c '. /opt/ros/${ROS_DISTRO}/setup.bash;             cd ${PRJ_DIR}; catkin_make'
RUN echo 'source ${PRJ_DIR}/devel/setup.bash' >> ~/.bashrc
RUN echo 'roslaunch tuw_p3dx aria_twist.launch'  >> ~/.bash_history
# clean apt
USER root
RUN rm -rf ${MY_HOME}/.ros
#RUN apt-get autoclean && apt-get autoremove && rm -rf /var/lib/apt/lists/*

USER ${MY_USER}
WORKDIR ${PRJ_DIR}
ENV HOME=${MY_HOME}
COPY --chown=${MY_USER}:${MY_USER} ./edge.sh /
ENTRYPOINT ["/edge.sh"]
CMD [/bin/bash]
