// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package logfields

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity"
)

func init() {
	generator.RegisterPlugin(NewPlugin(true))
}

type plugin struct {
	*generator.Generator
	generator.PluginImports
	useGogo bool
}

func NewPlugin(useGogo bool) generator.Plugin {
	return &plugin{useGogo: useGogo}
}

func (p *plugin) Name() string {
	return "logfields"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
	p.PluginImports = generator.NewPluginImports(p.Generator)
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	if !p.useGogo {
		vanity.TurnOffGogoImport(file.FileDescriptorProto)
	}
	for _, msg := range file.Messages() {
		if msg.GetOptions().GetMapEntry() {
			continue
		}
		p.generateFieldsExtractor(msg, gogoproto.IsProto3(file.FileDescriptorProto))
	}
}

func getLogFieldIfAny(field *descriptor.FieldDescriptorProto) *LogField {
	opts := field.GetOptions()
	if opts == nil {
		return nil
	}
	e, err := proto.GetExtension(opts, E_Logfield)
	if err != nil {
		return nil
	}
	return e.(*LogField)
}

// Convert a UpperCamelCase to lowerCamelCase.
// Handles initialisms as the first word: HTMLThing is converted to htmlThing.
func lowerCamel(varName string) string {
	firstNonUpper := 0
	for ; firstNonUpper < len(varName); firstNonUpper++ {
		if !('A' <= varName[firstNonUpper] && varName[firstNonUpper] <= 'Z') {
			break
		}
	}
	lastInitUpper := firstNonUpper - 1
	if lastInitUpper < 0 {
		return varName
	} else if lastInitUpper == 0 {
		return strings.ToLower(varName[:1]) + varName[1:]
	} else if lastInitUpper == len(varName)-1 {
		return strings.ToLower(varName)
	} else {
		return strings.ToLower(varName[:lastInitUpper]) + varName[lastInitUpper:]
	}
}

func (p *plugin) GetFieldVar(msg *generator.Descriptor, field *descriptor.FieldDescriptorProto) string {
	return lowerCamel(p.GetFieldName(msg, field)) + `Fields`
}

func (p *plugin) generateFieldsExtractor(msg *generator.Descriptor, proto3 bool) {
	p.P(`func (this *`, generator.CamelCaseSlice(msg.TypeName()), `) LogFields() map[string]string {`)
	defer p.P(`}`)
	p.In()
	defer p.Out()

	p.P(`if this == nil {`)
	p.In()
	p.P(`return map[string]string{}`)
	p.Out()
	p.P(`}`)

	var hasChildren bool
	for _, field := range msg.GetField() {
		if field.IsMessage() && !field.IsRepeated() {
			hasChildren = true
		}
	}
	if hasChildren {
		p.P(`subCount := 0`)
		for _, field := range msg.GetField() {
			if field.IsMessage() && !field.IsRepeated() {
				p.P(p.GetFieldVar(msg, field) + ` := this.` + p.GetFieldName(msg, field) + `.LogFields()`)
				p.P(`subCount += len(` + p.GetFieldVar(msg, field) + `)`)
			}
		}
	}

	p.P(`fields := map[string]string{`)
	p.In()
	for _, field := range msg.GetField() {
		logField := getLogFieldIfAny(field)
		if logField == nil {
			continue
		}
		if field.IsMessage() {
			p.Fail(fmt.Sprintf("LogField annotations cannot be applied to messages, %v.%v is a message", msg.GetName(), field.GetName()))
		}
		if field.IsRepeated() {
			p.Fail(fmt.Sprintf("LogField annotations can only be applied to singular fields, %v.%v is repeated", msg.GetName(), field.GetName()))
		}
		var fmtExpr string
		if proto3 {
			fmtExpr = `this.` + p.GetFieldName(msg, field)
		} else {
			fmtExpr = `this.Get` + p.GetFieldName(msg, field) + `()`
		}
		if field.IsString() {
			// no need to convert strings
		} else if field.IsBytes() {
			// pass through to string, let log handlers deal with non-printable bytes
			fmtExpr = `string(` + fmtExpr + `)`
		} else {
			fmtExpr = `fmt.Sprintf("%v", ` + fmtExpr + `)`
		}
		p.P(strconv.Quote(logField.Name) + `: ` + fmtExpr + `,`)
	}
	p.Out()
	p.P(`}`)

	if !hasChildren {
		p.P(`return fields`)
		return
	}

	p.P(`if subCount == 0 {`)
	p.In()
	p.P(`return fields`)
	p.Out()
	p.P(`}`)

	p.P(`res := make(map[string]string, subCount + len(fields))`)
	for _, field := range msg.GetField() {
		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE && field.GetLabel() != descriptor.FieldDescriptorProto_LABEL_REPEATED {
			p.P(`for k, v := range ` + p.GetFieldVar(msg, field) + ` {`)
			p.In()
			p.P(`res[k] = v`)
			p.Out()
			p.P(`}`)
		}
	}
	p.P(`for k, v := range fields {`)
	p.In()
	p.P(`res[k] = v`)
	p.Out()
	p.P(`}`)
	p.P(`return res`)
}
