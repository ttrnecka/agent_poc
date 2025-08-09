FROM golang:1.23-alpine AS build
ARG IMPORT_ZSCALER_CERT=false

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt
RUN if [ "$IMPORT_ZSCALER_CERT" = "true" ]; then \
      cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt && \
      update-ca-certificates \
    else \
      echo "Skipping Zscaler cert install"; \
    fi


WORKDIR /app

COPY go.work* .
COPY collector collector
COPY ingestion ingestion
COPY logger logger
COPY webapi webapi
COPY common common

RUN go mod download
RUN go build -o /collector ./collector

# final stage
FROM alpine
WORKDIR /
COPY --from=build /collector /collector
ENTRYPOINT ["/collector"]