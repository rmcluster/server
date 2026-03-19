FROM docker.io/golang:1.26.1@sha256:cdebbd553e5ed852386e9772e429031467fa44ca3a06735e6beb005d615e623d AS builder
WORKDIR /rama-swap

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rama-swap .

FROM quay.io/ramalama/ramalama:0.18.0@sha256:1d63a3e761113810c1efdf3b24195af22e8a4c6c095684d00755aebac7a01bc2
COPY --from=builder /rama-swap/rama-swap /usr/local/bin/rama-swap

ENTRYPOINT [ "env", "RAMALAMA_STORE=/app/store", "rama-swap", "-ramalama", "ramalama", "--nocontainer", ";", "-host", "0.0.0.0", "-port", "4917" ]
EXPOSE 4917
