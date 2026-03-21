FROM docker.io/golang:1.26.1@sha256:cdebbd553e5ed852386e9772e429031467fa44ca3a06735e6beb005d615e623d AS builder
WORKDIR /rama-swap

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rama-swap .

FROM ghcr.io/ggml-org/llama.cpp:server@sha256:ab3835d97af58fd1f553abf4e84ffd166a59b4ee24cdfe7733cb44d5cf8ca1c5
COPY --from=builder /rama-swap/rama-swap /usr/local/bin/rama-swap

ENTRYPOINT [ "env", "RAMALAMA_STORE=/app/store", "rama-swap", "-ramalama", "ramalama", "--nocontainer", ";", "-host", "0.0.0.0", "-port", "4917" ]
EXPOSE 4917
