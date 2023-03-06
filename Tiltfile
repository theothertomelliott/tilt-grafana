load('ext://helm_resource', 'helm_resource', 'helm_repo')
load('ext://configmap', 'configmap_create')
load("compose/prometheus/Tiltfile", "prometheus_compose_impl")

def grafana_compose(
    labels=["grafana"], 
    metrics_endpoints=[],
    dashboard_files=[],
    ):

    tfdir = os.path.dirname(__file__)
    docker_compose(os.path.join(tfdir, 'compose/docker-compose.yaml'))
    dc_resource('grafana', labels=labels, resource_deps=["Copy dashboards"])
    dc_resource('loki', labels=labels)
    dc_resource('tempo', labels=labels)
    dc_resource('promtail', labels=labels)
    dc_resource('mimir', labels=labels)

    # Set up the dashboards directory with the specified files
    dashboard_command = "rm -rf " + tfdir+"/compose/dashboards/dashboards/*.json"
    for index, dashboard in enumerate(dashboard_files):
        dashboard_command = dashboard_command + " && cp " + dashboard + " " + tfdir+"/compose/dashboards/dashboards/dashboard" + str(index) + ".json"
    local_resource(
            "Copy dashboards",
            cmd=dashboard_command,
            )

    logfile = tfdir+ "/compose/logs/tilt.log"
    local_resource('log-forwarder', serve_cmd="tilt logs -f | sed 's/│/\\|/g' > \"" + logfile + "\"", labels=labels)

    if len(metrics_endpoints) > 0:
        prometheus_compose_impl(endpoints=metrics_endpoints,labels=labels)

    return struct(
        otlp_grpc = "localhost:4317",
        otlp_http = "localhost:4318",
        zipkin = "localhost:9411",
        jaeger_http = "localhost:14268",
        jaeger_grpc = "localhost:14250",
        mimir = "http://localhost:9009/api/v1/push"
    )

def metrics_endpoint(name, port, path="/metrics"):
    return {"name": name, "port": port, "path": path}

def grafana_kubernetes(
    namespace="default", 
    labels=["grafana"],
    mimir_enabled=False,
    dashboard_files=[],
    ):

    tfdir = os.path.dirname(__file__)

    # Set up the dashboards with the specified files
    from_file = ['README.md=' + tfdir + '/README.md'] # Include a non-json file to allow for calling without dashboard files
    for index, dashboard in enumerate(dashboard_files):
        from_file.append(
            "dashboard" + str(index) + ".json" + "=" + dashboard
        )
    configmap_create(
        'grafana-dashboard-config', 
        namespace=namespace, 
        from_file=from_file,
    )

    helm_repo('grafana-helm', 'https://grafana.github.io/helm-charts', labels=labels)
    
    helm_resource('loki', 'grafana-helm/loki-stack', resource_deps=["grafana-helm"])
    helm_resource(
        'grafana', 
        'grafana-helm/grafana', 
        flags=[
            "-f", os.path.join(tfdir, 'kubernetes/grafana-values.yaml'),
            ], 
        resource_deps=["grafana-helm"],
    )
    helm_resource('tempo', 'grafana-helm/tempo', flags=["-f", os.path.join(tfdir, 'kubernetes/tempo-values.yaml')], resource_deps=["grafana-helm"])
    helm_resource('phlare', 'grafana-helm/phlare', resource_deps=["grafana-helm"])

    helm_repo('prometheus-community','https://prometheus-community.github.io/helm-charts', labels=labels)

    helm_resource('prometheus', 'prometheus-community/prometheus', flags=["-f", os.path.join(tfdir, 'kubernetes/prometheus-values.yaml')], resource_deps=["prometheus-community"])

    if mimir_enabled:
        k8s_yaml([
            os.path.join(tfdir,"kubernetes/mimir/deployment.yaml"),
            os.path.join(tfdir,"kubernetes/mimir/service.yaml"),
            os.path.join(tfdir,"kubernetes/mimir/mimir-config.yaml")
            ])

    k8s_resource(
        "grafana", 
        port_forwards="3000:3000",
        labels=labels
    )
    k8s_resource(
        "loki", 
        labels=labels
    )
    k8s_resource(
        "tempo",
        labels=labels
    )
    k8s_resource(
        "phlare",
        labels=labels
    )
    k8s_resource(
        "prometheus",
        port_forwards="9090:9090",
        labels=labels
    )
    if mimir_enabled:
        k8s_resource(
            "mimir",
            labels=labels
        )

    return struct(
        otlp_grpc = "tempo.default:4317",
        otlp_http = "tempo.default:4318",
        zipkin = "tempo.default:9411",
        jaeger_http = "tempo.default:14268",
        jaeger_grpc = "tempo.default:14250"
    )