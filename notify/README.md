# Notify-Srv Service

This is the Notify-Srv service

Generated with

```
micro new github.com/zhsyourai/teddy-backend/notify-srv --namespace=go.micro --fqdn=com.teddy.srv.notify --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.teddy.srv.notify
- Type: srv
- Alias: notify-srv

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