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

if [ -z "$ZK_SCENARIO_MANAGER_VERSION" ] || [ -z "$ZK_AXON_VERSION" ] || [ -z "$ZK_DAEMONSET_VERSION" ] || [ -z "$ZK_GPT_VERSION" ] || [ -z "$ZK_WSP_CLIENT_VERSION" ]  || [ -z "$ZK_OPERATOR_VERSION" ] || [ -z "$ZK_CLOUD_ADDR" ]
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

helm repo add zk-daemonset https://helm.zerok.ai/zk-client/zk-daemonset
helm repo update
helm upgrade zk-daemonset zk-daemonset/zk-daemonset --install --create-namespace --namespace zk-client --version $ZK_DAEMONSET_VERSION

helm repo add zk-gpt https://helm.zerok.ai/zk-client/zk-gpt
helm repo update
helm upgrade zk-gpt zk-gpt/zk-gpt --install --create-namespace --namespace zk-client --version $ZK_GPT_VERSION

helm repo add zk-wsp-client https://helm.zerok.ai/zk-client/zk-wsp-client
helm repo update
helm upgrade zk-wsp-client zk-wsp-client/zk-wsp-client --install --create-namespace --namespace zk-client --version $ZK_WSP_CLIENT_VERSION --set=global.zkcloud.host=$ZK_CLOUD_ADDR


helm repo add zk-operator https://helm.zerok.ai/zk-client/zk-operator
helm repo update
helm upgrade zk-operator zk-operator/zk-operator --install --create-namespace --namespace zk-client --version $ZK_OPERATOR_VERSION
