curl -X POST \
    http://localhost:8080/limit \
    -H 'content-type: application/json' \
    -d '{"ip":"127.0.0.3","func":"x.y.z"}'
