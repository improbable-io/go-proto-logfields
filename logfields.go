// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

// We generate a LogFields() method which returns a map. This method does the
// following for each field tagged with the logField annotation:
// 
// * Primitive fields - adds the name and value of the filed to the map
// * Messages - Calls LogFields on the message and merges the returned map into the current map 
// 
// Oneofs containing logfields and embedded messages are supported.
// Duplicate names in the same message are reported as an error by the generator.
// Repeated fields, and therefore maps, are not supported, and are ignored by the generator.

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
		logNames := map[string]struct{}{}
		for _, field := range msg.GetField() {
			logField := getLogFieldIfAny(field)
			if logField == nil {
				continue
			}
			if field.IsRepeated() {
				p.Fail(fmt.Sprintf("Cannot log repeated field %v.%v", msg.GetName(), field.GetName()))
			}
			if _, duplicate := logNames[logField.Name]; duplicate {
				p.Fail(fmt.Sprintf("Duplicate log field name %v in %v", logField.Name, msg.GetName()))
			}
			logNames[logField.Name] = struct{}{}
		}
	}

	for _, msg := range file.Messages() {
		if msg.GetOptions().GetMapEntry() {
			continue
		}

		// Split the generated code into sections grouped by the outermost-level message being processed
		p.P()
		p.generateLogsExtractor(msg, gogoproto.IsProto3(file.FileDescriptorProto))
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
	logField := e.(*LogField)
	if logField != nil && logField.Name == "" {
		logField = nil
	}
	return logField
}

