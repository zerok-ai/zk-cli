#! /bin/bash

kubectl apply -f ../ingress.yaml

if [[ $1 == "full" ]]; then
	kubectl delete -f highload-deployment.yaml
	kubectl delete -f highload-services.yaml
	kubectl delete -f highload-autoscaler.yaml
fi
