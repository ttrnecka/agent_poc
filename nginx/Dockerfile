# build stage
FROM node:20-alpine AS npm
ARG IMPORT_ZSCALER_CERT=false

COPY zscaler/zscaler-root-ca.crt /tmp/zscaler-root-ca.crt

RUN if [ "$IMPORT_ZSCALER_CERT" = "true" ]; then \
      cp /tmp/zscaler-root-ca.crt /usr/local/share/ca-certificates/zscaler-root-ca.crt && \
      npm config set cafile /usr/local/share/ca-certificates/zscaler-root-ca.crt \
    else \
      echo "Skipping Zscaler cert install"; \
    fi

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci  --force --loglevel verbose
COPY frontend/ .
RUN npm run build

FROM nginx:alpine

COPY --from=npm /app/dist /usr/share/nginx/html/
