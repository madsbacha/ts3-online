# ts3-online

```
version: '3.8'
services: 
  ts3-online:
    build: .
    environment: 
      - TS_HOST=host.docker.internal
      - TS_USERNAME=serveradmin
      - TS_PASSWORD=secret
    ports: 
      - 8080:8080
    restart: always
```
