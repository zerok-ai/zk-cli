#----------- Install NGINX proxy based ingress ----------- 

# Add the Helm chart for Nginx Ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

# Create namespace
kubectl create namespace ingress

# Install the Helm (v3) chart for nginx ingress controller
helm upgrade --install app-ingress ingress-nginx/ingress-nginx \
	--namespace ingress \
	--values ./yaml/values/nginx-ingress-values.yaml

# Print the Ingress Controller public IP address
kubectl get services --namespace ingress

#----------- 

# helm upgrade prometheus prometheus-community/kube-prometheus-stack \
# --namespace prometheus  \
# --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
# --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false