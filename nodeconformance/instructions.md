# How to submit node conformance results

## About the KubeEdge node conformance tests

The standard set of node conformance tests of KubeEdge is currently those defined by the `[sig-node] * [Conformance]` tag in the [KubeEdge e2e](https://github.com/kubeedge/kubeedge/tree/master/tests/e2e) suite, plus a part of kubernetes node conformance tests included at [here](https://github.com/kubeedge/kubeedge/blob/master/build/conformance/kubernetes/kube_node_conformance_test.go).

## Running

KubeEdge has provided the docker image of the node conformance test, which contains the scripts and related files of the node conformance test. Follow these steps to perform a node conformance test.

### Prerequisite

- At least one master node and one edge node exist on different VMS or physical machines.
- Enable the CloudStream and EdgeStream modules.

### Launch the conformance test container

Pull node conformance test image `kubeedge/nodeconformance` 

```
$ docker pull kubeedge/nodeconformance
```

Or build the image locally.

```
$ git clone https://github.com/kubeedge/kubeedge.git
$ cd kubeedge
$ docker build -t {image_name}:{tag_name} -f build/conformance/nodeconformance.Dockerfile .
```

Examples of running the conformance test containers:

```
docker run --env KUBECONFIG=/root/.kube/config  --env RESULTS_DIR=/tmp/results -v /root/.kube/config:/root/.kube/config -v /tmp/results:/tmp/results --network host -it kubeedge/nodeconformance
```

Description of container environment variables:

| Environment variables | The corresponding ginkgo parameters | Parameters description                                                                                                          |
|-----------------------|-------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|
| E2E_DRYRUN            | dryRun                              | If set, ginkgo will walk the test hierarchy without actually running anything. Best paired with -v.                             |
| E2E_SKIP              | skip                                | If set, ginkgo will only run specs that do not match this regular expression. Can be specified multiple times, values are ORed. |
| E2E_FOCUS             | focus                               | If set, ginkgo will only run specs that match this regular expression. Can be specified multiple times, values are ORed.        |
| RESULTS_DIR           | NA                                  | Output report path, default /tmp/results                                                                                        |
| REPORT_PREFIX         | NA                                  | Report file prefix, optional                                                                                                    |
| IMAGE_URL             | NA                                  | The name of the image used by the test case                                                                                     |
| TEST_WITH_DEVICE      | NA                                  | Whether to test device. The default value is false                                                                              |
| GINKGO_BIN            | NA                                  | `ginkgo` binary path，default /usr/local/bin/ginkgo                                                                              |
| E2E_EXTRA_ARGS        | NA                                  | Extra arguments                                                                                                                 |
| KUBECONFIG            | NA                                  | `kubeconfig` file path                                                                                                          |
| TEST_BIN              | NA                                  | `e2e.test` binary path，default /usr/local/bin/e2e.test                                                                          |

## Uploading

Prepare a PR to [https://github.com/kubeedge/community](https://github.com/kubeedge/community) to upload the report files to the directory `nodeconformance`. In the descriptions below, `X.Y` refers to the KubeEdge major and minor version, and `$dir` is a short subdirectory name to hold the results for your product.

Description: `Node Conformance results for vX.Y/$dir`

### Contents of the PR

For simplicity you can submit the tarball or extract the relevant information from the tarball to compose your submission.

If submitting test results for multiple versions, submit a PR for each product, i.e. one PR for vX.Y results and a second PR for vX.Z

```
vX.Y/$dir/README.md: A script or human-readable description of how to reproduce your results.
vX.Y/$dir/e2e.log: Test log output (from the container kubeedge/nodeconformance).
vX.Y/$dir/junit_conformance.xml: Machine-readable test log (from the container kubeedge/nodeconformance).
vX.Y/$dir/PRODUCT.yaml: See below.
```

#### PRODUCT.yaml

This file describes your product. It is YAML formatted with the following root-level fields. Please fill in as appropriate.

| Field                   | Description                                                                                                                                                                                                                               |
|-------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `vendor`                | Name of the legal entity that is certifying.                                                                                                                                                                                              |
| `name`                  | Name of the product being certified.                                                                                                                                                                                                      |
| `version`               | The version of the product being certified (not the version of KubeEdge it runs).                                                                                                                                                         |
| `website_url`           | URL to the product information website                                                                                                                                                                                                    |
| `repo_url`              | If your product is open source, this field is necessary to point to the primary GitHub repo containing the source. It's OK if this is a mirror. OPTIONAL                                                                                  |
| `documentation_url`     | URL to the product documentation                                                                                                                                                                                                          |
| `product_logo_url`      | URL to the product's logo, (must be in SVG, AI or EPS format -- not a PNG -- and include the product name). OPTIONAL. If not supplied, we'll use your company logo. Please see logo [guidelines](https://github.com/cncf/landscape#logos) |
| `type`                  | Is your product a distribution, hosted platform, or installer.                                                                                                                                                                            |
| `description`           | One sentence description of your offering                                                                                                                                                                                                 |
| `contact_email_address` | An email address which can be used to contact maintainers regarding the product submitted and updates to the submission process                                                                                                           |

Examples:

```
vendor: Huawei Cloud
name: IEF
version: v1.x.x
website_url: https://xxx
repo_url: https://xxx
documentation_url: https://xxx/docs
product_logo_url: https://xxx.svg
type: hosted platform
description: "Based on the KubeEdge and kubernetes ecosystems, IEF applies cloud native technologies to edge computing."
```

## Contact

If you have any problems about certifying or conformance testing, please file an issue in the [kubeedge/community](https://github.com/kubeedge/community). Questions and comments can also be sent to the working group's [slack channel](https://kubeedge.slack.com/archives/CKVUCM5ED). SIG Testing is the change controller of the conformance definition.