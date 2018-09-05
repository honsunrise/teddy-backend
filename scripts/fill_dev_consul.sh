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
|
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

value=$(cat <<EOF
host: "127.0.0.1"
port: 567
username: "admin"
password: "admin"
EOF
)

consul kv put teddy/config/mail "$value"

value=$(cat <<EOF
|
    -----BEGIN PRIVATE KEY-----
    MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC2ywaJHRT8JLBh
    hJfYo0NMlgZ3CyN7cz5daIr4w5cUjd2f/wuxdqXz47S22jF+R/riTaT/+VRFJl8J
    YiS6LSX2NMZG59tcWFvjn5oFPSoGcq6a2fz1E0M/c2QzAMAPzCC2gsqGzO81u4V5
    +9L7VKrV9ynqXzv7nbuhGTO5kWoP0Xhqp/4mJ+BIvIgKO9UVPqXNzChicC6JSEoi
    N0v2oTUw4j8kDwmA8MNcDjxiajEQQUWA6mY3gAN2yJopxrqpNqObNET439gyFrv9
    0sYXsG8CYWyD6vFSMGDdj9mIYdX60l6dv7JZTJTJ108j6s8WZMqGwdVMD7g8rNgY
    tMeor51BAgMBAAECggEAKdWnZkQQpHBlKbxl4D/lTCbdzervsPY8JLajb7Gb5ylc
    upxtea0U6A+KMXsYbrVcluR8SdUvUzAn+gbLLwzcLk//vQSdcLIMPbkuT9qivp0K
    lwgi25gQAPqQyRd33WWzavHeFiHa8Wo8byGSNNE41AVgQ3KOUNTVt1YEP4knQ/0i
    d65koNbQ4UWEiVAbG5+mEbQzGv3Hqu0NPMV9FZDeVrYqzuiJkMNwCE2eprcSXZB2
    z8tj6q6dlwpPlsPFezNtL/thJ+caoOJYtb9sS+00o1o2vxvXrmA2zFTQtBBZRFuQ
    d0vpQYhtlQxqsmapK/tJhPqXTZoo/8KuYAzqcY7oUQKBgQDrpkKeDV2F6j8ihDme
    +ncrNhVr5xLDOWtPGp7VqCaiycrDS9cfNE74rlQjzmt0aZxnSGHKO93/BLpyDylr
    aYMmZ/TqLR44RCad0w7kcW9P6GBVcfQYP0jKkz+6KnCG7vWbUKb3csZErFgrV+ca
    9Cvo85tf2wyG6cWu1JG/VuzqBQKBgQDGlDdJemBmn056/2al7/6P2mSktAPRrDAQ
    WbJFaAe1ClqNJmrNp4j6hitxrvMp38F7/b+obJQaOtwWByJMNhOS/SIDRHSWsiq8
    UNL399jT22+rtUsHqkxAANNbNXCYRbDic8hqqDuJdHUF4VN4AolY6InfuCkFkDdK
    zxkS0za/DQKBgBzDAziFSxfwOlp9JwdHbMoiZMTxxDF9zaIvDpnnVyfhV1U06YHO
    gaEKrgxcwnLH/SYCCKWFXxgkPJl1Tknk6/QBFjyK2zhk4Q28WAH78mkfZLqpGPDo
    sHrBNDMFwQxHGEUnt+lV4es52d0YcoWwrbdWHG27r7C70bwAB/YBpxL9AoGAGdAU
    a7m7pDtbEUP3zOQoe/yQjpRT1sKCMO3n7Xu7XL4uzSBMS9VWSfJ83Tc3pp7OYNa4
    PiV3Dv3NtBNTUwLIgpfi/ve8DAa25WnAMrmF9uwUVQao7SMm7D7vOnD05OZSOu2A
    BNU/f/uiZpRGrFfwEJ2RpoIi2vVHKrNG+Bp6iFUCgYAE1DoqCg6loGS/ylIAPGuj
    VqQBglTfRVlMJWnV6kfbdG/2gflXMjULM/TH0VwrJGKsOAnxonpgagdR2eG9A+sL
    1h2RqlZNAVWUohf52k6uXPXSveCl17kYqWg6QXsxjXgTSx2UsFS05O15terQczsG
    QmeHEzkcBA5kvxzIQLusIQ==
    -----END PRIVATE KEY-----

EOF
)

consul kv put teddy/config/jwt_pkcs8 "$value"