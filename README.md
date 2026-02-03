# av-control-api

`av-control-api` is a Go-based HTTP service that exposes a **REST API for controlling audiovisual (AV) devices by room**.

At a high level:

* Clients interact with the API using a `room` identifier.
* The service retrieves that room’s configuration from a **CouchDB-backed data service**.
* Based on the configured devices, the service uses a **driver registry** to communicate with physical AV hardware.
* The service supports both **reading current device state** and **applying desired state** across all devices in a room.

---

## Features

### Room-Based Control API

All primary endpoints are namespaced under:

```
/api/v1/room/:room
```

Supported operations include:

* Retrieve room configuration
  `GET /api/v1/room/:room`

* Retrieve current device state
  `GET /api/v1/room/:room/state`

* Apply desired device state
  `PUT /api/v1/room/:room/state`

* Retrieve device health information
  `GET /api/v1/room/:room/health`

* Retrieve room and device metadata
  `GET /api/v1/room/:room/info`

### Debug & Operational Endpoints

The service also exposes debug and operational endpoints:

* Health check
  `GET /debug/healthz`

* Runtime statistics
  `GET /debug/statz`

* Build and runtime info
  `GET /debug/infoz`

* Get current log level
  `GET /debug/logz`

* Set log level at runtime
  `GET /debug/logz/:level`
  (example levels: `debug`, `info`, `warn`, `error`)

---

## API Usage

### Setting Room State

**Endpoint**

```
PUT /api/v1/room/:room/state
```

**Request Body**

```json
{
  "devices": {
    "device-id-1": {
      "poweredOn": true,
      "blanked": false,
      "inputs": {
        "output-1": { "audioVideo": "input-2" }
      },
      "volumes": {
        "default": 40
      },
      "mutes": {
        "default": false
      }
    }
  }
}
```

**Response**

* `devices`: resulting device state as returned by drivers
* `errors`: per-device or per-field errors (partial success is supported)

Providing an `X-Request-ID` header is recommended for log correlation.

---

## Architecture Overview

The service is intentionally layered to keep responsibilities clear.

### Entry Point

```
cmd/av-control-api/
```

Responsibilities:

* Parse command-line flags
* Initialize logging
* Connect to the data service
* Load driver configuration
* Wire handlers and routes
* Start the HTTP server

### Handlers

```
handlers/
```

Handlers are intentionally thin and generally:

* Bind and validate request data
* Load room configuration via middleware
* Invoke state and driver logic
* Marshal and return HTTP responses

Middleware commonly includes:

* Request ID propagation
* Structured request logging
* Room configuration loading
* Host proxying

### State Orchestration

```
state/
```

The state layer coordinates:

* Reading current device state
* Applying requested state changes
* Aggregating per-device results and errors

### Drivers

```
drivers/
```

Drivers implement device-specific control logic and are registered by name.

Core driver implementations live under:

```
drivers/core/
```

These typically correspond to specific hardware vendors or protocols.

### Data Service

```
couch/
```

Room configuration is retrieved from CouchDB using the `kivik` CouchDB client library.

---

## Proxy Behavior

Room configuration may specify a **proxy host**.

If a request is received by a server that is not responsible for the room, the request may be transparently **proxied to the correct host** instead of being processed locally.

The `--host` flag identifies the current server instance for this purpose.

---

## Configuration

The service is configured primarily via command-line flags.

Common flags include:

* `--port` (`-P`)
  Port to listen on (default: `8080`)

* `--log-level` (`-L`)
  Log verbosity level

* `--host` (`-h`)
  Identifier for this server instance

* `--driver-config` (`-c`)
  Path to the driver configuration YAML file
  (default: `driver-config.yaml`)

* `--db-address`
  CouchDB server address

* `--db-username`
  CouchDB username

* `--db-password`
  CouchDB password

* `--db-insecure`
  Disable SSL for database connections

* `--cache-path`
  Optional filesystem cache location

---

## Running Locally

### Build

```bash
go build ./cmd/av-control-api
```

### Run

```bash
go run ./cmd/av-control-api \
  --port 8080 \
  --log-level info \
  --db-address "http://localhost:5984" \
  --db-username "admin" \
  --db-password "password" \
  --driver-config "./driver-config.yaml"
```

### Test

```bash
go test ./...
```

---

## Deployment

This repository includes deployment-related assets such as:

* `dockerfile` for container builds
* `terraform/` for infrastructure provisioning

Refer to those directories for environment-specific deployment details.

---

## Documentation

Additional documentation lives under:

```
docs/
```

These docs are structured for Antora and may reference older API paths. Verify against the current code before relying on them.

---

## License

See the `LICENSE` file for licensing information.