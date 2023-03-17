#!/usr/bin/env bash

# Run jupyter in foreground if no commands specified
if [ -z "$1" ]; then
    PARAMS=""
    if [ -z "$JUPYTERLAB_PORT" ]; then
        PARAMS="$PARAMS --port=$JUPYTERLAB_PORT"
    fi
    if [ -z "$JUPYTERLAB_TOKEN" ]; then
        PARAMS="$PARAMS --NotebookApp.token=$JUPYTERLAB_TOKEN"
    fi
    jupyter lab $PARAMS
else
    exec "$@"
fi
