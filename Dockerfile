FROM alpine:latest

ENV APPDIR /root/app
ENV PORT 8080

RUN mkdir -p $APPDIR

WORKDIR $APPDIR

COPY main $APPDIR/
COPY entrypoint.sh /

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]