func hasLogField(field *descriptor.FieldDescriptorProto) bool {
	return getLogFieldIfAny(field) != nil
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

func (p *plugin) GetFieldMethod(msg *generator.Descriptor, field *descriptor.FieldDescriptorProto) string {
	return lowerCamel(p.GetFieldName(msg, field)) + `LogFields`
}

func (p *plugin) getFieldFmtExpr(msg *generator.Descriptor, field *descriptor.FieldDescriptorProto, proto3 bool) string {
	expr := p.GetFieldName(msg, field)
	if proto3 {
		expr = `this.` + expr
	} else {
		expr = `this.Get` + expr + `()`
	}
	return p.getFmtExpr(expr, field)
}

func (p *plugin) getFmtExpr(fieldName string, field *descriptor.FieldDescriptorProto) string {
	fmtExpr := fieldName
	if field.IsString() {
		// no need to convert strings
	} else if field.IsBytes() {
		// pass through to string, let log handlers deal with non-printable bytes
		fmtExpr = `string(` + fmtExpr + `)`
	} else {
		fmtExpr = `fmt.Sprintf("%v", ` + fmtExpr + `)`
	}
	return fmtExpr
}

func (p *plugin) generateFieldsLiteralReturn(msg *generator.Descriptor, proto3 bool) {
	p.P(`return map[string]string{`)
	p.In()
	for _, field := range msg.GetField() {
		if field.IsMessage() {
			continue
		}
		if field.OneofIndex != nil {
			continue
		}
		logField := getLogFieldIfAny(field)
		if logField == nil {
			continue
		}

		expr := p.getFieldFmtExpr(msg, field, proto3)
		p.P(strconv.Quote(logField.Name), `: `, expr, `,`)
	}
	p.Out()
	p.P(`}`)
}

func findLoggedOneofField(msg *generator.Descriptor, oneofIndex int) *descriptor.FieldDescriptorProto {
	var loggedOneOfField *descriptor.FieldDescriptorProto
	for _, field := range msg.GetField() {
		if field.OneofIndex == nil || int(field.GetOneofIndex()) != oneofIndex {
			continue
		}
		if field.IsMessage() || hasLogField(field) {
			loggedOneOfField = field
			break
		}
	}
	return loggedOneOfField
}

func (p *plugin) generateOneOfSwitch(msg *generator.Descriptor, oneofProxy *descriptor.FieldDescriptorProto) string {
	oneOfVar := p.GetFieldVar(msg, oneofProxy)
	p.P(`var `, oneOfVar, ` map[string]string`)
	p.P(`switch f := this.`, p.GetFieldName(msg, oneofProxy), `.(type) {`)
	for _, field := range msg.GetField() {
		if field.OneofIndex == nil || field.GetOneofIndex() != oneofProxy.GetOneofIndex() {
			continue
		}
		// Oneof fields that can't generate log fields will use the default clause.
		if !field.IsMessage() && !hasLogField(field) {
			continue
		}
		p.P(`case *`, p.OneOfTypeName(msg, field), `:`)
		p.In()
		if field.IsMessage() {
			p.P(oneOfVar, ` = f.`, p.GetOneOfFieldName(msg, field), `.LogFields()`)
		} else {
			logName := getLogFieldIfAny(field).Name
			p.P(oneOfVar, ` = map[string]string{`, strconv.Quote(logName), `: `, p.getFmtExpr(`f.`+p.GetOneOfFieldName(msg, field), field), `}`)
		}
		p.Out()
	}
	p.P(`default:`)
	p.In()
	p.P(oneOfVar, ` = map[string]string{}`)
	p.Out()
	p.P(`}`)
	return oneOfVar
}

func (p *plugin) generateLogsExtractor(msg *generator.Descriptor, proto3 bool) {
	p.P(`func (this *`, generator.CamelCaseSlice(msg.TypeName()), `) LogFields() map[string]string {`)
	p.In()

	var needsBody bool
	for _, field := range msg.GetField() {
		if field.IsRepeated() {
			continue
		}
		if field.IsMessage() || hasLogField(field) {
			needsBody = true
			break
		}
	}
	if !needsBody {
		// If the message has nothing that might generate log fields, we can immediately return an empty map and skip everything else
		p.P(`return map[string]string{}`)
		p.Out()
		p.P(`}`)
		return
	}

	p.P(`// Handle being called on nil message.`)
	p.P(`if this == nil {`)
	p.In()
	p.P(`return map[string]string{}`)
	p.Out()
	p.P(`}`)

	canUseLiteral := true
	for _, field := range msg.GetField() {
		if field.IsMessage() || field.OneofIndex != nil {
			canUseLiteral = false
			break
		}
	}
	if canUseLiteral {
		// If there were no message fields, return a fields literal directly
		p.P(`// Generate fields for this message.`)
		p.generateFieldsLiteralReturn(msg, proto3)
		p.Out()
		p.P(`}`)
		return
	}

	p.P(`// Gather fields from oneofs and child messages.`)
	p.P(`var hasInner bool`)
	loggedOneOfs := map[int]string{}
	// Generate code to build a log field map for each oenof.
	for oneOfIndex, _ := range msg.GetOneofDecl() {
		// loggedOneOfField is used later as a proxy for the oneof
		loggedOneOfField := findLoggedOneofField(msg, oneOfIndex)
		if loggedOneOfField == nil {
			// We can skip generating code for the oneof.
			continue
		}

		// Generate a type-switch.
		oneOfVar := p.generateOneOfSwitch(msg, loggedOneOfField)
		loggedOneOfs[oneOfIndex] = oneOfVar
		// Keep track of whether any log fields were generated at runtime.
		p.P(`hasInner = hasInner || len(`, oneOfVar, `) > 0`)
	}

	// Generate code for embedded messages
	for _, field := range msg.GetField() {
		if field.OneofIndex != nil {
			continue
		} else if !field.IsMessage() {
			continue
		} else if field.IsRepeated() {
			continue
		}
		p.P(p.GetFieldVar(msg, field) + ` := this.` + p.GetFieldName(msg, field) + `.LogFields()`)
		// Keep track of whether any log fields were generated at runtime.
		p.P(`hasInner = hasInner || len(` + p.GetFieldVar(msg, field) + `) > 0`)
	}
	p.P(`if !hasInner {`)
	p.In()
	p.P(`// If no inner messages added any fields, avoid merging maps.`)
	p.generateFieldsLiteralReturn(msg, proto3)
	p.Out()
	p.P(`}`)

	p.P(`// Merge all the field maps.`)
	p.P(`res := map[string]string{}`)
	// Generally try and keep the order of fields intact.
	// Merge each oneof on encountering the first field that belongs to it.
	visitedOneOfs := map[int]struct{}{}
	for _, field := range msg.GetField() {
		if field.IsRepeated() {
			continue
		}
		if field.OneofIndex != nil {
			oneOfVar, logged := loggedOneOfs[int(field.GetOneofIndex())]
			if !logged {
				continue
			}
			if _, visited := visitedOneOfs[int(field.GetOneofIndex())]; visited {
				continue
			}
			visitedOneOfs[int(field.GetOneofIndex())] = struct{}{}
			p.P(`for k, v := range ` + oneOfVar + ` {`)
			p.In()
			p.P(`res[k] = v`)
			p.Out()
			p.P(`}`)
		} else if field.IsMessage() {
			p.P(`for k, v := range ` + p.GetFieldVar(msg, field) + ` {`)
			p.In()
			p.P(`res[k] = v`)
			p.Out()
			p.P(`}`)
		} else {
			logField := getLogFieldIfAny(field)
			if logField == nil {
				continue
			}
			quoted := strconv.Quote(logField.Name)
			p.P(`res[`, quoted, `] = `, p.getFieldFmtExpr(msg, field, proto3))
		}
	}

	p.P(`return res`)
	p.Out()
	p.P(`}`)
}
