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

#if [ -z "$PX_API_KEY" ] || [ -z "$PX_CLUSTER_KEY" ] || [ -z "$ZK_CLOUD_ADDR" ]
#then
#  echo "Invalid cli arguments. ERR #2"
#  exit 1
#fi

if [ -z "$ZK_SCENARIO_MANAGER_VERSION" ] || [ -z "$ZK_AXON_VERSION" ]
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

helm repo add zk-scenario-manager https://helm.zerok.ai/zk-client/zk-scenario-manager
helm repo update
helm upgrade zk-scenario-manager zk-scenario-manager/zk-scenario-manager --install --create-namespace --namespace zk-client --version $ZK_SCENARIO_MANAGER_VERSION

helm repo add zk-axon https://helm.zerok.ai/zk-client/zk-axon
helm repo update
helm upgrade zk-axon zk-axon/zk-axon --install --create-namespace --namespace zk-client --version $ZK_AXON_VERSION