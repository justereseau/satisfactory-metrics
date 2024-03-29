FROM golang:1.22-alpine as builder

ARG COMMIT_HASH

WORKDIR /go/src

COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY main.go ./main.go
COPY exporter/ ./exporter

RUN go mod download
RUN go build -o satisfactory-exporter -ldflags "-s -w" main.go

FROM scratch
LABEL maintainer="Sonic <sonic@djls.io>"
LABEL org.opencontainers.image.source=https://github.com/justereseau/satisfactory-metrics
LABEL org.opencontainers.image.description="Prometheus exporter for Satisfactory InGame data."
LABEL org.opencontainers.image.licenses=WTFPL

COPY --from=builder /go/src/satisfactory-exporter /bin/satisfactory-exporter

ENTRYPOINT [ "/bin/satisfactory-exporter" ]
