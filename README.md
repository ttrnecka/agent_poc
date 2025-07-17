To dev gui
cd agent_poc
npm run dev
cd ..
go run .

To deploy gui + backend

npm run build

docker build -t agent_poc:test .

docker image prune -f

docker run --rm -p "8888:8888" agent_poc:test

To dev collector
cd collector
go run . 
go run . --source collector1|collector2

