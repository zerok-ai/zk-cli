#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "checking helm binary"
if ! helm version; then
    echo "helm not available. ERR #4"
    exit 1
else
    echo "helm binary found."
fi

while [[ "$#" > "0" ]]
do
  case $1 in
    (*=*) eval $1;;
  esac
shift
done

#helm dependency update $THIS_DIR
#helm upgrade $APP_NAME --install $THIS_DIR/ --create-namespace --namespace zk-client --wait

helm repo add zk-redis https://helm.zerok.ai/zk-client/zk-redis
helm repo update
helm upgrade zk-redis zk-redis/zk-redis --install --create-namespace --namespace zk-client --version 0.1.0-alpha

helm repo add zk-postgres https://helm.zerok.ai/zk-client/zk-postgres
helm repo update
helm upgrade zk-postgres zk-postgres/zk-postgres --install --create-namespace --namespace zk-client --version 0.1.0-alpha --wait

