#!/bin/bash

value=$(cat <<EOF
mongodb:
  - address: "127.0.0.1:27017"
    username: ""
    password: ""
    auth_db: "admin"
EOF
)

consul kv put teddy/config/databases "$value"