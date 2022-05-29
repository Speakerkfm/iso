package tablewriter

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

type Table interface {
	SetHeader(keys []string)
	SetRowLine(line bool)
	Append(row []string)
	Render()
}

func NewWriter(writer io.Writer) Table {
	return tablewriter.NewWriter(writer)
}
