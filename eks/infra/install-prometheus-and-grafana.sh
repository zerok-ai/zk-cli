
# install Prometheus 
kubectl create namespace prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/prometheus --namespace prometheus


# install grafana
#-------------------------------------------
kubectl create namespace grafana
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana grafana/grafana  \
	–namespace grafana   \
	–set persistence.enabled=true  \
	–set adminPassword='1rmmIM7hxTa0' \
	–set datasources.'datasources\.yaml'.apiVersion=1   \
	–set datasources.'datasources\.yaml'.datasources[0].name=Prometheus  \
	–set datasources.'datasources\.yaml'.datasources[0].type=prometheus    \
	–set datasources.'datasources\.yaml'.datasources[0].url=http://prometheus-server.prometheus.svc.cluster.local     \
	–set datasources.'datasources\.yaml'.datasources[0].access=proxy     \
	–set datasources.'datasources\.yaml'.datasources[0].isDefault=true 

# validate
echo 'You can validate the cluster by running -- `kubectl get pods -n prometheus`'
echo 'access the Prometheus dashboard through -- `http:\\localhost:9090`'
echo 'and Grafana dashboard through -- `http:\\localhost:3000`'