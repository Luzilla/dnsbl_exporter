FROM alpine:3

ENV DNSBL_EXP_RBLS /etc/dnsbl-exporter/rbls.ini
ENV DNSBL_EXP_TARGETS /etc/dnsbl-exporter/targets.ini
ENV DNSBL_EXP_LISTEN :9211

# Add defaults
RUN mkdir -p /etc/dnsbl-exporter
ADD targets.ini "${DNSBL_EXP_TARGETS}"
ADD rbls.ini "${DNSBL_EXP_RBLS}"

ADD dnsbl-exporter /usr/local/bin

# TODO: create user account?

# This is for documentation
VOLUME /etc/dnsbl-exporter
EXPOSE 9211

ENTRYPOINT [ "/usr/local/bin/dnsbl-exporter" ]
