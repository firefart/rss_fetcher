#!/bin/sh

echo "Copying unit file"
cp /home/rss/rss_fetcher.service /etc/systemd/system/rss_fetcher.service
cp /home/rss/rss_fetcher.timer /etc/systemd/system/rss_fetcher.timer
echo "reloading systemctl"
systemctl daemon-reload
echo "enabling service"
systemctl enable rss_fetcher.timer
systemctl start rss_fetcher.timer
systemctl start rss_fetcher.service
systemctl status rss_fetcher.service
systemctl status rss_fetcher.timer
systemctl list-timers --all
