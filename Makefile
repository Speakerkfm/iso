build-server:
    env GOOS=linux GOARCH=amd64 go build -o bin/isoserver cmd/isoserver/main.go

build docker:
    docker build -t iso-server .

run docker:
    docker run --name my-iso-server -d -v example:/iso -p 82:82 iso-server