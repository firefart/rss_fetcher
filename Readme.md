# RSS Fetcher
This little GO program is intended to fetch all configured RSS or ATOM feeds every hour (configurable) and send new entries per E-Mail.

This project is mainly written because IFTT can not handle crt.sh feeds :/

Expected errors during execution are also sent via E-Mail to the E-Mail address configured in `config.json`.

For sending mails you should setup a local SMTP server like postfix to handle resubmission, signing and so on for you. SMTP authentication is currently not implemented.

## Installation on a systemd based system
* Build binary or download it
```bash
make
```
or
```bash
go get gopkg.in/gomail.v2
go get github.com/mmcdole/gofeed
go build
```
or get downloadlink from [https://github.com/FireFart/rss_fetcher/releases/](https://github.com/FireFart/rss_fetcher/releases/) and download
```
wget https://github.com/FireFart/rss_fetcher/archive/....
```

* Add a user to run the binary
```bash
adduser --system rss
```

* Copy everything to home dir
```bash
cp -R checkout_dir /home/rss/
```

* Modify run time (if you want to run it at other intervalls)
```bash
vim /home/rss/rss_fetcher.timer
```

* Edit the config
```
cp /home/rss/config.json.sample /home/rss/config.json
vim /home/rss/config.json
```

* Install the service and timer files
```bash
./install_service.sh
```

* Watch the logs
```bash
journalctl -u rss_fetcher.service -f
```
