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

if [ -z "$zkGptVersion" ] && [ -z "$zkPromtailVersion" ] && [ -z "$zkAxonVersion" ] && [ -z "$zkScenarioManagerVersion" ] && [ -z "$zkOtlpReceiverVersion" ] && [ -z "$zkDaemonsetVersion" ] && [ -z "$zkOperatorVersion" ] && [ -z "$zkWspClientVersion" ]
then
  echo "Invalid cli arguments. ERR #2"
  exit 1
fi

dep_names=()
dep_versions=()
if [ ! -z "$zkGptVersion" ]; then
  dep_names+=("zk-gpt")
  dep_versions+=($zkGptVersion)
fi

if [ ! -z "$zkPromtailVersion" ]; then
  dep_names+=("zk-promtail")
  dep_versions+=($zkPromtailVersion)
fi

if [ ! -z "$zkAxonVersion" ]; then
  dep_names+=("zk-axon")
  dep_versions+=($zkAxonVersion)
fi

if [ ! -z "$zkScenarioManagerVersion" ]; then
  dep_names+=("zk-scenario-manager")
  dep_versions+=($zkScenarioManagerVersion)
fi

if [ ! -z "$zkOtlpReceiverVersion" ]; then
  dep_names+=("zk-otlp-receiver")
  dep_versions+=($zkOtlpReceiverVersion)
fi

if [ ! -z "$zkDaemonsetVersion" ]; then
  dep_names+=("zk-daemonset")
  dep_versions+=($zkDaemonsetVersion)
fi

if [ ! -z "$zkOperatorVersion" ]; then
  dep_names+=("zk-operator")
  dep_versions+=($zkOperatorVersion)
fi

if [ ! -z "$zkWspClientVersion" ]; then
  dep_names+=("zk-wsp-client")
  dep_versions+=($zkWspClientVersion)
fi

for ((i = 0; i < ${#dep_names[@]}; ++i)); do
  dep_name="${dep_names[i]}"
  dep_version="${dep_versions[i]}"
  echo "Updating $dep_name to $dep_version"
  yq eval -i "(.dependencies[] | select(.name == \"$dep_name\").version) = \"$dep_version\"" $THIS_DIR/Chart.yaml
done

echo "Updated versions in Chart.yaml"
