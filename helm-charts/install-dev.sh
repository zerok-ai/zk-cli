#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [ "$#" -eq "0" ]; then
  echo "Invalid zk-client dev arguments. ERR #1"
  exit 1
fi

while [[ "$#" > "0" ]]
do
  case $1 in
    (*=*) eval $1;;
  esac
shift
done

if [ -z "$ZK_SCENARIO_MANAGER_VERSION" ] || [ -z "$ZK_AXON_VERSION" ] || [ -z "$ZK_PROMTAIL_VERSION" ] || [ -z "$ZK_OTLP_RECIEVER_VERSION" ] || [ -z "$ZK_DAEMONSET_VERSION" ] || [ -z "$ZK_WSP_CLIENT_VERSION" ]  || [ -z "$ZK_OPERATOR_VERSION" ] || [ -z "$ZK_APP_INIT_CONTAINERS_VERSION" ] || [ -z "$ZK_CLOUD_ADDR" ] || [ -z "$PX_CLUSTER_KEY" ] || [ -z "$PX_API_KEY" ] || [ -z "$PX_CLUSTER_ID" ]
then
  echo "Invalid zk-client dev arguments. ERR #2"
  exit 1
fi

# Check if GPT_ENABLED is true, then ZK_GPT_VERSION should be present
if [ "$GPT_ENABLED" = "true" ] && [ -z "$ZK_GPT_VERSION" ]
then
  echo "Invalid zk-client dev arguments. ERR #3"
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
helm repo add zk-scenario-manager https://helm.zerok.ai/zk-client/zk-scenario-manager
helm repo add zk-axon https://helm.zerok.ai/zk-client/zk-axon
helm repo add zk-daemonset https://helm.zerok.ai/zk-client/zk-daemonset
helm repo add zk-wsp-client https://helm.zerok.ai/zk-client/zk-wsp-client
helm repo add zk-operator https://helm.zerok.ai/zk-client/zk-operator
if [ "$GPT_ENABLED" = "true" ]
then
  helm repo add zk-gpt https://helm.zerok.ai/zk-client/zk-gpt
fi
helm repo add zk-promtail https://helm.zerok.ai/zk-client/zk-promtail
helm repo add zk-otlp-receiver https://helm.zerok.ai/zk-client/zk-otlp-receiver

# update
helm repo update

# install
helm upgrade zk-wsp-client zk-wsp-client/zk-wsp-client --install --create-namespace --namespace zk-client --version $ZK_WSP_CLIENT_VERSION --set=global.zkcloud.host=$ZK_CLOUD_ADDR
helm upgrade zk-operator zk-operator/zk-operator --install --create-namespace --namespace zk-client --version $ZK_OPERATOR_VERSION --set=global.zkcloud.host=$ZK_CLOUD_ADDR --set=global.zkcloud.keys.cluster_key=$PX_CLUSTER_KEY --set=global.zkcloud.keys.PX_API_KEY=$PX_API_KEY
helm upgrade zk-scenario-manager zk-scenario-manager/zk-scenario-manager --install --create-namespace --namespace zk-client --version $ZK_SCENARIO_MANAGER_VERSION
helm upgrade zk-axon zk-axon/zk-axon --install --create-namespace --namespace zk-client --version $ZK_AXON_VERSION
helm upgrade zk-daemonset zk-daemonset/zk-daemonset --install --create-namespace --namespace zk-client --version $ZK_DAEMONSET_VERSION
if [ "$GPT_ENABLED" = "true" ]
then
  helm upgrade zk-gpt zk-gpt/zk-gpt --install --create-namespace --namespace zk-client --version $ZK_GPT_VERSION
fi
helm upgrade zk-promtail zk-promtail/zk-promtail --install --create-namespace --namespace zk-client --version $ZK_PROMTAIL_VERSION --set=global.zkcloud.clusterId=$PX_CLUSTER_ID
helm upgrade zk-otlp-receiver zk-otlp-receiver/zk-otlp-receiver --install --create-namespace --namespace zk-client --version $ZK_OTLP_RECIEVER_VERSION


