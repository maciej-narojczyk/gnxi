
project(App)

cmake_minimum_required(VERSION 3.23)

find_package(Protobuf REQUIRED)

set( go_path "$ENV{GOPATH}" )

set( import_dirs 
    ${go_path}/src
)

file(GLOB_RECURSE proto_files RELATIVE ${CMAKE_CURRENT_SOURCE_DIR} "proto/*.proto")

# message("proto_files: ${proto_files}")

add_library(proto_yang 
    # proto/openconfig/openconfig_acl/openconfig_acl.proto
    ${proto_files}
    # proto/yext/yext.proto
    # proto/ywrapper/ywrapper.proto
    )

target_include_directories(proto_yang PUBLIC ${CMAKE_CURRENT_BINARY_DIR})
target_link_libraries(proto_yang PUBLIC ${Protobuf_LIBRARIES})

protobuf_generate(
    TARGET proto_yang 
    LANGUAGE cpp
    # IMPORT_DIRS ${import_dirs}
    OUT_VAR tmp
    )


# message( "tmp: ${tmp}")

add_executable(app main.cpp)

target_link_libraries(app proto_yang )

