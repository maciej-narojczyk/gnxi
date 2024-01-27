/* Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Binary gnmi_target implements a gNMI Target with in-memory configuration and telemetry.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"reflect"

	log "github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/google/gnxi/gnmi"
	"github.com/google/gnxi/gnmi/modeldata"
	"github.com/google/gnxi/gnmi/modeldata/gostruct"

	"github.com/google/gnxi/utils/credentials"

	pb "github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/proto/ywrapper"
	"github.com/openconfig/ygot/ygot"

	"github.com/google/gnxi/gnmi/modeldata/gostruct/proto/openconfig"
	"github.com/google/gnxi/gnmi/modeldata/gostruct/proto/openconfig/openconfig_interfaces"

	"google.golang.org/protobuf/encoding/protojson"
)

var (
	bindAddr   = flag.String("bind_address", ":9339", "Bind to address:port or just :port")
	configFile = flag.String("config", "", "IETF JSON file for target startup config")
)

type server struct {
	*gnmi.Server
}

func parseProto() {

	log.Info("parseProto")

	var str ywrapper.StringValue
	str.Value = "qwe"
	var conf openconfig_interfaces.Interfaces_Interface_Config
	conf.Name = &str

	var ifc openconfig_interfaces.Interfaces_Interface
	ifc.Config = &conf

	var ifk openconfig_interfaces.Interfaces_InterfaceKey
	ifk.Interface = &ifc
	ifk.Name = "asd"

	var ifc2 openconfig_interfaces.Interfaces
	ifc2.Interface = make([]*openconfig_interfaces.Interfaces_InterfaceKey,1);
	ifc2.Interface[0] = &ifk
	
	var dev openconfig.Device
	dev.Interfaces = &ifc2

	// msg := fmt.Sprintf("%s", dev)

	var opt protojson.MarshalOptions
	opt.EmitUnpopulated = true
	opt.EmitDefaultValues = true
	jsonDump, err := opt.Marshal(&dev)

	err = os.WriteFile("config.empty.json", jsonDump, 0644)
	if err != nil {
		msg := fmt.Sprintf("parseProto failed to write to file")
		log.Error(msg)
	}

	jsonDump, err = os.ReadFile("config.json")
	if err != nil {
		msg := fmt.Sprintf("parseProto failed to read file")
		log.Error(msg)
	}

	err = protojson.Unmarshal(jsonDump, &dev)
	if err != nil {
		msg := fmt.Sprintf("parseProto: error in Unmarshaling: %v", err)
		log.Error(msg)
	}

	// log.Info(msg)
}

func myCallback(config ygot.ValidatedGoStruct) error {
	log.Info("myCallback")

	jsonType := "Internal"
	_ = jsonType

	// jsonTree, err := ygot.ConstructInternalJSON(config)
	// jsonTree, err := ygot.ConstructIETFJSON(config, &ygot.RFC7951JSONConfig{AppendModuleName: false})
	// if err != nil {
	// 	msg := fmt.Sprintf("myCallback: error in constructing IETF JSON tree from config struct: %v", err)
	// 	log.Error(msg)
	// 	return status.Error(codes.Internal, msg)
	// }

	// jsonDump, err := json.Marshal(jsonTree)
	// if err != nil {
	// 	msg := fmt.Sprintf("myCallback: error in marshaling %s JSON tree to bytes: %v", jsonType, err)
	// 	log.Error(msg)
	// 	return status.Error(codes.Internal, msg)
	// }

	jsonDump, err := ygot.EmitJSON(config, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: false,
		 },
	})
	if err != nil {
		msg := fmt.Sprintf("myCallback: error in EmitJSON: %v", err)
		log.Error(msg)
		return status.Error(codes.Internal, msg)
	}	

	_ = jsonDump
	file, err := os.OpenFile("config.json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	file.WriteString(jsonDump)
	// err = os.WriteString("config.json", jsonDump, 0644)
	if err != nil {
		msg := fmt.Sprintf("myCallback failed to write to file")
		log.Error(msg)
		return status.Error(codes.Internal, msg)
	}

	parseProto()
	errors.New("an error msg")
	return nil
}

func newServer(model *gnmi.Model, config []byte) (*server, error) {
	s, err := gnmi.NewServer(model, config, myCallback)
	if err != nil {
		fmt.Print("dupa1\n")
		return nil, err
	}
	return &server{Server: s}, nil
}

// Get overrides the Get func of gnmi.Target to provide user auth.
func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Get request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	log.Infof("allowed a Get request: %v", msg)
	return s.Server.Get(ctx, req)
}

// Set overrides the Set func of gnmi.Target to provide user auth.
func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	log.Infof("allowed a Set request: %v", msg)
	return s.Server.Set(ctx, req)
}

// Set overrides the Subscribe func of gnmi.Target to provide user auth.
func (s *server) Subscribe(stream pb.GNMI_SubscribeServer) error {
	msg, ok := credentials.AuthorizeUser(stream.Context())
	if !ok {
		log.Infof("denied a Subscribe request: %v", msg)
		return status.Error(codes.PermissionDenied, msg)
	}
	log.Infof("allowed a Subscribe request: %v", msg)
	return s.Server.Subscribe(stream)
}

func main() {
	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*gostruct.Device)(nil)),
		gostruct.SchemaTree["Device"],
		gostruct.Unmarshal,
		gostruct.ΛEnum)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Supported models:\n")
		for _, m := range model.SupportedModels() {
			fmt.Fprintf(os.Stderr, "  %s\n", m)
		}
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Set("logtostderr", "true")
	flag.Parse()

	opts := credentials.ServerCredentials()
	g := grpc.NewServer(opts...)

	var configData []byte
	if *configFile != "" {
		var err error
		configData, err = ioutil.ReadFile(*configFile)
		if err != nil {
			log.Exitf("error in reading config file: %v", err)
		}
	}
	s, err := newServer(model, configData)
	if err != nil {
		log.Exitf("error in creating gnmi target: %v", err)
	}
	pb.RegisterGNMIServer(g, s)
	reflection.Register(g)

	log.Infof("starting to listen on %s", *bindAddr)
	listen, err := net.Listen("tcp", *bindAddr)
	if err != nil {
		log.Exitf("failed to listen: %v", err)
	}

	log.Info("starting to serve")
	if err := g.Serve(listen); err != nil {
		log.Exitf("failed to serve: %v", err)
	}
}
