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
zkctl activate -n 'zerok-demoapp' -r