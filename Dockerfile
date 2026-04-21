FROM docker.io/golang:1.26.2@sha256:5f3787b7f902c07c7ec4f3aa91a301a3eda8133aa32661a3b3a3a86ab3a68a36 AS builder
WORKDIR /build-dir

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rmd-server .

FROM ghcr.io/rmcluster/llama.cpp-rpc:server@sha256:8725f50f673998ea13d1b08210d9a81d890768a59efea8205402ff7a021b9a51
COPY --from=builder /build-dir/rmd-server /usr/local/bin/rmd-server

# llama.cpp's docker image puts the executables in /app
ENV PATH=/app:$PATH
ENTRYPOINT [ "rmd-server", "-host", "0.0.0.0", "-port", "4917" ]
EXPOSE 4917
