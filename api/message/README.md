# Message Service API

This is the Message service API

Generated with

```
micro new github.com/zhsyourai/teddy-backend/api/message --namespace=go.micro --fqdn=go.micro.api.message --type=api
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.api.message
- Type: api
- Alias: message

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
./message-api
```

Build a docker image
```
make docker
```