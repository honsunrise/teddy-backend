# Notify Service API

This is the Notify service API

Generated with

```
micro new github.com/zhsyourai/teddy-backend/api/notify --namespace=go.micro --fqdn=go.micro.api.notify --type=api
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.api.notify
- Type: api
- Alias: notify

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
./notify-api
```

Build a docker image
```
make docker
```