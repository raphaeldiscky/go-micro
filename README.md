# microservices-go

## Run the project

```bash
cd project
make up_build
make start
```

## Stop the project

```bash
make stop
```

## Generate proto files (example for logs.proto)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"

cd logger-service/logs
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    logs.proto

cd logger-service
go get google.golang.org/grpc
```

## Build and push all microservices

```bash
cd logger-service
docker build -f logger-service.dockerfile -t rdisckyzp/logger-service:1.0.0 .
docker push rdisckyzp/logger-service:1.0.0
```

## Deploy the project to swarm

```bash
docker stack deploy -c swarm.yaml go-micro
```

## Remove the project from swarm

```bash
docker stack rm go-micro
```
