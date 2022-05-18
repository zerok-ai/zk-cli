#! /bin/bash

ELB_HOSTNAME=$(kubectl get services --namespace ingress | grep 'app-ingress-ingress' | grep -v 'admission' | awk '{print $4}')
sed 's/<<ELB_HOSTNAME>>/'"$ELB_HOSTNAME"'/g' app-ingress-template.yaml > app-ingress.yaml

kubectl apply -f node-deployment.yaml
kubectl apply -f services.yaml
kubectl apply -f app-ingress.yaml

echo
echo "Application deployed at:"
echo "http://"$ELB_HOSTNAME"/info1"
echo
echo "Fetching Pod info: "
kubectl get pods --namespace app --show-labels

echo
echo "More commands: "
echo "  1. Label all the pods of the app"
echo "  kubectl label pods --namespace app --selector=app=load-test load=high"
echo
echo "  2. Set label to a pod"
echo "  kubectl label pods --namespace app <name of the pod> load=low --overwrite=true"
echo
