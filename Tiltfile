load('ext://helm_resource', 'helm_resource', 'helm_repo')

def grafana_compose(labels=["grafana"]):
    tfdir = os.path.dirname(__file__)
    docker_compose(os.path.join(tfdir, 'compose/docker-compose.yaml'))
    dc_resource('grafana', labels=labels)
    dc_resource('loki', labels=labels)
    dc_resource('tempo', labels=labels)
    dc_resource('promtail', labels=labels)

    logfile = tfdir+ "/compose/logs/tilt.log"
    local_resource('log-forwarder', serve_cmd="tilt logs -f | sed 's/│/\\|/g' > " + logfile, labels=labels)
    return struct(
        otlp_grpc = "localhost:4317",
        otlp_http = "localhost:4318",
        zipkin = "localhost:9411",
        jaeger_grpc = "localhost:14268" #?
    )

def grafana_kubernetes(namespace="default", labels=["grafana"]):
    tfdir = os.path.dirname(__file__)

    helm_repo('grafana-helm', 'https://grafana.github.io/helm-charts')
    helm_resource('loki', 'grafana/loki-stack')
    helm_resource('grafana', 'grafana/grafana', flags=["-f", os.path.join(tfdir, 'grafana-values.yaml')])
    helm_resource('tempo', 'grafana/tempo', flags=["-f", os.path.join(tfdir, 'tempo-values.yaml')])

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

    return struct(
        otlp_grpc = "tempo:4317",
        otlp_http = "tempo:4318",
        zipkin = "tempo:9411",
        jaeger_grpc = "tempo:14250"
    )