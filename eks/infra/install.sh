# create cluster with roles and permissions
echo '###################### Installing cluster'
sh ./create-and-configure-eks-cluster.sh

# install and configure kubernetes addons
echo '###################### Installing addons'
sh ./install-and-configure-kubernetes-addons.sh

# install prometheus and grafana
echo '###################### Installing prometheus and grafana'
sh ./install-prometheus-and-grafana.sh

# install ingress controller
echo '###################### Installing ingress controller'
sh ./install-ingress-controller.sh


