# ts3-online

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
