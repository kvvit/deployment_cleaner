# Kubernetes Deployments Cleaner

This is my first step in the big and fascinating world of Go programming.
The main idea of this application is to clean up Kubernetes deployments in the
development environment. It helps me delete outdated deployments automatically.
Because I am a big fan of Kubernetes, I decided to make this application for
my learning purposes of Go programming.

## Use case

This application can be used in the development kubernetes cluster, to
automatically delete outdated deployments in particular namespace.
In current version I commented the parts of the code that related to
deleting actions. This parts located in the module `pkg/deleteobjects`.

## Installation

This application can be installed via helm chart.:

```bash
    helm upgrade --install deployments-cleaner ./charts \
    -f ./charts/values.yaml -n your-namespace
```
