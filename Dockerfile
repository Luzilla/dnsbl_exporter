FROM scratch

# nobody / nogroup
ARG DNSBL_USER=65534
ARG DNSBL_GROUP=65534

USER $DNSBL_USER:$DNSBL_GROUP

ENV DNSBL_EXP_RESOLVER=ubound:53
ENV DNSBL_EXP_RBLS=/rbls.ini
ENV DNSBL_EXP_TARGETS=/targets.ini

COPY dnsbl-exporter rbls.ini targets.ini /

EXPOSE 9211

ENTRYPOINT ["/dnsbl-exporter"]
