#----------- Install NGINX proxy based ingress ----------- 

# Add the Helm chart for Nginx Ingress
echo '---------------------- Updating helm repo for ingress-nginx'
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

# Create namespace
echo '---------------------- Creating namespace `ingress`'
kubectl create namespace ingress

# Install the Helm (v3) chart for nginx ingress controller
echo '---------------------- Installing helm chart of ingress-nginx'
helm upgrade --install app-ingress ingress-nginx/ingress-nginx \
	--namespace ingress \
	--values ./yaml/values/nginx-ingress-values.yaml

# Print the Ingress Controller public IP address
kubectl get services --namespace ingress

#----------- 
# helm upgrade app-ingress ingress-nginx/ingress-nginx \
# 	--namespace ingress \
# 	-f ./yaml/values/nginx-ingress-values.yaml