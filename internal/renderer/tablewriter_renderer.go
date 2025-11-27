package renderer

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type TableWriterRenderer struct{}

func (t *TableWriterRenderer) RenderTable(file *os.File, headers []string, rows [][]string) (int, error) {
	table := tablewriter.NewWriter(file)
	table.Options(tablewriter.WithLineCounter())
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Row.Formatting = tw.CellFormatting{AutoWrap: tw.WrapTruncate, Alignment: tw.AlignLeft}
		// cfg.Row.Formatting.AutoFormat = tw.Success
		cfg.Row.ColMaxWidths = tw.CellWidth{Global: 50}
		cfg.Behavior.TrimSpace = tw.Fail
	})
	table.Header(headers)

	for _, row := range rows {
		table.Append(row)
	}
	table.Lines()

	err := table.Render()

	return table.Lines(), err
}
