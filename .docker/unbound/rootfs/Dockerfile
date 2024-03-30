FROM alpinelinux/unbound:latest

RUN apk add --update bind-tools drill \
    && chown -R unbound:unbound /etc/unbound

USER unbound
ADD --chown=unbound:unbound . /
ADD --chown=unbound:unbound --chmod=0644 https://www.internic.net/domain/named.cache /etc/unbound/root.hints

RUN ls -lah /etc/unbound

RUN /usr/sbin/unbound-checkconf

USER 100:101
EXPOSE 5353/udp

ENTRYPOINT ["unbound", "-dp"]
