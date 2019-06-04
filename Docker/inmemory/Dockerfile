FROM golang:1.12.5-alpine3.9 as base

ARG path=/go/src/github.com/crwnl3ss/micrograph/
ARG app_name=micrograph

RUN echo ${path}
WORKDIR ${path}
COPY . ${path}
RUN CGO_ENABLED=0 go test --count=1 -cover ./... && CGO_ENABLED=0 go build -o ${app_name}

FROM alpine:3.9
LABEL maintainer="metalgearangel@gmail.com"

ARG path=/go/src/github.com/crwnl3ss/micrograph/
ARG app_name=micrograph

RUN addgroup -g 1111 -S micrograph && adduser -u 1111 -S micrograph -G micrograph
USER micrograph:micrograph

COPY --from=base /go/src/github.com/crwnl3ss/micrograph/${app_name} /${app_name}
ENTRYPOINT ["./micrograph"]
CMD ["--help"]

# udp metric receiver
EXPOSE 6666/udp
# web  interface, api, grafana datasource
EXPOSE 8000/tcp

# HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://localhost/ || exit 1
