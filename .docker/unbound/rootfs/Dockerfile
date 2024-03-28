FROM alpinelinux/unbound:latest

RUN apk add --update bind-tools drill
ADD https://www.internic.net/domain/named.cache /etc/unbound/root.hints
ADD . /

RUN unbound-checkconf

EXPOSE 5353/udp
