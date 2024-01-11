#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "checking helm binary"
if ! helm version; then
    echo "helm not available. ERR #4"
    exit 1
else
    echo "helm binary found."
fi

helm dependency build $THIS_DIR
helm dependency update $THIS_DIR
helm upgrade zk-client $THIS_DIR/ --install --create-namespace --namespace zk-client --version $ZK_HELM_VERSION --dry-run

