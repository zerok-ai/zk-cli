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

if [ -z "$PX_API_KEY" ] || [ -z "$PX_CLUSTER_KEY" ] || [ -z "$ZK_CLOUD_ADDR" ]
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

helm dependency build $THIS_DIR
helm dependency update $THIS_DIR
helm --install --set=zk-axon.releaseName=zk-axon --set=zk-scenario-manager.releaseName=zk-scenario-manager --set=global.vizier.tag=$VIZIER_TAG --set=global.zkcloud.host=$ZK_CLOUD_ADDR --set=global.data.cluster_key=$PX_CLUSTER_KEY --set=global.data.PX_API_KEY=$PX_API_KEY upgrade $APP_NAME $THIS_DIR/ --create-namespace --namespace zk-client