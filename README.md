# Reused Mobile Devices

Server and Client for running an automatically managed llama RPC cluster.

## Run Server

```sh
go run .
```

## Run Client

The client will announce itself to the tracker and run the rpc server.
Replace `/path/to/rpc-server` with the path to the rpc server binary.
Replace `127.0.0.1:4917` with the ip and port of the tracker.

```sh
go run ./cmd/linux-client/ -cmd /path/to/rpc-server -tracker 127.0.0.1:4917 -- -c
```
