#!/usr/bin/env bash

# Run jupyter in foreground if no commands specified
if [ -z "$1" ]; then
    echo "Start Jupyter"
else
    exec "$@"
fi