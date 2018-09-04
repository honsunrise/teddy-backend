# Base Service API

This is the Base service API

Generated with

```
micro new github.com/zhsyourai/teddy-backend/api/base --namespace=go.micro --fqdn=go.micro.api.base --type=api
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.api.base
- Type: api
- Alias: base

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
./base-api
```

Build a docker image
```
make docker
```