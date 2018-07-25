# Heroku Demo

0. 安裝 Heroku CLI (指令列)
參考文件 [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli)

1. Heroku 登入
```shell
$ heroku login
```

2. 建立Heroku專案
```
heroku create
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

6. 建置DB (PostgreSQL)
```shell
$ heroku addons:create heroku-postgresql:hobby-dev
Creating heroku-postgresql:hobby-dev on ⬢ young-reef-51042... free
Database has been created and is available
 ! This database is empty. If upgrading, you can transfer
 ! data from another database with pg:copy
Created postgresql-angular-16538 as DATABASE_URL
```

7. 看DB連線設定
- https://data.heroku.com/
- 「Setting」 -> 「View Credentials」
