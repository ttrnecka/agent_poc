# build stage
FROM node:20-alpine AS npm

# COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt

# RUN npm config set cafile /usr/local/share/ca-certificates/zscaler-root-ca.crt
WORKDIR /app
COPY agent_poc/package*.json ./
RUN npm ci  --force --loglevel verbose
COPY agent_poc/ .
RUN npm run build

FROM golang:1.23-alpine AS build


COPY zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt
RUN update-ca-certificates

WORKDIR /app

COPY . .
COPY --from=npm /app/dist /app/agent_poc/dist

RUN go mod download
RUN go build -o /agent_poc

# final stage
FROM alpine 
WORKDIR /
COPY --from=build /agent_poc /agent_poc
COPY --from=build /app/data /data
EXPOSE 8888
ENTRYPOINT ["/agent_poc"]