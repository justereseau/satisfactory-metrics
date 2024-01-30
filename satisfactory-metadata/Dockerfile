FROM golang:1.21-alpine as builder

ARG COMMIT_HASH

WORKDIR /go/src

COPY go.mod ./go.mod
COPY go.sum ./go.sum
COPY main.go ./main.go
COPY cache_worker.go ./cache_worker.go

RUN go mod download
RUN go build -o satisfactory-metadata -ldflags "-s -w" main.go

FROM scratch
LABEL maintainer="Sonic <sonic@djls.io>"
LABEL org.opencontainers.image.source=https://github.com/justereseau/satisfactory-metrics
LABEL org.opencontainers.image.description="Metadata syncer for Satisfactory InGame data."
LABEL org.opencontainers.image.licenses=WTFPL

COPY --from=builder /go/src/satisfactory-metadata /bin/satisfactory-metadata

ENTRYPOINT [ "/bin/satisfactory-metadata" ]
