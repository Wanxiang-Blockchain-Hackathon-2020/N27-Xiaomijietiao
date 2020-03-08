#!/bin/bash

if [ "$PROC_NAME" == "" ]; then
    echo "No proc name set."
    exit -1
fi

echo "$config" > /opt/weyom/config/config.json

ldconfig
exec /opt/weyom/bin/$PROC_NAME /opt/weyom/config/config.json
