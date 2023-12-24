#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

while [[ "$#" > "0" ]]
do
  case $1 in
    (*=*) eval $1;;
  esac
shift
done

if [ -z "$ZK_HELM_VERSION" ]
then
  echo "Invalid zk-client ebpf arguments. ERR #2.0"
  exit 1
fi

if [ -z "$ZK_CLUSTER_ID" ]
then
  echo "Invalid zk-client ebpf arguments. ERR #2.1"
  exit 1
fi

if [ -z "$ZK_CLUSTER_NAME" ]
then
  echo "Invalid zk-client ebpf arguments. ERR #2.2"
  exit 1
fi

if [ -z "$ZK_EBPF_JWT_KEY" ]
then
  echo "Invalid zk-client ebpf arguments. ERR #2.3"
  exit 1
fi

#if [ -z "$ZK_EBPF_CA_CERT" ]
#then
#  echo "Invalid zk-client ebpf arguments. ERR #2.4"
#  exit 1
#fi
#
#if [ -z "$ZK_EBPF_CLIENT_CERT" ]
#then
#  echo "Invalid zk-client ebpf arguments. ERR #2.5"
#  exit 1
#fi
#
#if [ -z "$ZK_EBPF_CLIENT_KEY" ]
#then
#  echo "Invalid zk-client ebpf arguments. ERR #2.6"
#  exit 1
#fi
#
#if [ -z "$ZK_EBPF_SERVER_CERT" ]
#then
#  echo "Invalid zk-client ebpf arguments. ERR #2.7"
#  exit 1
#fi
#
#if [ -z "$ZK_EBPF_SERVER_KEY" ]
#then
#  echo "Invalid zk-client ebpf arguments. ERR #2.8"
#  exit 1
#fi

echo "checking helm binary"
if ! helm version; then
    echo "helm not available. ERR #4"
    exit 1
else
    echo "helm binary found."
fi

helm repo add zk-ebpf https://helm.zerok.ai/zk-client/zk-ebpf
helm repo update
#helm upgrade zk-ebpf zk-ebpf/zk-ebpf --install --create-namespace --namespace pl --version $ZK_HELM_VERSION --set=plClusterSecrets.clusterId=$ZK_CLUSTER_ID --set=plClusterSecrets.clusterName=$ZK_CLUSTER_NAME --set=plClusterSecrets.jwtSigningKey=$ZK_EBPF_JWT_KEY --set=serviceTlsCerts.caCrt=$ZK_EBPF_CA_CERT --set=serviceTlsCerts.clientCrt=$ZK_EBPF_CLIENT_CERT --set=serviceTlsCerts.clientKey=$ZK_EBPF_CLIENT_KEY --set=serviceTlsCerts.serverCrt=$ZK_EBPF_SERVER_CERT --set=serviceTlsCerts.serverKey=$ZK_EBPF_SERVER_KEY
helm upgrade zk-ebpf zk-ebpf/zk-ebpf --install --create-namespace --namespace zk-client --version $ZK_HELM_VERSION --set=plClusterSecrets.clusterId=$ZK_CLUSTER_ID --set=plClusterSecrets.clusterName=$ZK_CLUSTER_NAME --set=plClusterSecrets.jwtSigningKey=$ZK_EBPF_JWT_KEY