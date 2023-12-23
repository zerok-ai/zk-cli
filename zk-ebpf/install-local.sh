#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
if [ "$#" -eq "0" ]; then
  echo "Invalid zk-client local arguments. ERR #1"
  exit 1
fi

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

echo "ZK_HELM_VERSION=$ZK_HELM_VERSION"
echo "ZK_CLUSTER_ID=$ZK_CLUSTER_ID"
echo "ZK_CLUSTER_NAME=$ZK_CLUSTER_NAME"
echo "ZK_EBPF_JWT_KEY=$ZK_EBPF_JWT_KEY"

helm dependency build $THIS_DIR
helm dependency update $THIS_DIR
helm upgrade zk-ebpf $THIS_DIR/ --install --create-namespace --namespace pl --version $ZK_HELM_VERSION --set=plClusterSecrets.clusterId=$ZK_CLUSTER_ID --set=plClusterSecrets.clusterName=$ZK_CLUSTER_NAME --set=plClusterSecrets.jwtSigningKey=$ZK_EBPF_JWT_KEY