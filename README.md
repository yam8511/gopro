# Heroku Demo

0. 安裝 Heroku CLI (指令列)
參考文件 [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli)

1. Heroku 登入
```shell
$ heroku login
```

2. 建立Heroku專案
```
$ heroku create
Creating app... done, ⬢ young-reef-51042
https://young-reef-51042.herokuapp.com/ | https://git.heroku.com/young-reef-51042.git
```


3. 加入專案 git remote
```shell
$ heroku git:remote -a [專案名稱]
# heroku git:remote -a young-reef-51042
```

4. 推送專案 (用Dockerfile建置)
```shell
$ heroku container:login

# 預設推送 Dockerfile
$ heroku container:push web

# 推送 Dockerfile.demo
$ heroku container:push demo --recursive
```

5. 發佈專案 (建置好的專案)
```shell
$ heroku container:release web demo
```

6. 開啟頁面＆監看紀錄
```
$ heroku open
$ heroku logs -t
```

---

# 架設DB (PostgreSQL)

1. 架設DB (PostgreSQL)
```shell
$ heroku addons:create heroku-postgresql:hobby-dev
Creating heroku-postgresql:hobby-dev on ⬢ young-reef-51042... free
Database has been created and is available
 ! This database is empty. If upgrading, you can transfer
 ! data from another database with pg:copy
Created postgresql-angular-16538 as DATABASE_URL
```

2. 看DB連線設定
- https://data.heroku.com/
- 「Setting」 -> 「View Credentials」

3. 開啟 adminer 看資料庫
```shell
docker-compose up -d adminer
```
ps. 開啟 http://127.0.0.1:8080

4. 開API，建立使用者資料
- http://[專案的Domain]/api/create/user
可以看到回傳的User資料
