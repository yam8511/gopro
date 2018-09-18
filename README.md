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

---

## 建立Humpback(提供界面可以看)

```shell
# 啟動代理伺服器與界面
$ docker-compose up -d humpback-web
```

---

## 設定/etc/hosts

```shell
# 開啟 hosts 列表
$ sudo vi /etc/hosts

# 加入以下 host
127.0.0.1   registry
```

---

## 設定私有容器倉庫的來源

1. 開啟瀏覽器 - [http://127.0.0.1:8000](http://127.0.0.1:8000)
2. 輸入帳密登入 (帳號: admin 密碼 123456)
3. 設定[Manage] > [System Config] > [Enabel Private Registry]
4. 填入網址[http://registry:5000](http://registry:5000)
5. 點擊[Save]之後，重新整理可以看到「Hub」，點擊進去可看到映像檔名稱，即成功

---

## 額外補充

1. [Portainer](https://portainer.io) - 可以非常清楚以GUI方式呈現目前本機的容器所有狀態
2. [Rancher](https://Rancher.com) - 可以將Kubernetes的狀態以GUI方式呈現
3. [Humpback](https://humpback.github.io/humpback) - 主要特別在可以呈現私有容器倉庫，不過也可以呈現本機容器狀態
