# Kubernetes Deployments Cleaner

This is my first step in the big and fascinating world of Go programming.
The main idea of this application is to clean up Kubernetes deployments in the
development environment. It helps me delete outdated deployments automatically.
Because I am a big fan of Kubernetes, I decided to make this application for
my learning purposes of Go programming.

## Use case

This application can be used in the development kubernetes cluster, to
automatically delete outdated deployments in particular namespace.

## Installation

This application can be installed via helm chart.:

```bash
    helm upgrade --install deployments-cleaner ./charts \
    -f ./charts/values.yaml -n your-namespace
```

Environmet variables used in this application, with default values that
can be changed in `values.yaml`:

```yaml
WORK_START: 10
```

The time in which deploymets cleaner will stop to work, or begining of workday.

```yaml
WORK_END: 19
```

The time in which deploymets cleaner will start to work, or end of workday.

```yaml
TIME_TO_DELETE: 86400
```

Time in seconds that deploymets cleaner will wait before deleting outdated
deployments. Here is 86400 seconds that means 24 hours.

```yaml
DRY_RUN: true
```

Skip deleting actions if it set to `true`.

## Monitoring

This application can be monitored via Grafana. It returns deployments
names and timelive of each deployment in seconds in port `8080` in the
`/metrics` path.
