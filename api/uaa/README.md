# Uaa Service API

This is the Uaa service API

Generated with

```
micro new github.com/zhsyourai/teddy-backend/api/uaa --namespace=go.micro --fqdn=com.teddy.api.uaa --type=api
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.teddy.api.uaa
- Type: api
- Alias: uaa

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
./uaa-api
```

Build a docker image
```
make docker
```