# This will install prometheus and grafana in monitoring namespace
# install Prometheus 
kubectl create -f ./manifests/setup/
kubectl create -f ./manifests/

# validate
echo 'You can validate the cluster by running -- `kubectl get pods -n prometheus`'
echo 'access the Prometheus dashboard through -- `http:\\localhost:9090`'
echo 'and Grafana dashboard through -- `http:\\localhost:3000`'