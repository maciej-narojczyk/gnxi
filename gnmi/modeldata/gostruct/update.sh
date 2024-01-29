#!/bin/bash

echo dupa1

go get github.com/openconfig/ygot@latest
go get github.com/openconfig/ygot/ygot@latest

go get -u github.com/google/go-cmp/cmp
go get -u github.com/openconfig/gnmi/ctree
go get -u github.com/openconfig/gnmi/proto/gnmi
go get -u github.com/openconfig/gnmi/value
go get -u github.com/golang/glog
go get -u google.golang.org/protobuf/proto
go get -u github.com/kylelemons/godebug/pretty
go get -u github.com/openconfig/goyang/pkg/yang
go get -u google.golang.org/grpc

echo dupa2

rm -rf public yang
git clone https://github.com/openconfig/public.git
git clone https://github.com/YangModels/yang.git

go install github.com/openconfig/ygot/generator@latest

generator -generate_fakeroot -output_file generated.go -package_name gostruct -exclude_modules ietf-interfaces -path public,yang -generate_simple_unions -compress_paths=false public/release/models/interfaces/openconfig-interfaces.yang public/release/models/openflow/openconfig-openflow.yang public/release/models/platform/openconfig-platform.yang public/release/models/system/openconfig-system.yang

