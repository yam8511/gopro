# Expected System Flow

預期的完整系統流程

1. 版本控制系統 (Version Control System)
2. 容器管理系統 (Container Management System)
3. 自動測試部署 (CI/CD)
4. 日誌管理系統 (Log Management System)

p.s.  
照理來說，應該還需要有個「監控系統」，來監督與分析每台機器或應用程序的負載情況。  
但因為使用「容器管理系統」，有機會直接運用容器管理系統進行負載查看。  

---

## 前提需求

1. 找本機IP [192.168.xxx.xxx]

    ```shell
    # 使用 ip
    $ ip a | grep 192.168.
    # or 使用 ifconfig
    $ ifconfig | grep 192.168.
    ```

2. 設定 /etc/hosts

    p.s. 切記，這邊不指定 127.0.0.1，是有用意的

    ```shell
    # for Version Management
    192.168.xxx.xxx  gitea
    # for CI/CD
    192.168.xxx.xxx  drone-server
    # for Container Management
    192.168.xxx.xxx  registry
    # for Log Management
    192.168.xxx.xxx  fluentd kibana
    ```

3. 此文假設讀者已經熟悉 Git 相關的基本概念
4. 此文假設讀者已經熟悉 Docker 與 Docker Compose 等相關容器技術的基本操作
5. 此文假設讀者已經理解 CI/CD 相關的基本概念
6. 此文假設讀者知道為何需要收集日誌
7. 此文假設讀者是在 ubuntu、mint...等linux相關作業系統運行，Mac或Windows一概不負責XD

---

## 版本控制系統 (Version Control System)

這邊選擇使用 Gitea 當做個人的 Git Server，因為開源、設定簡單、用Golang撰寫的。

1. 啟動 gitea、pgsql

    ```shell
    docker-compose up -d gitea pgsql
    ```

