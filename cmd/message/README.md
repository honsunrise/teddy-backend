# Message-Srv Service

This is the Message-Srv service

Generated with

```
micro new github.com/zhsyourai/teddy-backend/message --namespace=go.micro --fqdn=com.teddy.srv.message --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.teddy.srv.message
- Type: srv
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
./message-srv
```

Build a docker image
```
make docker
```