FROM golang:1.16 AS builder
WORKDIR /temp
COPY . /temp
RUN go env -w GOSUMDB=on && go env -w GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 go build -o thingmocker

FROM alpine AS link
COPY --from=builder /temp/thingmocker /usr/bin/thingmocker
RUN echo "http://mirrors.aliyun.com/alpine/v3.10/main/" > /etc/apk/repositories \
    && apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone