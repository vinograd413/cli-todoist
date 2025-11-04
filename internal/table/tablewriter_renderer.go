package table

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableWriterRenderer struct{}

func (t *TableWriterRenderer) RenderTable(file *os.File, headers []string, rows [][]string) error {
	table := tablewriter.NewWriter(file)
	table.Header(headers)

	for _, row := range rows {
		table.Append(row)
	}

	err := table.Render()
	return err
}
