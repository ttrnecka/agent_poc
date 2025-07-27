# Podman compose setup on Ubuntu
https://linuxopsys.com/getting-started-with-podman-compose

# To dev GUI

## Start Vite
cd frontend
npm run dev

## Start webapi
cd webapi
go run .

# To build npm locally 
npm run build

# Zscaller
Use 

```set IMPORT_ZSCALER_CERT=true```

before starting the docker compose up

# To deploy GUI + backend

docker compose build --no-cache nginx webapi 

docker image prune -f

docker compose up -d nginx webapi

# To dev collector
cd collector
go run . 
go run . --source collector1|collector2

# Collector deploy

docker compose build --no-cache collector1

docker image prune -f

docker compose up -d collector1 collector2