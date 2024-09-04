#!/bin/env bash

set -eo pipefail

botToken=$1
if [ -z "$botToken" ]; then
    echo "Usage: $0 <bot-token>"
    exit 1
fi

cargo build --release
cp ./target/release/xray-stats-telegram-rs /usr/local/bin/xray-stats-telegram

sed "s|<bot-token>|$botToken|" xray-stats-telegram.service \
    > /etc/systemd/system/xray-stats-telegram.service

systemctl daemon-reload
systemctl enable xray-stats-telegram
systemctl restart xray-stats-telegram
