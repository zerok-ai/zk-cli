# Cluster

### Installation
```
sh install.sh
```
### Uninstall
```
sh uninstall.sh
```

# Monitoring

### Kubernetes dashboard
#####Start proxy
```
kubectl proxy
```

If you get the error `error: listen tcp 127.0.0.1:8001: bind: address already in use` run the following command and try to start the proxy again

```
kill -9 $(ps aux | grep "kubectl proxy" | grep -v "CVS" | awk '{print $2}')
```
#####Dashboard url
```
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/
```
#####Dashboard token
The following command will copy the token to clipboard
```
kubectl -n kube-system describe secret $(kubectl
-n kube-system get secret |  awk '/^deployment-controller-token-/{print $1}') | grep 'token:' | awk '{print $2}' | pbcopy
```
### Install Prometheus-Grafana
Run 
```
sh install-prometheus-and-grafana.sh
```

### Prometheus-Grafana

#####Port-forward Prometheus
```
kubectl port-forward service/prom-kube-prometheus-stack-prometheus -n monitoring 9090:9090
```

`http://localhost:9090/`

#####Port-forward Grafana
```
kubectl port-forward service/prom-grafana -n monitoring 3000:8080
```

`http://localhost:3000/`

Login: admin
Password: prom-operator


