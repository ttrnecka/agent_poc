# build stage
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /agent_poc

# final stage
FROM alpine
WORKDIR /
COPY --from=build /agent_poc /agent_poc
COPY --from=build /app/data /data
EXPOSE 8888
ENTRYPOINT ["/agent_poc"]