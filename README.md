# fasthttp-stream-server

HTTP server for receiving large file uploads via `POST multipart/form-data` without loading the file into RAM.

## Task

The assignment required implementing streaming transfer of data exceeding available RAM through POST multipart/form-data using fasthttp. This project implements the receiving side — an HTTP server that handles such uploads.

## How it works

`fasthttp` is configured with `StreamRequestBody: true` and `DisablePreParseMultipartForm: true`, which prevents the framework from buffering the request body. Instead, the body is read as a stream directly from the network socket.

A fixed-size read buffer (64 KB) is reused across requests via `sync.Pool`. At any point in time, only one buffer per active connection is allocated — regardless of file size. A 10 GB file upload uses the same ~64 KB of application memory as a 1 KB upload.

## Key design decisions

- **No temporary files** — data is read and discarded (or could be piped to a writer); nothing is written to disk by the server itself
- **Bounded memory** — peak RSS grows with the number of concurrent connections, not with file size
- **Configurable concurrency and body size limit** — via YAML config
- **Graceful shutdown** — handles `SIGINT`/`SIGTERM`, waits for active connections to finish
- **Structured logging** — `slog` with request ID propagation and error classification (client disconnect vs server error)

## Configuration

Config file path is set via `CONFIG_PATH` environment variable.

```yaml
env: prod
settings:
  level: info
  idleTimeout: 30s
  writeTimeout: 120s
  concurrency: 256
  maxBodySize: 107374182400  # 100 GB
```

## Run locally

```bash
CONFIG_PATH=config/config.local.yaml go run ./cmd/fasthttp
```

## Run with Docker

```bash
docker build -t fasthttp-stream-server .
docker run -p 8080:8080 fasthttp-stream-server
```

## Test

Generate a large test file:

```bash
dd if=/dev/urandom of=data.bin bs=1M count=10240
```

Upload it:

```bash
curl -X POST -F "file=@./data.bin" http://localhost:8080/upload
```

## Memory usage

RSS during load test with 500 total uploads (40 concurrent) of a 10 GB file stays between 9–16 MB.

![Memory usage (RSS)](memory-usage-tests.png)
