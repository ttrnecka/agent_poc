# To dev gui
cd agent_poc
npm run dev
cd ..
go run .

# To build npm locally 

npm run build


# To deploy gui + backend

docker build -t agent_poc:test . 

docker image prune -f

docker network create poc
docker run --rm --name server --network poc -p "8888:8888" agent_poc:test

# To dev collector
cd collector
go run . 
go run . --source collector1|collector2

# To build collector

## local docker based build
cd collector

DOCKER_BUILDKIT=1 docker build --output type=local,dest=./out .

Run export DOCKER_BUILDKIT=1

Run docker build --target bin --output bin/ .


## docker image build and run on the same nw as server
docker build -t collector:test .

docker image prune -f

docker run --rm --network poc collector:test --addr server:8888

docker run --rm --network poc collector:test --addr server:8888 --source collector2