# Tilt Grafana

A [Tilt Extension](https://docs.tilt.dev/extensions.html) for working with the LGTM stack from Grafana (Loki, Grafana, Tempo, Mimir),

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
grafana_kubernetes()
```

## Docker Compose

```
load('ext://tilt-grafana', 'grafana_compose')
grafana_compose()
```