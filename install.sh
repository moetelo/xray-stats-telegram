#!/bin/env bash

set -eo pipefail

USR_LOCAL_ETC=/usr/local/etc

DEFAULT_TRAFFIC_DATA_DIR=$USR_LOCAL_ETC/traffic-data
DEFAULT_USERS_JSON_PATH=$USR_LOCAL_ETC/xray-stats-telegram/users.json

trafficDataDir=$1
botToken=$2
usersJsonPath=$3

if [ "$1" == "--default" ]; then
    trafficDataDir=$DEFAULT_TRAFFIC_DATA_DIR
    usersJsonPath=$DEFAULT_USERS_JSON_PATH
fi

if [ -z "$trafficDataDir" ] || [ -z "$botToken" ] || [ -z "$usersJsonPath" ]; then
    echo "Usage: $0 <traffic-data-dir> <bot-token> <users-json-path>"
    echo "Usage: $0 --default <bot-token>"
    echo "(defaults to $DEFAULT_TRAFFIC_DATA_DIR and $DEFAULT_USERS_JSON_PATH)"
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
