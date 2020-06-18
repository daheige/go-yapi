package plugin

import (
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func init() {
	generator.RegisterPlugin(new(yapi))
}

type yapi struct {
	gen *generator.Generator
}

func (y *yapi) Name() string {
	return "go-yapi"
}

func (y *yapi) Init(g *generator.Generator) {
	y.gen = g
}

func (y *yapi) GenerateImports(file *generator.FileDescriptor) {}

// generate yapi file
func (y *yapi) Generate(file *generator.FileDescriptor) {
	y.gen.P("// NOT EDIT")
	for _, service := range file.FileDescriptorProto.Service {
		y.generateService(service, file.MessageType)
	}
}

// yapi read rpc interface request and reply to generate yapi files
func (y *yapi) generateService(service *pb.ServiceDescriptorProto, messages []*pb.DescriptorProto) {
	y.gen.P("// generate service ", service.Name)
	for _, method := range service.Method {
		y.gen.P("// generate method ", method.Name)
		y.generateMethod(method, messages)
	}
}

func (y *yapi) generateMethod(method *pb.MethodDescriptorProto, messages []*pb.DescriptorProto) {
	for _, msg := range messages {
		if msg.Name == method.InputType {
			y.gen.P("generate input request")
			y.generateMessage(msg)
		} else if msg.Name == method.OutputType {
			y.gen.P("generate output request")
			y.generateMessage(msg)
		}
	}
}

// object type description properties (fields) required
func (y *yapi) generateMessage(msg *pb.DescriptorProto) {
	y.gen.P(`{"type":"object","properties":{`)
	for _, filed := range msg.Field {
		y.generateFiled(filed)
	}
	y.gen.P("}")
}

//
func (y *yapi) generateFiled(filed *pb.FieldDescriptorProto) {
	y.gen.P(`"`, filed.Name, `":{`)

	y.gen.P("}")
}
