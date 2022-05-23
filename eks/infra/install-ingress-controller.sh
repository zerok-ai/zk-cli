#----------- Install NGINX proxy based ingress ----------- 
# Add the Helm chart for Nginx Ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

# Install the Helm (v3) chart for nginx ingress controller
helm install app-ingress ingress-nginx/ingress-nginx \
	--namespace ingress \
	--create-namespace --set controller.replicaCount=5 \
	--set controller.nodeSelector."beta\.kubernetes\.io/os"=linux \
	--set defaultBackend.nodeSelector."beta\.kubernetes\.io/os"=linux \
	--set controller.metrics.enabled=true \
	--set controller.metrics.serviceMonitor.enabled=true 

# Print the Ingress Controller public IP address
kubectl get services --namespace ingress
#----------- 