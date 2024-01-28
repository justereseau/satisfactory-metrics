FROM golang:1.21-alpine as builder

ARG COMMIT_HASH

WORKDIR /go/src

COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY main.go ./main.go
COPY exporter/ ./exporter

RUN go mod download
RUN go build -o satisfactory-metrics -ldflags "-s -w" main.go

FROM scratch
LABEL maintainer="Sonic <sonic@djls.io>"

COPY --from=builder /go/src/satisfactory-metrics /bin/satisfactory-metrics

ENTRYPOINT [ "/bin/satisfactory-metrics" ]
