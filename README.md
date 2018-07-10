# Traefik Quick Start
[官方教學文件](https://docs.traefik.io/#the-trfik-quickstart-using-docker)

1. docker-compose up -d reverse-proxy
see [http://127.0.0.1:8080](http://127.0.0.1:8080)

2. 啟動一個server
```
docker-compose up -d whoami
curl -H Host:whoami.docker.localhost http://127.0.0.1
```

3. 啟動兩個server，會發現有Load Balance
```
docker-compose up -d --scale whoami=2 whoami
curl -H Host:whoami.docker.localhost http://127.0.0.1
curl -H Host:whoami.docker.localhost http://127.0.0.1
```

4. 啟動三個server2，依照Host決定呼叫哪一方server
```
docker-compose up -d --scale whoami2=3 whoami2
curl -H Host:whoami2.docker.localhost http://127.0.0.1
```
