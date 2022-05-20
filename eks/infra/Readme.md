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

### Installation
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


