# create cluster with roles and permissions
sh ./create-and-configure-eks-cluster.sh

# install and configure kubernetes addons
sh ./install-and-configure-kubernetes-addons.sh

# install nginx ingress controller
sh ./install-ingress-controller.sh

# install prometheus and grafana
# sh install-prometheus-and-grafana.sh
