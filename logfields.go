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

		// Split the generated code into sections grouped by the outermost-level message being processed
		p.P()

		proto3 := gogoproto.IsProto3(file.FileDescriptorProto)
		for _, field := range msg.GetField() {
			p.generateFieldExtractor(msg, field, proto3)
		}
		p.generateLogsExtractor(msg, proto3)
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
	return lowerCamel(p.GetOneOfFieldName(msg, field)) + `Fields`
}

func (p *plugin) GetFieldMethod(msg *generator.Descriptor, field *descriptor.FieldDescriptorProto) string {
	return lowerCamel(p.GetOneOfFieldName(msg, field)) + `LogFields`
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

func (p *plugin) generateFieldExtractor(msg *generator.Descriptor, field *descriptor.FieldDescriptorProto, proto3 bool) {
	typeName := generator.CamelCaseSlice(msg.TypeName())
	funcName := p.GetFieldMethod(msg, field)
	fieldName := p.GetFieldName(msg, field)
	if field.IsMessage() && field.OneofIndex != nil {
		// Messages in a oneof are never repeated, but we need to do a typecast.
		p.P(`func (this *`, typeName, `) `, funcName, `() map[string][]string {`)
		p.In()
		p.P(`if f, ok := this.`, fieldName, `.(*`, p.OneOfTypeName(msg, field), `); ok {`)
		p.In()
		p.P(`return f.`, p.GetOneOfFieldName(msg, field), `.LogFields()`)
		p.Out()
		p.P(`}`)
		p.P(`return map[string][]string{}`)
		p.Out()
		p.P(`}`)
	} else if field.IsMessage() && field.IsRepeated() {
		// For repeated message fields, we need to gather the maps from each item
		p.P(`func (this *`, typeName, `) `, funcName, `() map[string][]string {`)
		p.In()
		p.P(`fields := map[string][]string{}`)
		p.P(`for _, msg := range this.`, fieldName, ` {`)
		p.In()
		p.P(`for k, v := range msg.LogFields() {`)
		p.In()
		p.P(`fields[k] = append(fields[k], v...)`)
		p.Out()
		p.P(`}`)
		p.Out()
		p.P(`}`)
		p.P(`return fields`)
		p.Out()
		p.P(`}`)
	} else if !field.IsMessage() && field.IsRepeated() && hasLogField(field) {
		// For repeated primitive fields, we need to format each item and put them in a slice
		p.P(`func (this *`, typeName, `) `, funcName, `() []string {`)
		p.In()
		p.P(`var vals []string`)
		p.P(`for _, val := range this.`, fieldName, ` {`)
		p.In()
		p.P(`vals = append(vals, `, p.getFmtExpr(`val`, field), `)`)
		p.Out()
		p.P(`}`)
		p.P(`return vals`)
		p.Out()
		p.P(`}`)
	}
}

func (p *plugin) generateFieldsLiteralReturn(msg *generator.Descriptor, proto3 bool) {
	fieldExpr := map[string]string{}
	// For output determinism, we track the ordering based on first occurrence of a log field name.
	var nameOrder []string
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

		if _, ok := fieldExpr[logField.Name]; !ok {
			nameOrder = append(nameOrder, logField.Name)
		}

		if field.IsRepeated() {
			expr := `this.` + lowerCamel(p.GetFieldName(msg, field)) + `LogFields()`
			if fieldExpr[logField.Name] == "" {
				fieldExpr[logField.Name] = expr
			} else {
				fieldExpr[logField.Name] = fmt.Sprintf(`append(%v, %v...)`, fieldExpr[logField.Name], expr)
			}
		} else {
			expr := p.getFieldFmtExpr(msg, field, proto3)
			if fieldExpr[logField.Name] == "" {
				fieldExpr[logField.Name] = fmt.Sprintf(`[]string{%v}`, expr)
			} else {
				fieldExpr[logField.Name] = fmt.Sprintf(`append(%v, %v)`, fieldExpr[logField.Name], expr)
			}
		}
	}

	p.P(`return map[string][]string{`)
	p.In()
	for _, name := range nameOrder {
		p.P(strconv.Quote(name), `: `, fieldExpr[name], `,`)
	}
	p.Out()
	p.P(`}`)
}

