FROM golang:1.23-alpine AS build
ARG IMPORT_ZSCALER_CERT=false

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt

ENV CGO_ENABLED=0

RUN if [ "$IMPORT_ZSCALER_CERT" = "true" ]; then \
      cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt && \
      update-ca-certificates \
    else \
      echo "Skipping Zscaler cert install"; \
    fi

RUN apk update && apk add bash ca-certificates dos2unix && update-ca-certificates 2>/dev/null

WORKDIR /app

COPY go.work* .
COPY collector collector
COPY ingestion ingestion
COPY logger logger
COPY webapi webapi
RUN go mod download
RUN go build -o /webapi ./webapi

WORKDIR /app/webapi/policies/brocade

RUN chmod +x build.sh
RUN dos2unix ./build.sh
RUN bash ./build.sh 1.0.0
RUN bash ./build.sh 1.0.1
RUN bash ./build.sh 1.0.2
RUN bash ./build.sh 1.0.3
RUN bash ./build.sh 1.0.4

# final stage
FROM alpine 
WORKDIR /
COPY --from=build /webapi /webapi
COPY --from=build /app/webapi/data /data
COPY --from=build /app/webapi/policies/brocade/ /data/policies
EXPOSE 8888
ENTRYPOINT ["/webapi"]