#create cluster
eksctl create cluster -f clusterconfig.yaml

#create a service account, a cluster admin and eks admin
kubectl apply -f eks-admin-service-account.yaml

#install metrics server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
