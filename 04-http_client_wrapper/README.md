# http-client-wrapper

## Purpose
- Understand:
    - What `http.Client` actually do.
    - How to wrap an existing client.
    - How requests flow through code.
    - How to add cross-cutting behavior (logging, headers)

## Flow
```
main() -> WrapperClient -> http.Client -> Internet
```

## Project Structure
```
mini-http-wrapper/
├── go.mod
└── main.go
```
