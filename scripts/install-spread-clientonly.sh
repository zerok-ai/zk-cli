#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [ "$#" -eq "0" ]; then
  echo "Invalid cli arguments. ERR #1"
  exit 1
fi

while [[ "$#" > "0" ]]
do
  case $1 in
    (*=*) eval $1;;
  esac
shift
done

if [ -z "$ZK_SCENARIO_MANAGER_VERSION" ] || [ -z "$ZK_OTLP_RECIEVER_VERSION" ] || [ -z "$ZK_OPERATOR_VERSION" ]
then
  echo "Invalid cli arguments. ERR #2"
  exit 1
fi


echo "checking helm binary"
if ! helm version; then
    echo "helm not available. ERR #4"
    exit 1
else
    echo "helm binary found."
fi

# add all helm repos
helm repo add zk-redis https://helm.zerok.ai/zk-client/zk-redis
helm repo add zk-scenario-manager https://helm.zerok.ai/zk-client/zk-scenario-manager
helm repo add zk-operator https://helm.zerok.ai/zk-client/zk-operator
helm repo add zk-otlp-receiver https://helm.zerok.ai/zk-client/zk-otlp-receiver

# update
helm repo update

# install
helm upgrade zk-redis zk-redis/zk-redis --install --create-namespace --namespace zk-client --version 0.1.0-alpha --wait
helm upgrade zk-operator zk-operator/zk-operator --install --create-namespace --namespace zk-client --version $ZK_OPERATOR_VERSION
helm upgrade zk-scenario-manager zk-scenario-manager/zk-scenario-manager --install --create-namespace --namespace zk-client --version $ZK_SCENARIO_MANAGER_VERSION
helm upgrade zk-otlp-receiver zk-otlp-receiver/zk-otlp-receiver --install --create-namespace --namespace zk-client --version $ZK_OTLP_RECIEVER_VERSION

