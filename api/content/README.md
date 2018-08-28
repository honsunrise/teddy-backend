# Message Service API

This is the Message service API

Generated with

```
micro new github.com/zhsyourai/teddy-backend/api/content --namespace=go.micro --fqdn=go.micro.api.content --type=api
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.api.content
- Type: api
- Alias: content

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
./content-api
```

Build a docker image
```
make docker
```