2. 開啟 gitea 安裝頁面，進行安裝 [[傳送門](http://gitea:3000/install)]

    2-1. 資料庫設定 (Database Settings)

        ```md
        DataType: PostgreSQL
        Host: pgsql:5432
        Username: root
        Password: qwe123
        Database Name: gitea
        ```

    2-2. 一般設定 (General Settings)

        ```md
        # SSH設定
        SSH Server Domain: gitea
        SSH Server Port: 10022

        # HTTP設定
        Gitea Base URL: http://gitea:3000/
        ```

    2-3. 可選設定 (Optional Setting)

        ```md
        Administrator Username: demo
        Password: qwe123
        Confirm Password: qwe123
        Email Address: demo@demo.com
        ```

    2-4. 按下「Install Gitea」

        ps. 安裝時候，需要等待一段時間，可以看Log的進度

        ```shell
        docker-compose logs -f gitea
        ```

3. 測試Git伺服器是否成功，並建立Repo，CI/CD即將用到

    3-1. 開啟 Gitea 頁面 [[傳送門](http://gitea:3000)]
    3-2. 建立一個Repo，取名為 demo
    3-3. 進入已經建立好的demo資料夾，並且推上 git

### Gitea 錯誤處理

如果有錯誤發生，可能有以下情況，一一確認之後，應該就會好起來的。

1. 如果頁面打不開，請確認 gitea 是否有啟動成功
2. 如果 install 失敗，確認 PostgreSQL 是否有啟動成功
3. 如果連線失敗，確認 /etc/hosts 有沒有設定「前提準備」所指定的
4. 如果上面3點都確認OK，請刪除 data 資料夾，在從頭重試看看

---

## 容器管理系統 (Container Management System)

專案開發上，個人選擇使用 [Docker](https://www.docker.com/) 搭配 [Docker Compose](https://docs.docker.com/compose/)，進行最簡易的容器建立。  
而容器倉庫，則使用 Docker 官方提供的 [Registry 映像檔](https://hub.docker.com/_/registry/)，來建立私人倉庫。  
那或許會問為什摩不使用 [Kubernetes](https://kubernetes.io/) 呢？ 用！當然用！  
只是目前重點放在「最簡易的容器建立」，所以採取 Docker Compose。  
但是如果是在正式環境，則推薦使用 Kubernetes！  

目前先來建立個人的容器倉庫：

1. 啟動 registry

    ```shell
    docker-compose up -d registry
    ```

    ps. 啟動後，測試倉庫是否啟動成功
    ```shell
    curl http://registry:5000/v2/_catalog
    # {"repositories":[]}
    ```

2. 啟動 humpback-web (容器倉庫的GUI)

    ```shell
    docker-compose up -d humpback-web
    ```

3. 開啟頁面，進行設定 [[傳送門](http://registry:5252)]

    3-1. 登入頁面

        ```md
        預設帳號 admin
        預設密碼 123456
        ```

    3-2. 設定容器倉庫的網址

        點選 Manage > System Config > Enabel Private Registry  
        輸入 `http://registry:5000`  
        點擊 Save  

    3-3. 點選 Hub，查看儲存的映像檔

4. 測試能不能正常上傳映像檔

    注意！因為 Docker 上傳映像檔，預設使用 HTTPS 協定連線  
    所以上傳映像檔，需要先設定 docker， 讓他可忽略安全性

    4-1. Docker 設定

        4-1-1. 打開 daemon.json

            ```shell
            sudo vi /etc/docker/daemon.json
            ```

        4-1-2. 貼上以下設定

            ```json
            {
                "insecure-registries":["registry:5000"]
            }
            ```

        4-1-3. 重啟Docker

            ```shell
            sudo service docker restart
            ```

    4-2. 建立映像檔

        ```shell
        docker build ./demo -t registry:5000/myimg
        ```

    4-3. 上傳映像檔

        ```shell
        docker push registry:5000/myimg
        ```

    4-4. 點選 Hub，查看儲存的映像檔 [[傳送門](http://registry:5252/hub)]

### Registry 錯誤處理

1. 請先確認 registry 有沒有啟動成功
2. 如果無法連線，確認 /etc/hosts 有沒有設定「前提準備」所指定的
3. 如果無法建立映像檔，請確認是否在本專案目錄底下
4. 如果頁面打不開，確認 humpback-web 是否有啟動成功
5. 如果上面4點都確認OK，請刪除 registry、humpback 資料夾，在從頭重試看看

---

## 自動測試部署 (CI/CD)

這邊選擇使用 Drone 當做個人的 CI/CD 工具，因為客製化插件非常簡單，使用上很彈性。

1. 啟動 drone-server、drone-agent

    ```shell
    docker-compose up -d drone-server drone-agent
    ```

    ps. 啟動後，須等待伺服器啟動，可以執行以下指令確認是否開啟完畢
    ```shell
    docker-compose logs -f drone-server
    ```

2. 開啟 drone-server [[傳送門](http://drone-server:8000)]，並以 gitea 帳號登入

    ```md
    Username: demo
    Password: qwe123
    ```

3. 在 drone 頁面，開啟 demo/demo 的開關 (右側的小圈圈)

4. 檢查 Gitea Webhook 是否有新連結 [[傳送門](http://gitea:3000/demo/demo/settings/hooks)]

5. 開始進行專案的CI/CD

    5-1. 修改「.drone.yml」裡面的「host」，改為自己電腦的IP。
    5-2. 請任意做一個 commit，然後 push 上去。
    5-3. 查看[Drone頁面](http://drone-server:8000/demo/demo)，是否有運行 Pipeline 流程
    5-4. 結果如果有叉叉，純屬自然現象。

### Drone 錯誤處理

1. 如果頁面打不開，確認 drone-server 是否有啟動成功
2. 如果一直呈現鬧鐘狀態，沒有運行，確認 drone-agent 是否有啟動成功
3. 如果無法連線，確認 /etc/hosts 有沒有設定「前提準備」所指定的
4. 確認 drone-agent 有掛載到本機 /var/run/docker.sock
5. 如果上面4點都確認OK，請刪除 drone 資料夾，在從頭重試看看

---

## 日誌管理系統 (Log Management System)

這邊選擇使用 Elasticsearch(資料庫)、Fluentd(處理工具)、Kibana(介面) 組合，當做系統的日誌管理系統

1. 啟動 elasticsearch、fluentd、kibana

    ```shell
    docker-compose up -d elasticsearch fluentd kibana
    ```

    ps. 啟動後，須等待伺服器啟動，可以執行以下指令確認是否開啟完畢
    ```shell
    docker-compose logs -f fluentd
    ```

2. 先打開網頁伺服器，產生訊息

    2-1. 開啟 web

        ```shell
        docker-compose up -d web
        ```

    2-2 開啟頁面 [[傳送門](http://127.0.0.1:8888)]


3. 開啟 Kibana [[傳送門](http://kibana:5601)]，並初始化設定

    3-1. Index pattern：「logstash-*」 改為 「fluentd-*」 # ---> 這邊需要「步驟2-2」才有辦法填  

    3-2. Time Filter field name：打勾「Expand index pattern when searching [DEPRECATED]」

    3-3. 點擊「Create」

4. 瀏覽日誌紀錄 [[傳送門](http://kibana:5601/app/kibana#/discover)]

---

## 總結

---

## 參考資料

1. 版本控制管理 (Version Management)
    - [Gitea安裝](https://docs.gitea.io/zh-tw/install-with-docker/)

2. 容器管理系統 (Container Management)
    - [編譯與上傳映像檔](http://blog.poetries.top/handbook/html/%E6%9C%8D%E5%8A%A1%E7%AB%AF/docker.html#t13Docker%E7%A7%81%E6%9C%89%E4%BB%93%E5%BA%93%E6%90%AD%E5%BB%BA)
    - [建立容器倉庫](https://humpback.github.io/humpback/#/run-registry)
    - [Humpback 頁面建立](https://humpback.github.io/humpback/#/run-humpback-web)

3. 自動測試部署 (CI/CD)
    - [Drone安裝](http://docs.drone.io/installation/)
    - [Drone使用](http://docs.drone.io/getting-started/)

4. 日誌管理系統 (Log Management)
