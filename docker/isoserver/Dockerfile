FROM golang:1.16 as build

RUN mkdir /app

COPY . /app

RUN cd /app && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/isoserver ./cmd/isoserver/main.go

FROM golang:1.16

COPY --from=build /bin/isoserver /isoserver

RUN mkdir /iso

ENTRYPOINT ["/isoserver"]
CMD ["-dir=/iso"]