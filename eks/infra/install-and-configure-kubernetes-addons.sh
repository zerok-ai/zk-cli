#--------- Metrics server ---------
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml


#--------- Kubernetes dashboard ---------

#install kubernetes dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.5.0/aio/deploy/recommended.yaml

#start commandline proxy in background
kill -9 $(ps aux | grep "kubectl proxy" | grep -v "CVS" | awk '{print $2}')
kubectl proxy &

echo "Secret copied to clipboard. Open the following link to access the dashboard"
echo "http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/"