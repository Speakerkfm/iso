FROM golang:1.16

RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool
RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

# Get the source from GitHub
RUN go get google.golang.org/grpc
# Install protoc-gen-go
RUN go get github.com/golang/protobuf/protoc-gen-go

RUN mkdir /iso

CMD cd /iso/plugin && \
    protoc --go_out=. --go_opt=paths=source_relative **/*.proto && \
    go mod init iso_plugin && \
    go mod tidy && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin spec.go && \
    cd .. && \
    mv ./plugin/spec.so spec.so && \
    rm -rf ./plugin
