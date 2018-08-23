# UAA-Srv Service

This is the Uaa-Srv service

Generated with

```
micro new github.com/zhsyourai/teddy-backend/uaa-srv --namespace=go.micro --fqdn=com.teddy.srv.uaa --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.teddy.srv.uaa
- Type: srv
- Alias: uaa-srv

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./micro-srv-srv
```

Build a docker image
```
make docker
```