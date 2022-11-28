# Tilt Grafana

A [Tilt Extension](https://docs.tilt.dev/extensions.html) for working with the LGTM stack from Grafana (Loki, Grafana, Tempo, Mimir).

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

## Docker Compose

```
load('ext://tilt-grafana', 'grafana_compose')
endpoints = grafana_compose()

...
```

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