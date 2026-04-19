FROM docker.io/golang:1.26.2@sha256:5f3787b7f902c07c7ec4f3aa91a301a3eda8133aa32661a3b3a3a86ab3a68a36 AS builder
WORKDIR /build-dir

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rmd-server .

FROM ghcr.io/project-panzerschreck/llama.cpp-rpc:server@sha256:296159e5b568ababefc8b317e9e6edb2f9f20986109fb9f0fa01e6b4e6bd70a9
COPY --from=builder /build-dir/rmd-server /usr/local/bin/rmd-server

# llama.cpp's docker image puts the executables in /app
ENV PATH=/app:$PATH
ENTRYPOINT [ "rmd-server", "-host", "0.0.0.0", "-port", "4917" ]
EXPOSE 4917
