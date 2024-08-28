package dns

import (
	"fmt"
	"io"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/olekukonko/ts"
)

func StructFormatter(myStruct any, tags ...string) (rows []table.Row) {
	v := reflect.ValueOf(myStruct)
	tv := reflect.TypeOf(myStruct)
	for i := 0; i < tv.NumField(); i++ {
		// see if the field is exportable
		if v.Field(i).CanInterface() {
			field := tv.Field(i).Name
			row := table.Row{field}
			for _, tag := range tags {
				row = append(row, GetTagValue(myStruct, field, tag))
			}

			switch tv.Field(i).Type.Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(v.Field(i).Interface())
				for j := 0; j < s.Len(); j++ {
					// slice of structs are special
					sv := reflect.TypeOf(s.Index(j).Interface())
					if sv.Kind() == reflect.Struct {
						childStruct := reflect.Indirect(s.Index(j)).Interface()
						for _, r := range StructFormatter(childStruct) {
							tmpRow := append(row, r...)
							rows = append(rows, tmpRow)
						}

					} else {
						tmpRow := append(row, fmt.Sprintf("%v", s.Index(j)))
						rows = append(rows, tmpRow)
					}
				}

			case reflect.Struct:
				childStruct := reflect.Indirect(reflect.ValueOf(myStruct)).
					FieldByName(tv.Field(i).Name).Interface()
				// Tags of child structs are not supported
				for _, r := range StructFormatter(childStruct) {
					tmpRow := append(row, r...)
					rows = append(rows, tmpRow)
				}
			default:
				tmpRow := append(row, fmt.Sprintf("%v", v.Field(i).Interface()))
				rows = append(rows, tmpRow)
			}
		}

	}
	return
}

// PrettyPrint tries to print a struct to io.writer as a table
// with optional list of tags being printed too
// note that this function should have been moved to `cmd` but
// since it's being used as a Marshaller, it's been moved to here
// to facilitate importing.
func PrettyPrint(myStruct any, w io.Writer, tags ...string) {
	t := table.NewWriter()
	headervals := []string{"param"}
	headervals = append(headervals, tags...)
	headers := table.Row{}
	for _, v := range headervals {
		headers = append(headers, v)
	}
	headers = append(headers, "value")
	t.AppendHeader(headers)

	size, _ := ts.GetSize()
	// wrap is a small function to trim outputs in tables
	wrap := text.Transformer(func(val interface{}) string {
		return text.WrapSoft(val.(string), size.Col()/(2+len(tags)))
	})
	t.SetAllowedRowLength(size.Col())

	rows := StructFormatter(myStruct, tags...)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendRows(rows, rowConfigAutoMerge)

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Transformer: wrap, AutoMerge: true},
		{Number: 2, Transformer: wrap, AutoMerge: true},
		{Number: 3, Transformer: wrap, AutoMerge: true},
	})
	// TODO: Style needs to change when the output is not a TTY
	// t.SetStyle(table.StyleColoredBright)
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	t.SetOutputMirror(w)
	t.Render()

}
