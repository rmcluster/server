FROM docker.io/golang:1.26.1@sha256:cdebbd553e5ed852386e9772e429031467fa44ca3a06735e6beb005d615e623d AS builder
WORKDIR /build-dir

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rmd-server .

FROM ghcr.io/ggml-org/llama.cpp:server@sha256:6cfa19f64b7a5ac67a783e7566bf57e036e120eb95d0d00670547cf96eb98780
COPY --from=builder /build-dir/rmd-server /usr/local/bin/rmd-server

ENTRYPOINT [ "rmd-server", "-host", "0.0.0.0", "-port", "4917" ]
EXPOSE 4917
