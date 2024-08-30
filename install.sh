#!/bin/env bash

set -eo pipefail

trafficDataDir=$1
botToken=$2
usersJsonPath=$3

if [ -z "$trafficDataDir" ] || [ -z "$botToken" ] || [ -z "$usersJsonPath" ]; then
    echo "Usage: $0 <traffic-data-dir> <bot-token> <users-json-path>"
    exit 1
fi

go build -o /usr/local/bin/xray-stats-telegram

sed "s|<traffic-data-dir>|$trafficDataDir|g" xray-stats-telegram.service \
    | sed "s|<bot-token>|$botToken|" \
    | sed "s|<users-json>|$usersJsonPath|" \
    > /etc/systemd/system/xray-stats-telegram.service

systemctl daemon-reload
systemctl enable xray-stats-telegram
systemctl restart xray-stats-telegram
