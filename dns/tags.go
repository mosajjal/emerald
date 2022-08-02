package dns

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"io"
	"reflect"

	"github.com/olekukonko/ts"
)

func GetTagValue(myStruct interface{}, myField string, myTag string) string {
	t := reflect.TypeOf(myStruct)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == myField {
			tag := field.Tag.Get(myTag)
			return tag
		}
	}

	return ""
}

// PrettyPrint tries to print a struct to io.writer as a table
// with optional list of tags being printed too
func PrettyPrint(myStruct any, w io.Writer, tags ...interface{}) {
	t := table.NewWriter()
	headers := table.Row{"parameter", "value"}
	headers = append(headers, tags...)
	t.AppendHeader(headers)

	v := reflect.ValueOf(myStruct)
	typeOfS := v.Type()

	size, _ := ts.GetSize()
	// wrap is a small function to trim outputs in tables
	wrap := text.Transformer(func(val interface{}) string {
		return text.WrapSoft(val.(string), size.Col()/(2+len(tags)))
	})
	t.SetAllowedRowLength(size.Col())

	for i := 0; i < v.NumField(); i++ {
		// see if the field is exportable
		if v.Field(i).CanInterface() {
			field := typeOfS.Field(i).Name
			row := table.Row{field}
			row = append(row, fmt.Sprintf("%+v", v.Field(i).Interface()))
			for _, tag := range tags {
				row = append(row, GetTagValue(myStruct, field, tag.(string)))
			}
			t.AppendRow(row)
		}

	}
	t.SetColumnConfigs([]table.ColumnConfig{{Number: 1, Transformer: wrap}, {Number: 2, Transformer: wrap}, {Number: 3, Transformer: wrap}})
	// todo: style needs to change when the output is not a TTY
	t.SetStyle(table.StyleColoredBright)
	t.SetOutputMirror(w)
	t.Render()

}
