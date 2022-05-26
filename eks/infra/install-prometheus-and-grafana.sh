# add helm repo for prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update 

# create namespace for monitoring stack
kubectl create namespace monitoring

# install kube-prometheus
helm upgrade --install prom prometheus-community/kube-prometheus-stack \
	--namespace monitoring \
	--values ./yaml/values/prometheus-grafana.yaml

# install Promtail
# helm upgrade --install promtail grafana/promtail -f ./yaml/values/promtail-values.yaml -n monitoring

# install Loki 
# helm upgrade --install loki grafana/loki-distributed -n monitoring
