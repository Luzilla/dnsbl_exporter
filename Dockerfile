FROM debian:stable-slim

# nobody / nogroup
ARG DNSBL_USER=65534
ARG DNSBL_GROUP=65534

ENV DNSBL_EXP_RESOLVER=ubound:53
ENV DNSBL_EXP_RBLS=/etc/dnsbl-exporter/rbls.ini
ENV DNSBL_EXP_TARGETS=/etc/dnsbl-exporter/targets.ini

COPY dnsbl-exporter /usr/bin/

RUN mkdir -p /etc/dnsbl-exporter
COPY rbls.ini /etc/dnsbl-exporter/
COPY targets.ini /etc/dnsbl-exporter/

RUN chown -R $DNSBL_USER:$DNSBL_GROUP /etc/dnsbl-exporter
USER $DNSBL_USER:$DNSBL_GROUP

EXPOSE 9211

ENTRYPOINT ["/usr/bin/dnsbl-exporter"]
