#! /bin/bash

ENVTEST_K8S_VERSION=$1
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m | sed -e 's/aarch64$/arm64/' -e 's/x86_64/amd64/' )"

PID=($(ps aux | grep "[b]in/k8s/${ENVTEST_K8S_VERSION}-${OS}-${ARCH}" | awk '{print $2}'))

if [ ! -z ${PID} ]; then
  kill ${PID[@]}
fi
