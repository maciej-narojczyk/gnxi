
go run ./gnmi_target.go \
    -bind_address :9339 \
    -config openconfig-openflow.json \
    -key ../certs/target.key \
    -cert ../certs/target.crt \
    -ca ../certs/ca.crt
