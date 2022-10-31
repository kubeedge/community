## ros-melodic image

This image provide a ROS iamge with GUI sipport. We referenced this [repository](https://github.com/turlucode/ros-docker-gui), but still made a lot of improvements to better match our application scenarios.

You can use the Dockerfile build the ros-melodic image:

```
docker build . -t jike5/ros-melodic:cpu -f docker/ros-melodic/Dockerfile 
```

