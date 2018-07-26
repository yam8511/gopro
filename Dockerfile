# FROM gitea/gitea:latest
FROM yam8511/govendor:alpine

RUN apk --no-cache add \
    bash \
    ca-certificates \
    curl \
    gettext \
    git \
    linux-pam \
    openssh \
    s6 \
    sqlite \
    su-exec \
    tzdata

RUN addgroup \
    -S -g 1000 \
    git && \
    adduser \
    -S -H -D \
    -h /data/git \
    -s /bin/bash \
    -u 1000 \
    -G git \
    git && \
    echo "git:$(dd if=/dev/urandom bs=24 count=1 status=none | base64)" | chpasswd

ENV USER git
ENV GITEA_CUSTOM /data/gitea

COPY ./data /data
COPY ./main.go ./main.go
COPY ./gitea ./gitea
# RUN ls -al
RUN go build -o app ./main.go
# RUN ls -al

CMD [ "./app" ]
