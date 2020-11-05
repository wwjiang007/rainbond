#!/bin/bash
set -e

if [ "$1" = "bash" ];then
    exec /bin/bash
elif [ "$1" = "version" ];then
    echo /root/rainbond-mesh-data-panel version
else
    env2file conversion -f /root/envoy_config.yaml
    cluster_name=${TENANT_ID}_${PLUGIN_ID}_${SERVICE_NAME}
    # start sidecar process
    /root/rainbond-mesh-data-panel&
    # start envoy process
    exec envoy -c /root/envoy_config.yaml --service-cluster ${cluster_name} --service-node ${cluster_name}     
fi
