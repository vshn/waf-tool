FROM docker.io/library/alpine:3.13 as runtime

ENTRYPOINT ["waf-tool"]

RUN \
    apk add --no-cache curl bash

COPY waf-tool /usr/bin/waf-tool
USER 1000:0
