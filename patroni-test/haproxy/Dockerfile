#FROM haproxy:2.6-alpine
#COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg
FROM ubuntu:latest

RUN apt-get update \
    && apt-get install -y haproxy rsyslog \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir /run/haproxy

COPY haproxy.cfg /etc/haproxy/haproxy.cfg

CMD haproxy -f /etc/haproxy/haproxy.cfg && tail -F /var/log/haproxy.log