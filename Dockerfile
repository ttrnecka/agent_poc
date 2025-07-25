# build stage
FROM node:20-alpine AS npm

COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt

RUN npm config set cafile /usr/local/share/ca-certificates/zscaler-root-ca.crt
WORKDIR /app
COPY agent_poc/package*.json ./
RUN npm ci  --force --loglevel verbose
COPY agent_poc/ .
RUN npm run build

FROM golang:1.23-alpine AS build

COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt

ENV CGO_ENABLED=0

RUN update-ca-certificates
RUN apk update && apk add bash ca-certificates dos2unix && update-ca-certificates 2>/dev/null

WORKDIR /app

COPY . .
COPY --from=npm /app/dist /app/agent_poc/dist

RUN go mod download
RUN go build -o /agent_poc

WORKDIR /app/policies/brocade

RUN chmod +x build.sh
RUN dos2unix ./build.sh
RUN bash ./build.sh 1.0.0
RUN bash ./build.sh 1.0.1
RUN bash ./build.sh 1.0.2
RUN bash ./build.sh 1.0.3

# final stage
FROM alpine 
WORKDIR /
COPY --from=build /agent_poc /agent_poc
COPY --from=build /app/data /data
COPY --from=build /app/policies/brocade/ /data/policies
EXPOSE 8888
ENTRYPOINT ["/agent_poc"]