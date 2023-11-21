#!/bin/bash

THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

while [[ "$#" > "0" ]]
do
  case $1 in
    (*=*) eval $1;;
  esac
shift
done

if [ -n "$install" ] && [ "$install" = "true" ]; then

  if [ -n "$pgurl" ] && [ -n "$redisAuth" ]; then
    echo "Install debug helm"
    helm repo add zk-debug https://helm.zerok.ai/zk-client/zk-debug
    helm repo update
    helm upgrade zk-debug zk-debug/zk-debug --install --create-namespace --namespace zk-client --set=zkDebug.redisAuth=$redisAuth --set=zkdebug.postgrest.postgresUri=$pgurl
  else
    echo "Invalid zk-client debug arguments. ERR #1"
    exit 1
  fi

elif [ -n "$delete" ] && [ "$delete" = "true" ]; then

  echo "delete debug helm"
  helm uninstall zk-debug --namespace zk-client

fi

