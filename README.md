# Traefik Quick Start

[官方教學文件](https://docs.traefik.io/#the-trfik-quickstart-using-docker)

1. docker-compose up -d reverse-proxy
    see [http://127.0.0.1:8080](http://127.0.0.1:8080)

2. 啟動一個server
    ```shell
    docker-compose up -d whoami
    curl -H Host:demo1 http://127.0.0.1
    ```

3. 啟動兩個server，會發現有Load Balance
    ```shell
    docker-compose up -d --scale whoami=2 whoami
    curl -H Host:demo1 http://127.0.0.1
    curl -H Host:demo1 http://127.0.0.1
    ```

4. 啟動三個server2，依照Host決定呼叫哪一方server
    ```shell
    docker-compose up -d --scale whoami2=3 whoami2
    curl -H Host:demo2 http://127.0.0.1
    ```
