# ts3-online

[![go report card](https://goreportcard.com/badge/github.com/madsbacha/ts3-online "go report card")](https://goreportcard.com/report/github.com/madsbacha/ts3-online)

The previous version, which is implemented in python and
uses a redis server, is located here:
https://github.com/madsbacha/ts3-online-counter

## docker-compose
```
version: '3.8'
services: 
  ts3-online:
    build: .
    environment: 
      - TS_HOST=host.docker.internal:10011
      - TS_USERNAME=serveradmin
      - TS_PASSWORD=secret
      - EXCLUDE_USERNAMES=SinusBot
    ports: 
      - 8080:8080
    restart: always
```

`EXCLUDE_USERNAMES` is a comma separated list of usernames to exclude from the list of online users.
