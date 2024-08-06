#!/usr/bin/env bash

if ! go build -o ./picklebot ./; then
    echo "Failed to build picklebot"
    exit 1
fi

