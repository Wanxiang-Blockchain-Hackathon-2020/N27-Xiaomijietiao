FROM golang:1.8
MAINTAINER liangzhaohao "liang.zhaohao@qq.com"

RUN mkdir -p /data \
    && mkdir -p /opt/weyom \
    && mkdir -p /opt/weyom/bin \
    && mkdir -p /opt/weyom/config


VOLUME /data/logs

COPY ioubackend /opt/weyom/bin/ioubackend
ADD ./start.sh /opt/weyom/
RUN chmod +x /opt/weyom/start.sh
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

WORKDIR /opt/weyom

CMD ["/opt/weyom/start.sh"]


