FROM golang:1.12.5-alpine3.9 as base

ARG path=/go/src/github.com/crwnl3ss/micrograph/
ARG app_name=micrograph

WORKDIR ${path}
COPY . ${path}
RUN CGO_ENABLED=0 go test ./... && CGO_ENABLED=0 go build -o ${app_name}

FROM alpine3.9
LABEL maintainer="metalgearangel@gmail.com"
USER mg:1111
COPY --from=base ${path}{app_name}h /${app_name}
ENTRYPOINT ["./app"]

# udp metric receiver
EXPOSE 6666/udp
# web  interface, api, grafana datasource
EXPOSE 8000/tcp

# --volume myvolume:/path/in/container
VOLUME /

# HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://localhost/ || exit 1
