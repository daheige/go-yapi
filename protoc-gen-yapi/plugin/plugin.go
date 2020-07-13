package plugin

import (
	pb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/lanceryou/go-yapi/protoc-gen-yapi/generator"
	"strings"
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
	//y.gen.P("// NOT EDIT")
	for _, service := range file.FileDescriptorProto.Service {
		y.generateService(service, file.MessageType)
	}
}

// yapi read rpc interface request and reply to generate yapi files
func (y *yapi) generateService(service *pb.ServiceDescriptorProto, messages []*pb.DescriptorProto) {
	//y.gen.P("// generate service ", service.Name)
	y.gen.P("{")
	for i, method := range service.Method {
		//y.gen.P("// generate method ", method.Name)
		y.gen.P(`"`, method.Name, `"`, ":{")
		// type name
		y.generateMethod(method, messages)
		y.gen.P("}")
		if len(service.Method) > 1 && i != len(service.Method)-1 {
			y.gen.P(",")
		}
	}
	y.gen.P("}")
}

func (y *yapi) generateMethod(method *pb.MethodDescriptorProto, messages []*pb.DescriptorProto) {
	msg := matchMessage(messages, *method.InputType)
	if msg != nil {
		y.gen.P(`"`, *msg.Name, `"`, ":")
		y.gen.P("{")
		y.generateMessage(msg, messages)
		y.gen.P("}")
		y.gen.P()
	}

	msg = matchMessage(messages, *method.OutputType)
	if msg != nil {
		//y.gen.P("// generate output request", *msg.Name)
		y.gen.P(`,"`, *msg.Name, `"`, ":")
		y.gen.P("{")
		y.generateMessage(msg, messages)
		y.gen.P("}")
		y.gen.P()
	}
}

// object type description properties (fields) required
func (y *yapi) generateMessage(msg *pb.DescriptorProto, messages []*pb.DescriptorProto) {
	if msg == nil {
		return
	}

	y.gen.P(`"type":"object","properties":{`)
	for i, filed := range msg.Field {
		y.generateFiled(filed, messages)
		if len(msg.Field) > 1 && i != len(msg.Field)-1 {
			y.gen.P(",")
		}
	}
	y.gen.P("}")
}

//
func (y *yapi) generateFiled(field *pb.FieldDescriptorProto, msgs []*pb.DescriptorProto) {
	y.gen.P(`"`, field.Name, `":{`)
	typ := jsonType(field)
	y.gen.P(`"type":"`, typ, `",`)
	switch typ {
	case "array":
		y.generateArray(field, msgs)
	case "object":
		y.generateMessage(matchMessage(msgs, *field.Name), msgs)
	case "string", "number", "boolean":
		y.gen.P(`"description": `, `"comments"`)
	}
	y.gen.P("}")
}

func (y *yapi) generateArray(field *pb.FieldDescriptorProto, msgs []*pb.DescriptorProto) {
	y.gen.P(`"items": {`)
	if isMessage(field) {
		y.generateMessage(matchMessage(msgs, *field.TypeName), msgs)
	} else {
		y.gen.P(`"type":"`, field.TypeName, `",`)
		y.gen.P(`"description": `, `"comments"`)
	}
	y.gen.P("}")
}

/*
"dropshipper_info": {
			"type": "object",
			"properties": {
				"option": {
					"type": "string",
					"description": "\"true\" \"false\""
				}
			},
			"required": ["option"]
		},

"offset": {
			"type": "number",
			"description": "查询"
		},
		"field_1": {
			"type": "array",
			"items": {
				"type": "string"
			}
		},
	"field_1":{
            "type":"array",
            "items":{
                "type":"object",
                "properties":{
                    "xx":{
                        "type":"string"
                    }
                },
                "required":[
                    "xx"
                ]
            }
        }
*/
func jsonType(field *pb.FieldDescriptorProto) string {
	if *field.Label == pb.FieldDescriptorProto_LABEL_REPEATED {
		return "array"
	}

	switch *field.Type {
	case pb.FieldDescriptorProto_TYPE_STRING, pb.FieldDescriptorProto_TYPE_BYTES:
		return "string"
	case pb.FieldDescriptorProto_TYPE_BOOL:
		return "boolean"
	case pb.FieldDescriptorProto_TYPE_MESSAGE:
		return "object"
	default:
		return "number"
	}
}

func matchMessage(msgs []*pb.DescriptorProto, name string) *pb.DescriptorProto {
	for _, msg := range msgs {
		if strings.Contains(name, *msg.Name) {
			return msg
		}
	}

	panic(name)

	return nil
}

func isRepeated(field *pb.FieldDescriptorProto) bool {
	return field.Label != nil && *field.Label == pb.FieldDescriptorProto_LABEL_REPEATED
}

func isMessage(field *pb.FieldDescriptorProto) bool {
	return *field.Type == pb.FieldDescriptorProto_TYPE_MESSAGE
}

func isString(field *pb.FieldDescriptorProto) bool {
	return *field.Type == pb.FieldDescriptorProto_TYPE_STRING
}
