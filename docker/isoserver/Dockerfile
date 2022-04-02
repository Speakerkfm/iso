FROM golang:1.16 as build

RUN mkdir /app

COPY . /app

RUN cd /app && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /bin/isoserver ./cmd/isoserver/main.go

FROM nginx

COPY --from=build /bin/isoserver /isoserver
COPY --from=build /app/iso_nginx.conf /iso_nginx.conf

# make dir for mounting
RUN mkdir /iso

CMD nginx -c /iso_nginx.conf && nginx -s reload && /isoserver -dir=iso