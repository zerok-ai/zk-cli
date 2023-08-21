#!/bin/bash
THIS_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "checking helm binary"
if ! helm version; then
    echo "helm not available. ERR #4"
    exit 1
else
    echo "helm binary found."
fi

helm dependency update $THIS_DIR
#helm --install $APP_NAME $THIS_DIR/ --create-namespace --namespace zk-client
helm upgrade $APP_NAME --install $THIS_DIR/ --create-namespace --namespace zk-client --wait
