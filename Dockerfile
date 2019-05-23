FROM golang:1.12.5-alpine3.9 as base

WORKDIR /tmp/micrograph/

COPY . /tmp/micrograph/

RUN go test ./... && go build -o micrograph

FROM alpine3.9

COPY --from=base /tmp/micrograph/micrograph /opt/micrograph/app

EXPOSE 6666/tcp 6665/udp

