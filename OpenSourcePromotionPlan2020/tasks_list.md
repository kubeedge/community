## Tasks

* [1. Build an AI example to run AI application on KubeEdge platform](#1-build-an-ai-example-to-run-ai-application-on-kubeedge-platform)
* [2. Add script to install KubeEdge cloud and edge side](#2-add-script-to-install-kubeedge-cloud-and-edge-side)
* [3. Integrate chaosMesh to test kubeedge](#3-integrate-chaosmesh-to-test-kubeedge)
* [4. New theme and format for the website](#4-new-theme-and-format-for-the-website)
* [5. Support list-watch from edgecore for applications on the edge](#5-support-list-watch-from-edgecore-for-applications-on-the-edge)

### 1. Build an AI example to run AI application on KubeEdge platform

- **Project Description**: Build a sample example to run AI application on KubeEdge platform, using the device API, and ideally running inference on the edge and training on the cloud, showing cloud-edge collaboration.

- **Difficulty**: medium

- **Mentor**: Kevin Wang(@kevin-wangzefeng): kevinwzf0126@gmail.com

- **Output**:
   1. Build an example to run AI application
   2. The device API 

- **Tech. Requirement**:
   1. KubeEdge
   2. AI technology, such as face recognition

- **Related Repo**:
   1. https://github.com/kubeedge/kubeedge
   2. https://github.com/kubeedge/examples


### 2. Add script to install KubeEdge cloud and edge side

- **Project Description**: Write a script to install KubeEdge cloud and edge with one click, including K8s, Runtime. **Note**: Need more discussion.

- **Difficulty**: low

- **Mentor**: Fei Xu(@fisherxu): fisherxu1@gmail.com

- **Output**: 
   1. Write a script to install KubeEdge 

- **Tech. Requirement**:
   1. Shell
   2. Kubernetes
   3. KubeEdge

- **Related Repo**:
   1. https://github.com/kubeedge/kubeedge

### 3. Integrate chaosMesh to test kubeedge

- **Project Description**: Now we have unit tests, integration tests and e2e tsets that guarantee the KubeEdge is running well, but to better identify system vulnerabilities and improve resilience, let's use chaosMesh to injects various types of faults into the KubeEdge, like low quality network, etc. to test KubeEdge.

- **Difficulty**: Medium

- **Mentor**: Kevin Wang(@kevin-wangzefeng): kevinwzf0126@gmail.com

- **Output**:
   1. Test the reliability and stability of the system to ensure production availability.
   2. Output the report of the test result

- **Tech. Requirement**:
   1. KubeEdge
   2. chaos-mesh

- **Related Repo**:
   1. https://github.com/pingcap/chaos-mesh
   2. https://github.com/kubeedge/kubeedge

### 4. New theme and format for the website

- **Project Description**: Redesign and reimplement the official website to show the key features of KubeEdge.

- **Difficulty**: Medium

- **Mentor**: Kevin Wang(@kevin-wangzefeng): kevinwzf0126@gmail.com

- **Output**:
   1. New website styling
   2. Key values/features of KubeEdge highlighted in prominent position
   3. redundant content section removed in homepage 
   4. add support for hosting multiple releases and multi-languages for documents
   5. add community information including partners, end users, success stories.


- **Tech. Requirement**:
   1. KubeEdge
   2. css
   3. website layout

- **Related Repo**:
   1. https://github.com/kubeedge/website

### 5. Support list-watch from edgecore for applications on the edge

- **Project Description**: Some applications running on the edge side need to connect to the k8s master through list-watch interface, but on the edge it cannot directly connect to the k8s master. So we need forward the list-watch requests to the k8s master through the cloud-side channel.


- **Difficulty**: Hard

- **Mentor**: Kevin Wang(@kevin-wangzefeng): kevinwzf0126@gmail.com

- **Output**:
   1. Implement the list-watch interface on the edge side 
   2. Forward the list-watch requests and get the data through the reliable transmission channel

- **Tech. Requirement**:
   1. KubeEdge
   2. Kubernetes

- **Related Repo**:
   1. https://github.com/kubeedge/kubeedge
   2. https://github.com/kubernetes/kubernetes
