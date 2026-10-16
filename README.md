# apicheck

A fast CLI in Go for REST API health checks. Define your APIs in a YAML file and get a 
pass/fail report in one command.

## Demo

```
$ apicheck run
API NAME                      STATUS   STATUS_CODE   TIME
DummyJSON-GetProduct          PASS     200           103ms
JSONPlaceholder-GetPost       PASS     200           155ms
JSONPlaceholder-GetAllPosts   PASS     200           156ms
PostmanEcho-PUT               PASS     200           322ms
PostmanEcho-DELETE            PASS     200           329ms
JSONPlaceholder-CreatePost    PASS     201           397ms
DummyJSON-UserLogin           PASS     200           431ms
PostmanEcho-GET               PASS     200           732ms
HTTPBin-Get                   PASS     200           926ms
HTTPBin-Delay2s               PASS     200           2.933s
```

## Features

- GET, POST, PUT and DELETE, with custom headers and JSON bodies
- All checks run concurrently, so a run takes about as long as the slowest API (1,000 local checks in ~185 ms)
- Global (`~/.apicheck.yaml`) or per-project (`.apicheck.yaml`) config

## Install

Requires Go 1.25+.

```bash
git clone https://github.com/Atharv3221/apicheck.git
cd apicheck
go build -o apicheck ./cmd/apicheck
```

## Usage

```bash
apicheck init                # create a local .apicheck.yaml
apicheck add                 # add an API interactively
apicheck list                # list configured APIs
apicheck run                 # run all checks
apicheck run -n <name>...    # run only the named checks
apicheck remove <name>...    # remove APIs
apicheck scope               # show which config is in use
```

The global config is created on first run. A `.apicheck.yaml` in the current directory takes priority.

## Config

```yaml
apis:
  - name: Auth Login
    url: https://example.com/api/login
    method: POST
    requestBody: '{"username": "admin", "password": "123"}'

  - name: Profile
    url: https://example.com/api/me
    method: GET
    header:
      Authorization: Bearer <token>
```

If `requestBody` is set and there is no `Content-Type` header, `application/json` is sent.

## Status

| Status | Meaning |
| --- | --- |
| PASS | 2xx response |
| FAIL | Non-2xx response or connection error |

## License

MIT, see [LICENSE](LICENSE).
