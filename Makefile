build:
    go build -o ./bin/isoserver cmd/isoserver/main.go

build docker:
    docker build -t iso-plugin -f ./docker/isoplugin/Dockerfile .
    docker build -t iso-server -f ./docker/isoserver/Dockerfile .

docker server:
    docker run --rm -v $(pwd)/example:/iso iso-plugin
    docker run --name my-iso-server -d -v $(pwd)/example:/iso -p 82:82 iso-server