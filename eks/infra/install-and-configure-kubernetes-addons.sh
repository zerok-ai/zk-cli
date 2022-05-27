#--------- Metrics server ---------
echo '---------------------- Installing metrics-server'
kubectl apply -f ./yaml/cluster/metrics-server.yaml


#--------- Kubernetes dashboard ---------

#install kubernetes dashboard
# echo '---------------------- Installing dashboard'
# kubectl apply -f ./yaml/cluster/dashboard.yaml

#start commandline proxy in background
echo '---------------------- Starting proxy'
kill -9 $(ps aux | grep "kubectl proxy" | grep -v "CVS" | awk '{print $2}')
kubectl proxy &

echo "Secret copied to clipboard. Open the following link to access the dashboard"
echo "http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/"