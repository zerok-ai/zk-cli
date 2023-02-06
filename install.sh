alias zkctl="go run main.go"

# check the cluster name at the following locations
# 1. capacitor's install.sh
# 1. variables.sh

# install zerok operator
zkctl install operator

# install backend
zkctl install backend

# perform post backend installation tasks
zkctl install backend postsetup

# do rolling restart
kubectl label namespace zerok-demoapp zk-injection=enabled
kubectl scale deploy demo-shop-deployment --replicas=0 -n zerok-demoapp
kubectl scale deploy order --replicas=0 -n zerok-demoapp
kubectl scale deploy product --replicas=0 -n zerok-demoapp
kubectl scale deploy inventory --replicas=0 -n zerok-demoapp
kubectl scale deploy demo-shop-deployment --replicas=1 -n zerok-demoapp
kubectl scale deploy order --replicas=1 -n zerok-demoapp
kubectl scale deploy product --replicas=1 -n zerok-demoapp
kubectl scale deploy inventory --replicas=1 -n zerok-demoapp