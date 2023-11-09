#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

while [[ "$#" > "0" ]]
do
  case $1 in
    (*=*) eval $1;;
  esac
shift
done

if [ -z "$PX_API_KEY" ]
then
  echo "Invalid zk-client arguments. ERR #2.1"
  exit 1
fi

if [ -z "$PX_CLUSTER_KEY" ]
then
  echo "Invalid zk-client arguments. ERR #2.2"
  exit 1
fi

if [ -z "$ZK_CLOUD_ADDR" ]
then
  echo "Invalid zk-client arguments. ERR #2.3"
  exit 1
fi

if [ -z "$PX_CLUSTER_ID" ]
then
  echo "Invalid zk-client arguments. ERR #2.4"
  exit 1
fi

echo "checking helm binary"
if ! helm version; then
    echo "helm not available. ERR #4"
    exit 1
else
    echo "helm binary found."
fi

helm repo add zk-client https://helm.zerok.ai/zk-client/zk-cli
helm repo update
helm upgrade zk-client zk-client/zk-client --install --create-namespace --namespace zk-client --version $ZK_HELM_VERSION --set=zk-scenario-manager.obfuscate=$OBFUSCATE_ENABLED --set=global.zkcloud.host=$ZK_CLOUD_ADDR --set=global.zkcloud.keys.cluster_key=$PX_CLUSTER_KEY --set=global.zkcloud.clusterId=$PX_CLUSTER_ID --set=global.zkcloud.keys.PX_API_KEY=$PX_API_KEY --set=global.gpt.enabled=$GPT_ENABLED
