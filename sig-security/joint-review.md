# Joint-review of KubeEdge

This joint-review relied heavily on the [self-assessment](https://github.com/kubeedge/community/tree/master/sig-security/self-assessment.md) provided by the KubeEdge project.

## Table of Contents

* [Metadata](#metadata)
  * [Security links](#security-links)
* [Overview](#overview)
  * [Background](#background)
  * [Goals](#goals)
  * [Non-goals](#non-goals)
* [Joint-review use](#joint-review-use)
* [Intended use](#intended-use)
* [Project design](#project-design)
  * [Functions and features](#functions-and-features)
    * [Security functions and features](#security-functions-and-features)
* [Configuration and set-up](#configuration-and-set-up)
* [Project compliance](#project-compliance)
* [Security analysis](#security-analysis)
* [Secure development practices](#secure-development-practices)
* [Security issue resolution](#security-issue-resolution)
  * [Closed security issues and vulnerabilities](#closed-security-issues-and-vulnerabilities)
* [Hands-on review](#hands-on-review)
* [Roadmap](#roadmap)
* [Appendix](#appendix)

## Metadata

|                   |                                                              |
| ----------------- | ------------------------------------------------------------ |
| Software          | https://github.com/kubeedge/kubeedge                         |
| Website           | https://kubeedge.io                                          |
| Security Provider | No                                                           |
| Incubation PR     | https://github.com/cncf/toc/pull/461                         |
| Languages         | Go                                                           |
| SBOM              | Check [go.mod](https://github.com/kubeedge/kubeedge/blob/master/go.mod) for libraries, packages, versions used by the project |

### Security links

| Doc                          | url                                                          |
| ---------------------------- | ------------------------------------------------------------ |
| Security file                | https://github.com/kubeedge/community/tree/master/sig-security<br />https://github.com/kubeedge/community/tree/master/team-security |
| Default and optional configs | [cloudcore](https://github.com/kubeedge/kubeedge/blob/master/pkg/apis/componentconfig/cloudcore/v1alpha1/default.go)<br />[edgecore](https://github.com/kubeedge/kubeedge/blob/master/pkg/apis/componentconfig/edgecore/v1alpha1/default.go) |

## Overview

KubeEdge is an open source system for extending native containerized application orchestration capabilities to hosts at the edge. It's built upon Kubernetes and provides fundamental infrastructure support for networking, application deployment and metadata synchronization between cloud and edge.

Since joining CNCF, KubeEdge has attracted more than [1000+ Contributors](https://kubeedge.devstats.cncf.io/d/18/overall-project-statistics-table?orgId=1) from 80+ different Organizations with 3,900+ Commits, and got 5,600+ Stars on Github and 1,600+ Forks. It has been adopted by China Mobile, China Telecom, Raisecom, WoCloud, Xinghai IoT, KubeSphere, HUAWEI CLOUD, Harmony Cloud, DaoCloud, SAIC Motor, SF Express, etc.

### Background

As 5G and AI technologies grow, enterprises are increasingly demanding intelligent upgrade, and the application scenarios of edge computing are becoming more and more extensive. It is challenging to realize unified management and control of edge computing resources and cloud-edge synergy, such as how to run intelligent applications and algorithms on edge devices with limited resources (e.g. cameras and drones), how to solve the problems caused by the access of mass heterogeneous edge devices in intelligent transportation and intelligent energy, and how to ensure high reliability of services in off line scenarios.

KubeEdge provides solutions for cloud-edge synergy and has been widely adopted in industries including transportation, energy, Internet, CDN, manufacturing, smart campus, etc. KubeEdge provides:

- Seamless cloud edge communication for both metadata and data.
- Edge autonomy: autonomous operation of edge even during disconnection from cloud.
- Low resource readiness: functioning with limited resources (memory, bandwidth, compute capacity, etc).
- Simplified device communication: easy communication between applications and devices in IoT and IIoT.

### Goals

The main goals of KubeEdge are as follows:

- Building an open edge computing platform with cloud native technologies.
- Helping users extending their business architecture, applications, services, etc. from cloud to edge but ensuring unified user experience.
- Implementing extensible architecture based on Kubernetes.
- Integrating with CNCF projects, including (but not limited to) containerd, cri-o, Prometheus, Envoy, etc.
- Seamlessly develop, deploy, and run complex workloads at the edge with optimized resources.

### Non-goals
KubeEdge does not take care of getting through the underlying physical network of edge nodes and ensuring the reliable communication of the cloud network infrastructure. However, when edge nodes are disconnected, KubeEdge supports offline autonomy for these nodes.

## Joint-review use

This document does not intend to provide a security audit of KubeEdge and is not intended to be used in lieu of a security audit. This document provides users of KubeEdge with a security focused understanding of KubeEdge and when taken with the self-assessment provide the community with the TAG-Security Review of the project. Both of these documents may be used and references as part of a security audit.

## Intended Use

### Target Users

KubeEdge is intended for:

- Platform Implementer: Based on the Kubernetes capabilities, KubeEdge extends cloud native to the edge, fully compatible with the Kubernetes ecosystem. KubeEdge can be built into types of devices, providing large-scale device management capabilities.
- Enterprise Operator: KubeEdge provides reports and unified metrics that can be consumed by a variety of dashboard platforms, allowing an enterprise operator to monitor and manage resources, for example, Prometheus is widely used in this scenario. This allows an Enterprise Operator to make administrative decisions and take action according to an organization’s policies.
- End User/Developer: KubeEdge supports routing management with the help of Kubernetes [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/#customresourcedefinitions) and a router module. Users can deliver their custom messages between cloud and edge.

### Use Cases

Please refer to [Case Studies](#case-studies).

### Operation

- [How to manage custom messages which are delivered between cloud and edge](https://kubeedge.io/en/docs/developer/custom_message_deliver/).
- [How to manage devices](https://kubeedge.io/en/docs/developer/device_crd/).

## Project Design

The following diagram shows the logical architecture for KubeEdge: 

<img src="assets/kubeedge_arch.png">

KubeEdge consists of below major components:

- CloudCore (in the cloud)
  - CloudHub: a websocket server responsible for watching changes at the cloud side, caching and sending messages to EdgeHub.
  - EdgeController: an extended Kubernetes controller which manages edge nodes and pods metadata so that the data can be targeted to a specific edge node.
  - DeviceController: an extended Kubernetes controller which manages devices so that the device metadata and devicetwin data can be synced between edge and cloud.
  - SyncController: an extended Kubernetes controller responsible for reliable data transmission between cloud and edge. It periodically checks meta data that persists on edge and cloud, and triggers reconcile if necessary.
  - DynamicController: an extended kubernetes controller based on Kubernetes dynamic client, which allows clients on the edge node to list/watch common Kubernetes resource and custom resources.
- EdgeCore (on the edge)
  - EdgeHub: a websocket client responsible for interacting with cloud services for edge computing (EdgeHub and CloudHub are symmetric components for edge-cloud communications). This includes syncing cloud resource updates to the edge, and reporting edge host and device status changes to the cloud.
  - MetaManager: the message processor between edged and edgehub. It is also responsible for storing/retrieving metadata to/from a lightweight database (SQLite).
  - Edged: an agent that runs on edge nodes and manages containerized applications.
  - DeviceTwin: responsible for storing devicetwin data and syncing devicetwin data to the cloud. It also provides query interfaces for applications.
  - EventBus: a MQTT client to interact with MQTT servers (mosquitto), offering pub-sub messaging capabilities to other components.
  - ServiceBus: a HTTP client to interact with HTTP servers (REST), offering HTTP client capabilities to components of cloud to reach HTTP servers running at edge.
  - MetaServer: MetaServer starts an HTTPS server with mutual TLS certs and acts as an edge api-server for Kubernetes operators. It proxies the Kubernetes resource request to the dynamic controller in the cloud.

### Functions and features

The features of KubeEdge include:

- **Kubernetes Native API at Edge**

  Autonomic Kube-API Endpoint at Edge, support to run third-party plugins and applications that depends on Kubernetes APIs on edge nodes.

- **Seamless Cloud-Edge Coordination**

  Bidirectional communication, able to talk to edge nodes located in private subnet. Supports both metadata and data.

- **Edge Autonomy**

  Metadata persistent per node, no list-watch needed during node recovery, and getting ready faster. Autonomous operation of edge even during disconnection from cloud.

- **Low Resource Readiness**

  Optimized usage of resources at the edge. Memory footprint down to about 70MB.

- **Simplified Device Communication**

  Easy communication between application and devices for IoT and Industrial Internet.

- **Heterogenous**

  Native support of x86, ARMv7, ARMv8. Heterogeneous access capability avoids a lot of adaptation work and enables customers to quickly access and deploy services. It can help extend more edge scenarios.

#### Security functions and features

As described in [KubeEdge Security Audit](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-security-audit-2022.pdf) under section `KubeEdge trust architecture`, we analyzed the security functions and features thoroughly. We isolate these parts into components because they form boundaries in the system where distinct trust relationships meet:

- CloudCore. CloudCore is connected to EdgeCore by way of EdgeHub.
- EdgeCore. The HTTP server in EdgeCore.
- MQTT broker part of EdgeCore that communicates with Devices. The MQTT broker part of EdgeCore accepts inputs from the Devices. Privileges flow from low to high in that an attacker in control of a device should not be able to cause adversarial affect on EdgeCore.
- Edged part EdgeCore which involves running of pods and, thus, containers.

## Configuration and Set-Up

### Defaults

KubeEdge initial the default configurations for components [cloudcore](https://github.com/kubeedge/kubeedge/blob/master/pkg/apis/componentconfig/cloudcore/v1alpha1/default.go) and [edgecore](https://github.com/kubeedge/kubeedge/blob/master/pkg/apis/componentconfig/edgecore/v1alpha1/default.go).

Please see the [configuration file](https://kubeedge.io/en/docs/setup/config/#modification-of-the-configuration-file) for more details about how to configure `cloudcore` and `edgecore`. 

### Advanced Security

As described in [KubeEdge threat model and security protection analysis paper](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-threat-model-and-security-protection-analysis.md#attack-surface-of-external-malicious-attackers), we have taken the following security hardening measures based on the default security principle. For security proposes, we have disabled these features by default, we also don't suggest enabling them without enforcing extra security reinforcement, or, unless they know the potential risks and impacts.

**On the CloudCore side:**

- The router HTTP server in CloudCore exposes CloudCore locally on the system. In this context, CloudCore should be protected against attempts from the local system to thwart KubeEdge. The router module is disabled by default.

**On the EdgeCore side:**

- The ServiceBus HTTP server in EdgeCore exposes EdgeCore locally on the system. In this context, EdgeCore should be protected against attempts from the local system to thwart KubeEdge. The ServiceBus module is disabled by default.
- The MetaServer starts a HTTP server and acts as an edge api-server for Kubernetes operators. The MetaServer module is disabled by default.

More recommendations about advanced security can be found at [security protection analysis paper](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-threat-model-and-security-protection-analysis.md#recommendations).

## Project Compliance

Not applicable

### Existing Audits

To conduct a more comprehensive security assessment of the KubeEdge project, the KubeEdge community cooperated with Ada Logics Ltd. and the Open Source Technology Improvement Fund (OSTIF) to perform a holistic security audit of KubeEdge and output a security auditing report in July, 2022. For more details please refer to the [KubeEdge security audit](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-security-audit-2022.pdf).

## Security Analysis

### Attacker Motivations

Attacks may come in the following ways:

- Attackers may try to gain access to the hosts which run the CloudCore component so that they can take control of KubeEdge control plane.
- Another attacker vector includes DOS’ing the CloudCore infrastructure, preventing edge nodes from connecting to the cloud. 
- Attackers may try to gain the access to the edge node which run EdgeCore, and exploit the edgecore.db which stores meta data on the edge node.

### Predisposing Conditions

According to [KubeEdge security audit in 2022](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-security-audit-2022.pdf) under section `KubeEdge trust architecture`, in order to understand the potential attack scope and the severity of a given attack we separate the system into components of trust. A component of trust specifically means a shared level of trust over data and a shared level of authority in terms of the system functioning properly. We then match these components with the potential actors and use this as a reference on classifying risk throughout the system.

We divide the system into the following components of trust: 

- CloudCore, more details in section [Project Design](#project-design).

- EdgeCore, more details in section [Project Design](#project-design).

- HTTP server of EdgeCore, including ServiceBus server.

- MQTT broker of EdgeCore that communicates with devices.

- Edged of EdgeCore which involves running of pods and, thus, containers.

We isolate these parts into components because they form boundaries in the system where distinct trust relationships meet: 

- CloudCore is connected to EdgeCore by way of EdgeHub. EdgeCore has a separate level of authority in the overall system, in that an attacker in control of EdgeCore should not be able to negatively affect other edge nodes, e.g. by way of manipulation of CloudCore. In this sense, the trust relationship flows from low to high in that EdgeCore has lower authority over KubeEdge than CloudCore, and CloudCore should be considered the highest level of trust in the KubeEdge ecosystem.
- The HTTP server in EdgeCore is exposed locally on the system. The HTTP server in EdgeCore exposes EdgeCore locally on the system and the applications on the local system do not have a higher level of authority over the full KubeEdge cluster. In this context, EdgeCore should be protected against attempts from the local system to thwart KubeEdge. In this sense, the trust relationship flows from low to high in that the local applications have lower authority over the KubeEdge ecosystem than EdgeCore.
- The MQTT broker part of EdgeCore accepts inputs from the Devices. Privileges flow from low to high in that an attacker in control of a device should not be able to cause adversarial affect on EdgeCore. The devices can reach EdgeCore via the MQTT broker. There is a trust boundary where the devices themselves should not be able to negatively affect the overall KubeEdge system. For example, if an attacker controls a device, the overall system should be protected against possible attacks from this device. In this sense, the trust relationship flows from low to high in that the devices have a lower authority over the KubeEdge system than EdgeCore.
- The edged of EdgeCore handles the running of Pods. edged handles the running of Pods. The containers in these pods have a low authority over the KubeEdge system and privileges thus flow from low to high. An attacker in control of some components in the containers should not be able to grow a foothold of more of the KubeEdge system.

**Notes: The MetaServer of component EdgeCore now starts a Secure HTTPS Server and provide mutual certificate authentication and RBAC authentication to ensure communication security. This security fix was hardened since we finished the security audit and has already been on the [master branch](https://github.com/kubeedge/kubeedge/tree/master) and will be available in the coming release v1.12.**

### Expected Attacker Capabilities

All internal components of KubeEdge use TLS to communicate with each other. There are currently no known vulnerabilities and operators have the ability to exploit TLS for the communications between internal components. As for the components (ServiceBus of EdgeCore) that communicate with external applications, TLS is not enforced, but the service is only exposed to the local host and the KubeEdge team has documented the risks of an insecure setup and made it harder for operators to deploy insecurely.

Areas that are potentially susceptible to attacks would require credential handling, for example, token-based single sign-on from edge to cloud. However, all of these use some forms of tokens based authentication and we assume the attacker cannot easily break existing cryptographic algorithms like AES or SHA256.

### Attack Risks and Effects

- For data saved in one database on the edge node, if an attacker attacks an edge database, it will cause denial of service in off line scenarios.
- External apps can request Kubernetes resources to MetaServer through list-watch, then the request will be passed to CloudCore and finally to the api-server of the cluster. If an app was exploited by an attacker, a DOS attack to list-watch server can be launched.

### Security Degradation

Exploiting vulnerabilities in management on the cloud can lead to maliciously deleting resources or bringing the whole cluster down. If an attacker intercepts tokens used for communicating with CloudCore, he can register fake nodes and render the CloudCore unusable. If attackers gain access to the EdgeCore database, he could manipulate critical metadata such as adding or deleting existing pods or change the host mode or image pull address of pods.

### Compensating Mechanisms

According to the trust boundary, we list 17 corresponding mitigations currently available in the KubeEdge project and 11 security reinforce recommendations for users and developers. Please see the details at [KubeEdge threat-model and mitigations](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-threat-model-and-security-protection-analysis.md#kubeedge-threat-model-and-mitigations).

## Threat Model

KubeEdge security team did an [overall system security analysis of KubeEdge](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-threat-model-and-security-protection-analysis.md#kubeEdge-threat-model-and-mitigations), mainly based on the [STRIDE threat modeling](https://en.wikipedia.org/wiki/STRIDE_(security)) and the [KubeEdge security audit](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-security-audit-2022.pdf). 

Please see the detailed threat model analysis paper of KubeEdge at section [KubeEdge threat-model and mitigations](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-threat-model-and-security-protection-analysis.md#kubeEdge-threat-model-and-mitigations).

### Identity Theft

Traffic between the trust boundary has the risk of identity theft, protecting your project from identity theft is a matter of treating all private keys and data as top secret information. More security protection analysis details are described at [Attack surface of external malicious attackers](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-threat-model-and-security-protection-analysis.md#attack-surface-of-external-malicious-attackers).

### Compromise

As a default, EdgeCore trusts all incoming data it receives via ServiceBus and the MQTT Broker. However, because these will often be deployed in untrusted environments and contexts, we include a section that goes further into detail regarding these two parts of EdgeCore.

The MQTT Broker and the ServiceBus receive data from sources that are often not manufactured or maintained by the KubeEdge cluster admin. Furthermore, the cluster admin will not have access to the source code of either the apps that connect to the ServiceBus or the software that runs on the devices that connect to the MQTT broker. Therefore, while EdgeCore trusts input to the ServiceBus and the MQTT Broker, the external services connecting to the ServiceBus and the MQTT Broker expose a critical attack surface that malicious actors will attempt to gain control over as a medium to get control over EdgeCore.

The detail implications of this attack surface are described at the [KubeEdge security audit](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-security-audit-2022.pdf) in section `EdgeCore: MQTT Broker and ServiceBus security`.

### Denial of Service

No incoming requests to the ServiceBus and MQTT Broker, even if these were controlled by a malicious actor, should cause denial of service to EdgeCore or CloudCore. More details are described at the [KubeEdge security audit](https://github.com/kubeedge/community/blob/master/sig-security/sig-security-audit/KubeEdge-security-audit-2022.pdf) in section `EdgeCore: MQTT Broker and ServiceBus security`.

## Secure Development Practices

KubeEdge has achieved the passing level criteria for [CII Best Practices](https://bestpractices.coreinfrastructure.org/en/projects/3018).

### Development Pipeline

All code is maintained in [GitHub](https://github.com/kubeedge/kubeedge) and changes must be reviewed by maintainers and must pass all Unit Tests, e2e tests, CI Fuzz, static checks, verifications on gofmt, go lint, go vet, vendors, and fossa checks. Code changes are submitted via Pull Requests (PRs) and must be signed. Commits to the `master` branch directly are not allowed.

### Communication Channels

- Internal. How do team members communicate with each other?
  

Team members communicate with each other frequently through [Slack Channel](https://join.slack.com/t/kubeedge/shared_invite/enQtNjc0MTg2NTg2MTk0LWJmOTBmOGRkZWNhMTVkNGU1ZjkwNDY4MTY4YTAwNDAyMjRkMjdlMjIzYmMxODY1NGZjYzc4MWM5YmIxZjU1ZDI), [KubeEdge sync meeting](https://zoom.us/my/kubeedge), and team members will open a new [issue](https://github.com/kubeedge/kubeedge/issues) to further discuss if necessary.

- Inbound. How do users or prospective users communicate with the team?
  

Users or prospective users usually communicate with the team through [Slack Channel](https://kubeedge.slack.com/archives/CDXVBS085), you can open a new [issue](https://github.com/kubeedge/kubeedge/issues) to get further help from the team, and [KubeEdge mailing list](https://groups.google.com/forum/#!forum/kubeedge) is also available. Besides, we have regular [community meeting](https://zoom.us/my/kubeedge) (includes SIG meetings) alternative between Europe friendly time and Pacific friendly time. All these meetings are publicly accessible and meeting records are uploaded to YouTube.

Regular Community Meetings:

  - Europe Time: **Wednesdays at 16:30-17:30 Beijing Time**. ([Convert to your timezone.](https://www.thetimezoneconverter.com/?t=16%3A30&tz=GMT%2B8&))
  - Pacific Time: **Wednesdays at 10:00-11:00 Beijing Time**. ([Convert to your timezone.](https://www.thetimezoneconverter.com/?t=10%3A00&tz=GMT%2B8&))

- Outbound. How do you communicate with your users? (e.g. flibble-announce@ mailing list)
  

KubeEdge communicates with users through [Slack Channel](https://kubeedge.slack.com/archives/CUABZBD55), [issues](https://github.com/kubeedge/kubeedge/issues), [KubeEdge sync meetings](https://zoom.us/my/kubeedge), [KubeEdge mailing list](https://groups.google.com/forum/#!forum/kubeedge). As for security issues, we provide the following channel:

- Security email group

  You can email to [kubeedge security team](mailto:cncf-kubeedge-security@lists.cncf.io) to report a vulnerability, and the team will disclose to the distributors through the [distributor announcement list](mailto:cncf-kubeedge-distrib-announce@lists.cncf.io) (more details [here](https://github.com/kubeedge/kubeedge/security/policy)).
  
### Ecosystem

KubeEdge helps users extend their business architecture, applications, services, etc. from cloud to edge while ensuring same user experience, implements extensible architecture based on Kubernetes and integrates with CNCF projects, including (but not limited to) containerd, cri-o, Prometheus, Envoy, etc. 

KubeEdge also integrates project KubeSphere to align with the cloud native ecosystem. KubeSphere is a distributed operating system for cloud-native application management, using Kubernetes as its kernel. It provides a plug-and-play architecture, allowing third-party applications to be seamlessly integrated into its ecosystem.

More KubeEdge adopters are listed in [ADOPTERS File](https://github.com/kubeedge/kubeedge/blob/master/ADOPTERS.md). More Vendors are listed in section [Related Projects / Vendors](#related-projects-/-vendors).

## Security Issue Resolution

### Responsible Disclosures Process

KubeEdge project vulnerability handling related processes are recorded in [Security Policy](https://github.com/kubeedge/kubeedge/security/policy). Related security vulnerabilities can be reported and communicated via email `cncf-kubeedge-security@lists.cncf.io`.

### Incident Response

See the [KubeEdge releases page](https://github.com/kubeedge/kubeedge/releases) for information on supported versions of KubeEdge. Once the fix is confirmed, the Security Team will patch the vulnerability in the next patch or minor release, and backport a patch release into the latest three minor releases.

The release of low to medium severity bug fixes will include the fix details in the patch release notes. Any public announcements sent for these fixes will be linked to the release notes.

### Closed security issues and vulnerabilities

According to the vulnerability release process, the KubeEdge team published the vulnerability in Security Advisories on Github. All the vulnerabilities have been assigned CVE IDs, focusing on the problems of memory exhaustion and access exceptions. For details, please see: https://github.com/kubeedge/kubeedge/security/advisories?state=published

## Hands-on review

The hands-on review is a lightweight review of the project's internal security as well as the current recommendation configuration, deployment, and interaction with regard to security.  Hands-on reviews are subject to security reviewer availability and expertise.  They are not intended to serve as an audit or formal assessment and are no gurantee of the actual security of the project.

**KubeEdge did not receive a hands-on review from TAG-Security.**

*If a hands-on review was performed, the below format should be used for reporting details*

|                    |                     |
| ------------------ | ------------------- |
| Date of review     | mmddyyyy-mmddyyyy   |
| Hands-on reviewers | name, github handle |

| Finding Number | Finding name | Finding Notes | Reviewer |
| -------------- | ------------ | ------------- | -------- |
|                |              |               |          |

### Hands-on review result

General comments and summary of the hands-on review with any recommendations worth
 noting.  If nothing found use the below example:

> TAG-Security's hands-on review did not reveal any significant or notable security findings for [project]. This outcome does not indicate that none exist, rather that none were discovered.

## Roadmap

* **Project Next Steps.** *Link to your general roadmap, if available, then list prioritized next steps that may have an impact on the risk profile of your project, including anything that was identified as part of this review.*
* **CNCF Requests.** *In the initial draft, please include whatever you believe the CNCF could assist with that would increase security of the ecosystem.*
  - Assess KubeEdge CI workflows and provide security hardening methods.
  - According to the threat model, provide a way to harden security from the overall perspective.

## Appendix

### Known Issues Over Time.

For details, please see: 

https://github.com/kubeedge/kubeedge/security/advisories?state=published

### Case Studies

KubeEdge has been widely adopted in industries including transportation, energy, Internet, CDN, manufacturing, smart campus etc.

* Tiansuan Constellation Program is a cloud native satellite computing platform initiated by the Shenzhen Institute of BUPT with Spacety Co., Ltd. KubeEdge brings in orbit-earth collaboration for the Tiansuan-1 Satellite in the platform, enabling new service scenarios such as image inference, incremental deep learning, and federated learning.
* In China’s highway electronic toll collection (ETC) system, KubeEdge helps manage nearly 100,000 edge nodes and more than 500,000 edge applications in 29 of China’s 34 provinces, municipalities, and autonomous regions. With these applications, the system processes more than 300 million data records daily and supports continuous update of ETC services on highways. Time used passing through toll stations is reduced from 29s to 3s for trucks, and from 15s to 2s for cars.
* As the the world’s longest sea crossing bridge, the [Hong Kong–Zhuhai–Macao bridge](https://en.wikipedia.org/wiki/Hong_Kong–Zhuhai–Macau_Bridge) is an open-sea fixed link that spans the Lingding and Jiuzhou channels, stretching for 55 km. KubeEdge helps manage the edge nodes deployed on the bridge that collect up to 14 different types of sensor data, including light intensity, carbon dioxide, atmospheric pressure, noise, temperature, humidity, PM 2.5 (fine particulate matter), PM 10 (particulate matter), rain and snow, acceleration, angular velocity, Euler angle, magnetic field, and sound. In addition, KubeEdge enables business apps, as well as device mapper and AI interference programs, to be deployed at the edge. See more details [here](https://www.altoros.com/blog/kubeedge-monitoring-edge-devices-at-the-worlds-longest-sea-bridge/).
* In the corporation with the Shanghai automotive industry, on the next-generation car-cloud collaborative architecture, SAIC MAXUS full-size luxury, smart, pure electric MPV MIFA 9 has become the first ever vehicle with KubeEdge Inside. Based on the capabilities provided by KubeEdge, such as light-weighted architecture, optimized for unstable cloud-edge network at scale, simplified heterogeneous IoT device management, edge autonomy in the case of disconnection, the platform can manage 100,000 + nodes and millions of edge devices. 200,000 new vehicles per year are now installed with KubeEdge.
* e-Cloud of China Telecom uses KubeEdge to manage CDN edge nodes, automatically deploy and upgrade CDN edge services, and implement edge service disaster recovery (DR) when they migrates their CDN services to the cloud. See more details [here](https://www.cncf.io/blog/2022/03/18/e-cloud-large-scale-cdn-using-kubeedge).
* China Mobile On-line Marketing Service Center, a secondary organ of the China Mobile Communications Group, which holds the world’s largest call center with 44,000 agents, 900 million users, and 53 customer service centers, builds a cloud-edge synergy architecture consisting of two centers and multiple edges based on KubeEdge. See more details [here](https://www.cncf.io/blog/2021/08/16/china-mobile-kubeedge-based-customer-service-platform-featuring-edge-cloud-synergy).
* Xinghai IoT is an IoT company that provides comprehensive smart building solutions by leveraging a construction IoT platform, intelligent hardware, and AI. Xinghai IoT built a smart campus with cloud-edge-device synergy based on KubeEdge and its own Xinghai IoT cloud platform, greatly improving the efficiency of campus management. With AI assistance, nearly 30% of the repetitive work is automated. In the future, Xinghai IoT will continue to collaborate with KubeEdge to launch KubeEdge-based smart campus solutions.
* [SF Express](https://www.sf-express.com/), the largest integrated logistics service provider in China and the fourth largest in the world, has completely transformed their supply chains. Based on KubeEdge, SF Express smart supply chain system realizes edge-cloud synergy, edge data collection, real-time edge data processing, massive IoT device management, and creates a more stable, secure, and reliable edge industrial IoT system.

### Related Projects / Vendors

#### Related Projects 

As the CNCF's first Cloud Native Edge Computing incubating project (announced in 2020), KubeEdge continues to focus on extending cloud native capabilities to the edge. At the same time, KubeEdge can also be integrated with other CNCF projects to enrich cloud native ecosystem.

As there are many related cloud native edge projects, this article will narrow down to the following CNCF edge projects. **Akri** is an open source project that exposes leaf devices as resources in a Kubernetes cluster. It leverages and extends the Kubernetes [device plugin framework](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/), which was created with the cloud in mind and focuses on advertising static resources such as GPUs and other system hardware. Akri took this framework and applied it to the edge, where there is a diverse set of leaf devices with unique communication protocols and intermittent availability.

**K3s** is a CNCF sandbox project that delivers a lightweight certified Kubernetes distribution and can be run at the edge. A K3s user can manipulate Kubernetes resources by calling the K3s API on the server node. A server node is defined as a machine (bare-metal or virtual) running the `k3s server` command. A worker node is defined as a machine running the `k3s agent` command. Agent nodes are registered with a websocket connection initiated by the `k3s agent` process, and the connection is maintained by a client-side load balancer running as part of the agent process.

**OpenYurt** and **SuperEdge** also provide cloud native capabilities in the edge computing area and overlap with the KubeEdge project.

On the edge side, KubeEdge can manage kinds of applications by adapting to runtimes listed as follows.

As KubeEdge adapts to the Container Runtime Interface (CRI), KubeEdge can be integrated easily with runtimes like **Containerd**, **CRI-O**, etc. In the FaaS (Function as a Service) area, **OpenFunction** can be integrated to provide capabilities that let you focus on your business logic without having to maintain the underlying runtime environment and infrastructure. You can generate event-driven and dynamically scaling serverless workloads by simply submitting business-related source code in the form of functions.

As WebAssembly is increasingly used in Edge Computing scenarios where it is difficult to deploy Linux containers or when the application performance is vital, KubeEdge provides a lightweight alternative to Linux containers in edge native environments by integrating with **WasmEdge**. WasmEdge is a lightweight, high-performance, and extensible WebAssembly runtime for cloud native, edge, and decentralized applications. It powers serverless apps, embedded functions, microservices, smart contracts, and IoT devices.

For the images registry infrastructure, KubeEdge can work with **Harbor** to manage images efficiently. Harbor is an open source trusted cloud native registry project that stores, signs, and scans content. Harbor extends the open source Docker Distribution by adding the functionalities usually required by users such as security, identity and management. Having a registry closer to the build and run environment can improve the image transfer efficiency.

#### Vendors and products

The KubeEdge project is also working with many other vendors to integrate KubeEdge into their solutions to enrich the Cloud Native Edge Computing ecosystem.

EdgeStack, an intelligent edge computing platform of HarmonyCloud, builds a edge-cloud synergy system through the integration with KubeEdge, which can support access of millions of edge nodes and devices. It is designed for large-scale, massive access, low bandwidth, low latency, high performance, and high stability of edge computing.  At present, it has landed in communications, transportation, finance, and many other fields.

The DaoCloud Edge Computing platform uses KubeEdge to manage edge applications. With cloud-edge synergy, this platform responds to edge application requests in real time, and constantly monitors node, device, and application status. The cloud is responsible for edge node registration, management, and application and configuration delivery. At the edge runs edge applications to which edge autonomy is available. For devices, multi-protocol access is supported and standard APIs are provided to connect to devices.

Based on KubeEdge, Intelligent EdgeFabric (IEF) of Huawei Cloud embeds cloud native into edge computing, supporting ultimate lightweight deployments, edge intelligence, and powerful computing power. It has been widely used in smart campus, industrial quality inspection, mining, and smart transportation where the collaboration between edge and cloud yields brilliant results. 

KubeSphere Enterprise (KSE), is an enterprise container management platform from KubeSphere open source community. Based on KubeEdge, applications and workloads are uniformly distributed and managed on the cloud and edge nodes, meeting the needs for application delivery, O&M, and management on a large number of edge and device devices.

EMQ is a software provider of open-source IoT data infrastructure. As the core product of EMQ, EMQX is a reliable open-source MQTT messaging platform, supports 100M concurrent IoT device connections per cluster while maintaining 1M message per second throughput and sub-millisecond latency, uses KubeEdge in edge side to manage middle-ware deployments.

As a distributed cloud, Inspur Cloud can provide users with cloud products and service capabilities, and extend cloud native capabilities to the edge based on KubeEdge.

Click2Cloud is a provider of cloud-based products and services. It offers managed cloud to optimize, manage, and enhance cloud efficiency. In the edge computing area, KubeEdge helps Click2Cloud provide users with more edge computing solutions.

Inovex is able to better fulfill the vision of transforming the potentials provided by digitization into excellent solutions by integrating KubeEdge.
