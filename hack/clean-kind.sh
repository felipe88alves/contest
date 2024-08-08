#! /bin/bash

CLUSTERS=$(kind get clusters)

if [ ! -z ${CLUSTERS} ]; then
  kind delete cluster --name ${CLUSTERS}
fi
