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

value=$(cat <<EOF
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
EOF
)

consul kv put teddy/config/casbin "$value"