func (p *plugin) generateLogsExtractor(msg *generator.Descriptor, proto3 bool) {
	p.P(`func (this *`, generator.CamelCaseSlice(msg.TypeName()), `) LogFields() map[string][]string {`)
	p.In()

	var hasLoggedOneOf bool
	for _, field := range msg.GetField() {
		if field.OneofIndex == nil {
			continue
		}
		if field.IsMessage() || hasLogField(field) {
			hasLoggedOneOf = true
			break
		}
	}

	if !hasLoggedOneOf {
		var needsBody bool
		for _, field := range msg.GetField() {
			if field.IsMessage() || hasLogField(field) {
				needsBody = true
				break
			}
		}
		if !needsBody {
			// If the message has nothing that might generate log fields, we can immediately return an empty map and skip everything else
			p.P(`return map[string][]string{}`)
			p.Out()
			p.P(`}`)
			return
		}
	}

	p.P(`// Handle being called on nil message.`)
	p.P(`if this == nil {`)
	p.In()
	p.P(`return map[string][]string{}`)
	p.Out()
	p.P(`}`)

	if !hasLoggedOneOf {
		var hasChildren bool
		for _, field := range msg.GetField() {
			if field.IsMessage() {
				hasChildren = true
				break
			}
		}
		if !hasChildren {
			// If there were no message fields, return a fields literal directly
			p.P(`// Generate fields for this message.`)
			p.generateFieldsLiteralReturn(msg, proto3)
			p.Out()
			p.P(`}`)
			return
		}
	}

	p.P(`// Gather fields from child messages.`)
	p.P(`var hasInner bool`)
	for _, field := range msg.GetField() {
		if !field.IsMessage() && field.OneofIndex != nil && hasLogField(field) {
			p.P(`hasInner = hasInner || this.`, p.GetFieldName(msg, field), ` != nil`)
			continue
		} else if !field.IsMessage() {
			continue
		} else if field.IsRepeated() {
			p.P(p.GetFieldVar(msg, field), ` := this.`, p.GetFieldMethod(msg, field), `()`)
		} else if field.OneofIndex != nil {
			p.P(p.GetFieldVar(msg, field), ` := this.`, p.GetFieldMethod(msg, field), `()`)
		} else {
			p.P(p.GetFieldVar(msg, field) + ` := this.` + p.GetFieldName(msg, field) + `.LogFields()`)
		}
		p.P(`hasInner = hasInner || len(` + p.GetFieldVar(msg, field) + `) > 0`)
	}
	p.P(`if !hasInner {`)
	p.In()
	p.P(`// If no inner messages added any fields, the fields map is complete.`)
	p.generateFieldsLiteralReturn(msg, proto3)
	p.Out()
	p.P(`}`)

	p.P(`// Merge all the field maps.`)
	p.P(`res := map[string][]string{}`)
	for _, field := range msg.GetField() {
		if field.IsMessage() {
			p.P(`for k, v := range ` + p.GetFieldVar(msg, field) + ` {`)
			p.In()
			p.P(`res[k] = append(res[k], v...)`)
			p.Out()
			p.P(`}`)
			continue
		}
		logField := getLogFieldIfAny(field)
		if logField == nil {
			continue
		}
		quoted := strconv.Quote(logField.Name)
		if field.IsRepeated() {
			p.P(`res[`, quoted, `] = append(res[`, quoted, `], this.`, p.GetFieldMethod(msg, field), `()...)`)
		} else if field.OneofIndex != nil {
			p.P(`if f, ok := this.`, p.GetFieldName(msg, field), `.(*`, p.OneOfTypeName(msg, field), `); ok {`)
			p.In()
			p.P(`res[`, quoted, ` ] = append(res[`, quoted, `], `, p.getFmtExpr(`f.`+p.GetOneOfFieldName(msg, field), field), `)`)
			p.Out()
			p.P(`}`)
		} else {
			p.P(`res[`, quoted, `] = append(res[`, quoted, `], `, p.getFieldFmtExpr(msg, field, proto3), `)`)
		}
	}

	p.P(`return res`)
	p.Out()
	p.P(`}`)
}
