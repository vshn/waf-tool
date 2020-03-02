FROM docker.io/library/alpine:3.11 as runtime

ENTRYPOINT ["waf-tuning-tool"]

RUN \
    apk add --no-cache curl bash

COPY waf-tuning-tool /usr/bin/image-cleanup
USER 1000:0
