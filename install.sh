#!/bin/env bash

set -eo pipefail

botToken=$1

go build -o /usr/local/bin/xray-stats-telegram

if [ -n "$botToken" ]; then
    sed "s|<bot-token>|$botToken|" xray-stats-telegram.service \
        > /etc/systemd/system/xray-stats-telegram.service
fi

mkdir -p /usr/local/etc/xray-stats-telegram
touch /usr/local/etc/xray-stats-telegram/admins
touch /usr/local/etc/xray-stats-telegram/users

systemctl enable xray-stats-telegram
systemctl restart xray-stats-telegram
