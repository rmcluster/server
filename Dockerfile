FROM docker.io/golang:1.26.2@sha256:b54cbf583d390341599d7bcbc062425c081105cc5ef6d170ced98ef9d047c716 AS builder
WORKDIR /build-dir

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rmd-server .

FROM ghcr.io/rmcluster/llama.cpp-rpc:server@sha256:0fcff462704e7e1fa5d692109934ca79927d7811f64fe6b070dfd0fd82e758a9
COPY --from=builder /build-dir/rmd-server /usr/local/bin/rmd-server

# llama.cpp's docker image puts the executables in /app
ENV PATH=/app:$PATH
ENTRYPOINT [ "rmd-server", "-host", "0.0.0.0", "-port", "4917" ]
EXPOSE 4917
