package dns

import (
	"fmt"
	"io"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/olekukonko/ts"
)

// PrettyPrint tries to print a struct to io.writer as a table
// with optional list of tags being printed too
// note that this function should have been moved to `cmd` but
// since it's being used as a Marshaller, it's been moved to here
// to facilitate importing.
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
			// TODO: it could be possible to show lists better
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
