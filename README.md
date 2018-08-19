# Telegram Bot

1. 先登入註冊Telegram的帳號

2. 建立機器人 [教學文章](https://hackmd.io/QglPYeBDTi6pvr7PStWJ-Q#%E7%94%B3%E8%AB%8B-Telegram-Bot)
只須完成「申請Telegram Bot」的步驟即可

3. 先在專案目錄底下新增 .env 檔案，並設定自己的BOT_TOKEN
```shell
# echo BOT_TOKEN=[自己的機器人Token] >> .env
echo BOT_TOKEN=4474ASF....... >> .env
```

4. 啟動服務
```shell
docker-compose up --build
```

5. 在自己的機器人裡面打文字或指令
```telegram
Hi~ Bot
> Hi~ Bot , 接收訊息 By Zuolar

/hello
> Hello~ Zuolar

/myid
> Your ID is 1234567

/cmd [....]
> You enter [....]

/start
> Hello (出現一組按鈕可以點擊)

/cancel
> Reset (關閉按鈕)
```

6. 或是傳送檔案或圖片給自己的機器人
