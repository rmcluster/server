FROM docker.io/golang:1.27rc2@sha256:2317c8e806fe884a0d5f3d80d596dcd1369a4c211b7eba8d507718245e9e5831 AS builder
WORKDIR /build-dir

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rmd-server .

FROM ghcr.io/rmcluster/llama.cpp-rpc:server@sha256:45ca3a255f43bbead13202635ee97f37e39ac96605672319f57f617f9a6e2264
COPY --from=builder /build-dir/rmd-server /usr/local/bin/rmd-server

# llama.cpp's docker image puts the executables in /app
ENV PATH=/app:$PATH
ENV RMD_METADATA_DB_PATH=/var/lib/rmd/metadata.db
ENV RMD_MODEL_STORAGE_DIR=/var/lib/rmd/models
RUN mkdir -p /var/lib/rmd
RUN mkdir -p /var/lib/rmd/models
VOLUME ["/var/lib/rmd"]
ENTRYPOINT [ "rmd-server", "-host", "0.0.0.0", "-port", "4917" ]
EXPOSE 4917
