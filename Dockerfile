FROM alpine:3

ENV DNSBL_EXP_RESOLVER=ubound:53
ENV DNSBL_EXP_RBLS=/etc/dnsbl-exporter/rbls.ini
ENV DNSBL_EXP_TARGETS=/etc/dnsbl-exporter/targets.ini

COPY dnsbl-exporter /usr/bin/

RUN mkdir -p /etc/dnsbl-exporter
COPY rbls.ini /etc/dnsbl-exporter/
COPY targets.ini /etc/dnsbl-exporter/

EXPOSE 9211

ENTRYPOINT ["/usr/bin/dnsbl-exporter"]
