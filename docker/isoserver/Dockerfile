FROM golang:1.16 as build

RUN mkdir /app

COPY . /app

RUN cd /app && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/isoserver ./cmd/isoserver/main.go

FROM nginx

COPY --from=build /bin/isoserver /isoserver

# make dir for mounting
RUN mkdir /iso

# apply nginx config and start server
CMD nginx -c /iso/$REVERSE_PROXY_CONFIG_FILE_NAME && \
    nginx -s reload && \
    /isoserver -dir=iso