package renderer

import (
	"os"

	"github.com/aquasecurity/table"
)

type AquaSecTableRenderer struct{}

func (t *AquaSecTableRenderer) RenderTable(file *os.File, headers []string, rows [][]string) (int, error) {
	table := table.New(file)
	table.SetHeaders(headers...)
	table.AddRows(rows...)

	table.Render()

	return table.RowCount(), nil
}
