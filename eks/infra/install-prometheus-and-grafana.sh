# create namespace for monitoring stack
kubectl create namespace monitoring

# add helm repo for prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update 

# install kube-prometheus
helm install prom prometheus-community/kube-prometheus-stack -n monitoring --values ./yaml/values/prometheus-grafana.yaml