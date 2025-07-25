FROM golang:1.23-alpine AS build

COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN update-ca-certificates

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /ingest_service

# final stage
FROM alpine
WORKDIR /
COPY --from=build /ingest_service /ingest_service
ENTRYPOINT ["/ingest_service"]