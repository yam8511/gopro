# 容器私有倉庫

[參考網址](http://blog.poetries.top/handbook/html/%E6%9C%8D%E5%8A%A1%E7%AB%AF/docker.html#t13Docker%E7%A7%81%E6%9C%89%E4%BB%93%E5%BA%93%E6%90%AD%E5%BB%BA)

---

## 使用Docker官方映像檔建立私有倉庫

```shell
# 啟動容器
$ docker-compose up -d registry
# 測試倉庫
$ curl http://127.0.0.1:5000/v2/_catalog
{"repositories":[]}
```

---

## 建立一個映像檔(Image)

```shell
# 建立映像檔，詳細請看 docker-compose.yml
$ docker-compose up --build -d my-image
# 查看有沒有映像檔
$ docker images | grep gopro-image-test
127.0.0.1:5000/gopro-image-test    alpine    d23ad2567c5a    53 seconds ago    4.41MB
```

---

## 上傳映像檔到私有倉庫

```shell
# 上傳映像檔
$ docker push 127.0.0.1:5000/gopro-image-test:alpine
The push refers to repository [127.0.0.1:5000/gopro-image-test]
...

# 檢查是否有成功
$ curl http://127.0.0.1:5000/v2/_catalog
{"repositories":["gopro-image-test"]}

# 也可以嘗試下載映像檔
$ docker pull 127.0.0.1:5000/gopro-image-test:alpine
```
