FROM golang:1.19.5-alpine as builder
RUN apk add --update alpine-sdk

ENV APP_HOME /ecies-encrypt
RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME
ADD . $APP_HOME

RUN cd $APP_HOME && \
  go mod download && \
  go build -o ecies-encrypt ecies-encrypt.go && \
  go run keys.go

FROM kong:3.1.1-alpine

USER root
RUN mkdir -p /ecies-encrypt
COPY --from=builder /ecies-encrypt/ecies-encrypt /usr/local/bin/
COPY --from=builder /ecies-encrypt/ecies.pk.key /ecies-encrypt/
USER kong

ENTRYPOINT ["/docker-entrypoint.sh"]
EXPOSE 8000 8443 8001 8444
STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health
CMD ["kong", "docker-start"]