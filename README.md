
**Healthcheck** is a go web-server that runs on *localhost:8792*.

## Installation

```
go run main.go
```

## Key functionality

By curling "/ping" developer can cumulatively check that all web-servers specified in [config](server/config.go) return 200 response by fetching servers uri "/ping".

## Information

I've modified original code to test additional functionality. This tool is licensed and can be used for personal purposes only.
