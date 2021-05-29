# vcreport

[![build](https://github.com/invit/vcreport/actions/workflows/build.yml/badge.svg)](https://github.com/invit/vcreport/actions/workflows/build.yml)

CLI to display [version-checker](https://github.com/jetstack/version-checker) metrics in a human-readable way. 

## Installation

Downloadable binaries are available from the [releases page](https://github.com/invit/vcreport/releases/latest).

## Usage

```
Usage:
  vcreport metrics-url [flags]

Flags:
  -a, --all       Show all images, not just outdated ones
  -b, --brief     Just show images, but no pods
  -h, --help      help for vcreport
  -v, --version   version for vcreport
```

_metrics-url_ is the full URL to the metrics endpoint of version-checker.

### Examples

* Run locally

```shell
$ kubectl port-forward service/version-checker 8080:8080 --namespace=version-checker &
$ vcreport -a -b http://localhost:8080/metrics
+----------------------------------------------------------+----------------------+
|                          IMAGE                           |       VERSION        |
+----------------------------------------------------------+----------------------+
| some-vendor/image                                        | 2.0.0 > 2.1.0        |
+----------------------------------------------------------+----------------------+
| quay.io/jetstack/version-checker                         | v0.2.1 (Up to date)  |
+----------------------------------------------------------+----------------------+
| redis                                                    | 6.2.2-alpine > 6.2.3 |
+----------------------------------------------------------+----------------------+
| us.gcr.io/k8s-artifacts-prod/autoscaling/vpa-recommender | 0.9.2 (Up to date)   |
+----------------------------------------------------------+----------------------+
```

* Run in Kubernetes cluster

```shell
$ kubectl run --namespace=version-checker -i --tty --rm vcreport --image=ghcr.io/invit/vcreport/vcreport:latest --restart=Never -- http://version-checker:8080/metrics
+----------------------------------------------------------+----------------------+--------------------------------------+
|                          IMAGE                           |       VERSION        |                 PODS                 |
+----------------------------------------------------------+----------------------+--------------------------------------+
| some-vendor/image                                        | 2.0.0 > 2.1.0        | namespace/pod-1/container            |
|                                                          |                      | namespace/pod-1/container            |
+----------------------------------------------------------+----------------------+--------------------------------------+
| redis                                                    | 6.2.2-alpine > 6.2.3 | some-other-namespace/pod-1/container |
+----------------------------------------------------------+----------------------+--------------------------------------+
pod "vcreport" deleted
```

## Build

On Linux:

```
$ git clone github.com/invit/vcreport 
$ cd vcreport
$ make 
```

## License

vcreport is licensed under the [MIT License](http://opensource.org/licenses/MIT).
