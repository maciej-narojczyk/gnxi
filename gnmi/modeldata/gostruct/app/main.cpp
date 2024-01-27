
#include <iostream>
#include <fstream>
#include <cassert>
#include <string>

#include <google/protobuf/util/json_util.h>
#include "proto/openconfig/openconfig.pb.h"

int main()
{
    std::cout << "start" << std::endl;
    std::ifstream in("./config.in.json",std::ios::in );
    assert(in.is_open());
    std::string str((std::istreambuf_iterator<char>(in)),
                 std::istreambuf_iterator<char>());

    
    ::openconfig::Device config;

    ::google::protobuf::util::JsonParseOptions opt;
    opt.ignore_unknown_fields = false;
    opt.case_insensitive_enum_parsing = true;

    auto ret = ::google::protobuf::util::JsonStringToMessage(str,&config,opt);
    if ( not ret.ok())
    {
        std::cout << "error: " << ret.ToString() << std::endl;
    }

    std::string out;
    auto ret2 = ::google::protobuf::util::MessageToJsonString(config,&out);
    if ( not ret.ok())
    {
        std::cout << "error2: " << ret2.ToString() << std::endl;
    }

    std::ofstream o("./config.out.json", std::ios::out);
    o << out;


    std::cout << "end" << std::endl;
    return 0;
}
