# Captcha-Srv Service

This is the Captcha-Srv service

Generated with

```
micro new github.com/zhsyourai/teddy-backend/captcha --namespace=go.micro --fqdn=com.teddy.srv.captcha --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: com.teddy.srv.captcha
- Type: srv
- Alias: captcha

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
./captcha-srv
```

Build a docker image
```
make docker
```