# RSS Fetcher

This little GO program is intended to fetch all configured RSS or ATOM feeds every hour (configurable) and send new entries per E-Mail.

This project is mainly written because IFTT can not handle crt.sh feeds :/

Expected errors during execution are also sent via E-Mail to the E-Mail address configured in `config.json`.

For sending mails you should setup a local SMTP server like postfix to handle resubmission, signing and so on for you. SMTP authentication is currently not implemented.

The program keeps the last date of the last entry per feed in it's database to compare it to on the next fetch.
We can't just use the current date because crt.sh is caching it's feeds and they do not appear at the time written in the feed.

## Installation on a systemd based system

- Build binary or download it

```bash
make
```

or

```bash
go get
go build
```

or

```bash
make_linux.bat
make_windows.bat
```

- Add a user to run the binary

```bash
adduser --system rss
```

- Copy everything to home dir

```bash
cp -R checkout_dir /home/rss/
```

- Modify run time (if you want to run it at other intervalls)

```bash
vim /home/rss/rss_fetcher.timer
```

- Edit the config

```bash
cp /home/rss/config.json.sample /home/rss/config.json
vim /home/rss/config.json
```

- Install the service and timer files

```bash
./install_service.sh
```

- Watch the logs

```bash
journalctl -u rss_fetcher.service -f
```

## Config Sample

```json
{
  "timeout": 10,
  "mailserver": "localhost",
  "mailport": 25,
  "mailfrom": "RSS <a@a.com>",
  "mailto": "People <b@b.com>",
  "mailonerror": true,
  "mailtoerror": "c@c.com",
  "database": "rss.db",
  "globalignorewords": ["ignore1", "ignore2"],
  "feeds": [
    {
      "title": "Certificates *.aaa.com",
      "url": "https://crt.sh/atom?q=%25.aaa.com",
      "ignorewords": ["[Precertificate]", "ignore2"]
    },
    {
      "title": "Certificates *.bbb.com",
      "url": "https://crt.sh/atom?q=%25.bbb.com"
    }
  ]
}
```
