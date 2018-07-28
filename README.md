# Docker-Compose Demo

1. 啟動所有服務
```shell
$ docker-compose up -d

# 啟動指定服務
# docker-compose up -d [服務名稱...]
# docker-compose up -d web pgsql
```

2. 查看服務狀態
```shell
$ docker-compose ps

# 查看指定服務
# docker-compose ps [服務名稱...]
# docker-compose ps web pgsql
```

3. 查看服務CPU用量
```shell
$ docker-compose top

# 查看指定服務
# docker-compose top [服務名稱...]
# docker-compose top web pgsql
```

4. 停止/重啟服務
```shell
# 停止
$ docker-compose stop [服務名稱...]
# docker-compose stop web

# 重啟
$ docker-compose restart [服務名稱...]
# docker-compose restart web
```

5. 水平擴展服務
```shell
docker-compose up -d --scale [服務名稱]=[數量]
# docker-compose up -d --scale web=3  --scale adminer=2
```

6. 關閉所有服務
```shell
$ docker-compose down
```

---

# 觀念
1. 所有的「服務名稱」皆等同容器的「IP」
2. 容器內的世界，用容器內服務的「名稱」與「Port」互Call
3. 容器外的世界，只能Call「對外開的Port」
