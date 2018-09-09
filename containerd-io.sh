#!/bin/bash

set -exuo pipefail

~/go/bin/dlv --help

CONTAINER_ID="$(crictl ps --label io.kubernetes.container.name=$1 -q)"
SHIM_PID=$(ps -eaf | grep ${CONTAINER_ID} | grep -v grep | awk '{print $2}')
ps --ppid ${SHIM_PID} -f --forest

STDOUT=$(find /run/containerd/io.containerd.grpc.v1.cri/containers/${CONTAINER_ID} -type p -name "*stdout")

bash -exuc "trap 'pkill cat' EXIT: cat /dev/urandom | tr -dc 'a-zA-Z0-9' > ${STDOUT}" &

set +e
~/go/bin/dlv --headless -l=:2345 --api-version=2 attach ${SHIM_PID}

kill $(jobs -p)
