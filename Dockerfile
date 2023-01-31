FROM golang:1.19.5-alpine as builder

ENV APP_HOME /ecies-encrypt
RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME
ADD . $APP_HOME

RUN cd $APP_HOME && \
  go mod download && \
  go build -o ecies-encrypt ecies-encrypt.go

FROM kong:3.1.1-alpine

USER root
COPY --from=builder /ecies-encrypt/ecies-encrypt /usr/local/bin/

USER kong
ENTRYPOINT ["/docker-entrypoint.sh"]
EXPOSE 8000 8443 8001 8444
STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health
CMD ["kong", "docker-start"]