# Drone + Gitea

0. 前言
此範例是在自己的本機架設Git伺服器(Gitea)，再架設Drone連上自己的Git伺服器。
若不想在自己本機架設Git伺服器，可跳過Git伺服器架設的步驟。(直接跳到步驟5)
但須注意，必須讓外部Git伺服器**有辦法連線**到自己的Drone伺服器

1. 先設定一下 hosts
```
$ sudo vi /etc/hosts
```
加入以下字串
```
127.0.0.1 drone-server
```

2. 啟動 gitea、pgsql
```
docker-compose up -d gitea pgsql
```

3. 開啟 gitea 安裝頁面，進行安裝
- http://127.0.0.1:3000/install

    3-1. 資料庫設定 (Database Settings)
    ```
    DataType: PostgreSQL
    Host: pgsql:5432
    Username: root
    Password: qwe123
    Database Name: gitea
    ```

    3-2-1. 找本機IP [192.168.xxx.xxx]
    ```
    # 使用 ip
    $ ip a | grep 192.168.
    # or 使用 ifconfig
    $ ifconfig | grep 192.168.
    ```

    3-2-2. 一般設定 (General Settings)
    ```
    Gitea Base URL: http://192.168.xxx.xxx:3000/
    # 例如: http://192.168.2.137:3000/
    ```

    3-3. 可選設定 (Optional Setting)
    ```
    Administrator Username: demo
    Password: qwe123
    Confirm Password: qwe123
    Email Address: demo@demo.com
    ```

    3-4. 按下「Install Gitea」
    ps. 安裝時候，需要等待一段時間，可以看Log的進度
    ```
    docker-compose logs -f gitea
    ```


4. 登入 Gitea， 建立一個Repo， 命名「demo」
- http://127.0.0.1:3000/user/login?redirect_to=


5. 啟動 drone-server、drone-agent
```
docker-compose up -d drone-server drone-agent
```

6. 開啟 drone-server，並以 gitea 帳號登入
- http://drone-server:8000
```
Username: demo
Password: qwe123
```

7. 在 drone 頁面，開啟 demo/demo 的開關 (右側的小圈圈)
ps. 檢查 Gitea Webhook 是否有新連結
- http://127.0.0.1:3000/demo/demo/settings/hooks


8. 進入 demo 資料夾，上傳 git
```
cd demo/
git init
git commit -m "first commit"
git remote add origin http://192.168.xxx.xxx:3000/demo/demo.git
git push -u origin master
```
ps. 192.168.xxx.xxx 換成自己的本機IP
請參考 [demo/demo](http://127.0.0.1:3000/demo/demo.git)


9. 查看 drone 自動測試狀態、以及 git 專案呈現的狀態
- http://drone-server:8000/demo/demo/1
- http://127.0.0.1:3000/demo/demo
ps. 若沒有看見東西，可以重新整理看看
