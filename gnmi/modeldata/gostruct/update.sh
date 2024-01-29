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


rm -rf proto
generator -generate_fakeroot -output_file generated.go -package_name gostruct -exclude_modules ietf-interfaces -path public,yang -generate_simple_unions -compress_paths=false public/release/models/interfaces/openconfig-interfaces.yang public/release/models/openflow/openconfig-openflow.yang public/release/models/platform/openconfig-platform.yang public/release/models/system/openconfig-system.yang

echo dupa3

go run $GOPATH/pkg/mod/github.com/openconfig/ygot@v0.29.18/proto_generator/protogenerator.go \
    -generate_fakeroot \
    -base_import_path="proto" \
    -go_package_base="github.com/google/gnxi/gnmi/modeldata/gostruct/proto" \
    -ywrapper_path="proto/ywrapper" \
    -yext_path="proto/yext" \
    -path=public,yang \
    -output_dir=proto \
    -enum_package_name=enums -package_name=openconfig \
    -exclude_modules=ietf-interfaces \
    -compress_paths=false \
    public/release/models/interfaces/openconfig-interfaces.yang \
    public/release/models/openflow/openconfig-openflow.yang \
    public/release/models/platform/openconfig-platform.yang \
    public/release/models/system/openconfig-system.yang \

mkdir -p proto/ywrapper
cp $GOPATH/pkg/mod/github.com/openconfig/ygot@v0.29.18/proto/ywrapper/ywrapper.proto proto/ywrapper/ywrapper.proto
mkdir -p proto/yext
cp $GOPATH/pkg/mod/github.com/openconfig/ygot@v0.29.18/proto/yext/yext.proto proto/yext/yext.proto


echo dupa4
