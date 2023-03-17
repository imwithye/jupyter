#!/bin/bash

ARCH=$(uname -p)

cd /tmp/
wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-${ARCH}.sh
bash Miniconda3-latest-Linux-${ARCH}.sh -b -p /opt/miniconda3
rm -r /tmp/Miniconda3-latest-Linux-${ARCH}.sh
