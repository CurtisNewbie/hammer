FROM alpine:3.17

LABEL author="Yongjie Zhuang"
LABEL descrption="Hammer - Thumbnail Generation Service"

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add tzdata \
    ffmpeg

WORKDIR /usr/src/

# binary is pre-compiled
COPY hammer_build ./app_hammer

ENV TZ=Asia/Shanghai

CMD ["./app_hammer", "configFile=/usr/src/config/conf.yml"]
