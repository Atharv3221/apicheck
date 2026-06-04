# apicheck

A fast, lightweight CLI tool written in Go for running and monitoring REST API health checks.
Define your APIs in a config file and get a clear pass/fail report in seconds.

---

## Demo

```
API NAME                              STATUS   STATUS_CODE   TIME
JSONPlaceholder - Get Single Post     PASS     200           275ms
JSONPlaceholder - Get All Posts       PASS     200           277ms
DummyJSON - Get Single Product        PASS     200           423ms
Postman Echo - GET Request            PASS     200           477ms
GoREST - Create User                  FAIL     401           541ms
DummyJSON - User Login                FAIL     400           679ms
HTTPBin - Delay 2 Seconds             PASS     200           3.603s
```

---

## Features

- Run health checks across multiple APIs in one command
- Supports **GET, POST, PUT, DELETE** HTTP methods
- Configure custom **headers** and **request bodies** per API
- Scoped configurations — use global or local config files
- Clean, tabular terminal output with status code and response time
- Concurrent execution for fast results

---

## Project Structure

```
.
├── cmd/
│   └── apicheck/
│       └── main.go              # Entry point
├── internal/
│   ├── commands/
│   │   ├── commands.go          # Root command setup
│   │   ├── add.go               # `add` command
│   │   ├── list.go              # `list` command
│   │   ├── remove.go            # `remove` command
│   │   ├── run.go               # `run` command
│   │   └── scope.go             # Scope (global/local) handling
│   ├── configparser/
│   │   ├── apiStructure.go      # API config data structures
│   │   ├── load.go              # Config loading
│   │   └── write.go             # Config writing
│   ├── constants/
│   │   └── methods/
│   │       └── method.go        # HTTP method constants
│   ├── execution/
│   │   ├── execute.go           # HTTP execution logic
│   │   └── terminal_output.go   # Table rendering
│   ├── initialization/
│   │   ├── default_config.yaml  # Default config template
│   │   ├── global.go            # Global config initialization
│   │   ├── initialize.go        # Init logic
│   │   └── local.go             # Local config initialization
│   └── const.go                 # Shared constants
├── go.mod
├── go.sum
└── README.md
```

---

## Installation

**Prerequisites:** Go 1.18+

```bash
git clone https://github.com/Atharv3221/apicheck.git
cd apicheck
go build -o apicheck ./cmd/apicheck
```

Optionally move the binary to your PATH:

```bash
sudo mv apicheck /usr/local/bin/
```

---

## Usage

### Initialize a config

```bash
apicheck init          # Initialize a local config in the current directory
apicheck init --global # Initialize a global config
```

### Add an API

```bash
apicheck add
```

You'll be prompted for:
- **Name** — a label for the API
- **URL** — the endpoint to hit
- **Method** — GET, POST, PUT, or DELETE
- **Headers** — optional, as `key: value` pairs
- **Body** — optional JSON body

### List configured APIs

```bash
apicheck list
```

### Run all checks

```bash
apicheck run
```

### Remove an API

```bash
apicheck remove <name>
```

---

## Config Format

APIs are stored in a YAML config file. Example entries:

```yaml
# POST with a JSON body and no headers
- name: Auth Login
  url: https://example.com/api/login
  method: POST
  body: '{"username": "admin", "password": "123"}'

# GET with custom headers
- name: xyz
  url: http://xyz.com
  method: GET
  headers:
    name: abc

```

---

## Status Reference

| STATUS | Meaning                                      |
|--------|----------------------------------------------|
| PASS   | Request succeeded with a 2xx status code     |
---

## License

See [LICENSE](./LICENSE) for details.
