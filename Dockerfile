FROM yam8511/govendor:alpine

# 設置環境變數
ENV SITE=web
# 以下為 PostgreSQL 的設定，可打開註解進行設定
# ENV POSTGRES_HOST=???
# ENV POSTGRES_PORT=5432
# ENV POSTGRES_DB=???
# ENV POSTGRES_USER=???
# ENV POSTGRES_PASSWORD=???

# 編譯應用程式
COPY . /go/src/app
WORKDIR /go/src/app
RUN govendor init
RUN govendor add +external
RUN govendor fetch -v +missing
RUN go build -o app

# 指定使用者，避免以root身份執行
# 如果是Alpine OS的映像檔，打開以下註解
RUN adduser -D zuolar
USER zuolar
# 如果是Ubuntu OS的映像檔，打開以下註解
# RUN useradd -m zuolar
# USER zuolar

# 開始運作程序
CMD [ "./app" ]
