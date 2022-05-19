kubectl delete -f ../deployment.yaml
kubectl delete -f ../autoscaler.yaml

kubectl delete -f highload-deployment.yaml
kubectl delete -f highload-autoscaler.yaml

sleep 10

kubectl apply -f ../deployment.yaml
kubectl apply -f ../autoscaler.yaml
kubectl apply -f highload-deployment.yaml
kubectl apply -f highload-autoscaler.yaml
