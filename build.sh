#!/usr/bin/env bash

push uaa
make docker
pop

push content
make docker
pop

push captcha
make docker
pop

push message
make docker
pop

push api/uaa
make docker
pop

push api/base
make docker
pop

push api/content
make docker
pop

push api/message
make docker
pop