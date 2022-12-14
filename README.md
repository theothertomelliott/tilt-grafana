# Tilt Grafana

A [Tilt Extension](https://docs.tilt.dev/extensions.html) for working with the LGTM stack from Grafana ([Loki](https://grafana.com/oss/loki/), [Grafana](https://grafana.com/oss/grafana/), [Tempo](https://grafana.com/oss/tempo/), [Mimir](https://grafana.com/oss/mimir/)). [Prometheus](https://prometheus.io/) is included to scrape metrics for Mimir. Finally, as a bonus, [Phlare](https://grafana.com/docs/phlare/latest/) is also installed when running under Kubernetes.

# Resource Requirements

Running the full LGTM stack locally can use almost 3GB of RAM in a steady state. In order to provide your other workloads with enough resources to run, it is recommended that you provide your local Kubernetes cluster with at least 8GB of RAM.

You will typically set the total memory available to your cluster in your Docker configuration, with a few
examples below:

* [Docker Desktop](https://docs.docker.com/desktop/settings/mac/#resources)
* [Colima](https://github.com/abiosoft/colima#customization-examples)

# Usage

To import this repo for use, add the following to your `Tiltfile`:

```
v1alpha1.extension_repo(name='tilt-grafana', url='http://github.com/theothertomelliott/tilt-grafana')
v1alpha1.extension(name='tilt-grafana', repo_name='tilt-grafana', repo_path='')
```

Then you can import the appropriate functions either for Kubernetes or Docker Compose.

## Kubernetes

```
load('ext://tilt-grafana', 'grafana_kubernetes')
endpoints = grafana_kubernetes()

...
```

Mimir can use a lot of memory, so it is disabled in Kubernetes by default. You can enable it by setting
`mimir_enabled=True`:

```
load('ext://tilt-grafana', 'grafana_kubernetes')
endpoints = grafana_kubernetes(mimir_enabled=True)
```

## Docker Compose

```
load('ext://tilt-grafana', 'grafana_compose')
endpoints = grafana_compose()

...
```

You may also specify a set of local endpoints to be scraped by Prometheus with the `metrics_endpoints` parameter:

```
endpoints = grafana_compose(
    metrics_endpoints=[
        metrics_endpoint(name='jobname', port=1234)
    ]
)
```

The above example will result in Prometheus scraping metrics at http://localhost:1234/metrics. These
metrics will then be forwarded to Mimir.

If the `metrics_endpoints` parameter is omitted, no Prometheus instance will be created.

# Endpoints

The `endpoints` struct is returned by both `grafana_kubernetes` and `grafana_compose` and contains details
of the ingestion endpoints that your code may need to be able to send telemetry to the Grafana stack.

Fields within this struct are currently:

* otlp_grpc
* otlp_http
* zipki
* jaeger_grpc

For example, the test compose-based Tiltfile at `test/compose` uses the `otlp_http` value as an environment
variable for its synthetic load generator program:

```
endpoints = grafana_compose()

local_resource(
    "generator",
    serve_cmd="go run ../generator",
    serve_env={
        # Pass the otlp_http endpoint value to the generator so it can send traces.
        "OTEL_OTLP_HTTP_ENDPOINT": endpoints.otlp_http,
    },
    labels=["generator"]
)
```