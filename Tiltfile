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
        port_forwards="4318:4318",
        labels=labels
    )