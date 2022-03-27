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

# make dir for mounting
RUN mkdir /iso

CMD cd /iso/$PLUGIN_DIR && \
    protoc --go_out=. --go_opt=paths=source_relative **/*.proto && \
    go mod init $PLUGIN_MODULE_NAME && \
    go mod tidy && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o $PLUGIN_FILE_NAME $PLUGIN_GO_FILE_NAME && \
    cd .. && \
    mv ./$PLUGIN_DIR/$PLUGIN_FILE_NAME $PLUGIN_FILE_NAME && \
    rm -rf /iso/$PLUGIN_DIR
