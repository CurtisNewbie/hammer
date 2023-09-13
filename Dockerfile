FROM golang:1.18-alpine as build
LABEL author="Yongjie Zhuang"
LABEL descrption="Hammer - Image processing service"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache vips \
    vips-dev \
    tzdata \
    gcc \
    libc-dev \
    glib-dev \
    pkgconfig \
    glib

WORKDIR /go/src/build/

# for golang env
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# dependencies
COPY go.mod .
COPY go.sum .

RUN go mod download

# build executable
COPY . .
RUN go build -o main

# ---------------------------------------------

FROM alpine:3.17

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache vips \
    vips-dev

WORKDIR /usr/src/
COPY --from=build /go/src/build/main ./app_hammer
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Shanghai

CMD ["./app_hammer", "configFile=/usr/src/config/app-conf-prod.yml"]
