#! /bin/bash

if [[ $1 == "full" ]]; then
	kubectl apply -f highload-deployment.yaml
	kubectl apply -f highload-services.yaml
	# kubectl apply -f highload-autoscaler.yaml
	kubectl apply -f highload-service-monitor.yaml
fi

kubectl apply -f highload-ingress.yaml
