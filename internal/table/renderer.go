package table

import "os"

type TableRenderer interface {
	// RenderTable receives headers and rows, and prints or returns the table string
	RenderTable(file *os.File, headers []string, rows [][]string) error
}
