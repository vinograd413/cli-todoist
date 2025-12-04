package table

import (
	"cliTodoist/internal/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/renderer"
	"cliTodoist/internal/tasks"
	"cliTodoist/internal/util"
	"slices"

	"fmt"
	"os"
	"strconv"
	"time"
)

type Table struct {
	Renderer       renderer.TableRenderer
	CursorPosition int
	TaskList       []*tasks.Task
	SelectedTasks  []*tasks.Task
}

func (t *Table) DisplayTableMultiSelection(i input.Input) error {
	var err error
	var tableLines int

	tableLines, err = t.PrintTasksAsTableWithSelection(i.File())
	if err != nil {
		return err
	}

	oldState, err := i.SetRawMode()
	if err != nil {
		return err
	}

	for {
		keyInput, err := i.ReadKey()
		if err != nil {
			return err
		}
		switch keyInput {
		case util.MoveUp:
			if t.CursorPosition > 0 {
				t.CursorPosition--
			}
		case util.MoveDown:
			if t.CursorPosition < len(t.TaskList)-1 {
				t.CursorPosition++
			}
		case util.Enter:
			i.RestoreMode(oldState)
			return nil
		case util.Space:
			if !slices.Contains(t.SelectedTasks, t.TaskList[t.CursorPosition]) {
				t.SelectedTasks = append(t.SelectedTasks, t.TaskList[t.CursorPosition])
			} else {
				t.SelectedTasks = util.RemoveFirst(t.SelectedTasks, t.TaskList[t.CursorPosition])
			}
		case util.Escape:
			i.RestoreMode(oldState)
			t.SelectedTasks = nil
			return nil
		}
		i.RestoreMode(oldState)
		fmt.Printf(util.CursorUpFormat, tableLines)
		fmt.Print(util.ClearDown)
		tableLines, err = t.PrintTasksAsTableWithSelection(i.File())
		oldState, err = i.SetRawMode()
		if err != nil {
			return err
		}
	}
}

func (t *Table) DisplayTableSingleSelection(i input.Input) error {
	var err error
	var tableLines int

	tableLines, err = t.PrintTasksAsTableWithSelection(i.File())
	if err != nil {
		return err
	}

	oldState, err := i.SetRawMode()
	if err != nil {
		return err
	}

	for {
		keyInput, err := i.ReadKey()
		if err != nil {
			return err
		}
		switch keyInput {
		case util.MoveUp:
			if t.CursorPosition > 0 {
				t.CursorPosition--
			}
		case util.MoveDown:
			if t.CursorPosition < len(t.TaskList)-1 {
				t.CursorPosition++
			}
		case util.Enter:
			t.SelectedTasks = append(t.SelectedTasks, t.TaskList[t.CursorPosition])
			i.RestoreMode(oldState)
			return nil
		case util.Escape:
			i.RestoreMode(oldState)
			t.SelectedTasks = nil
			return nil
		}
		i.RestoreMode(oldState)
		fmt.Printf(util.CursorUpFormat, tableLines)
		fmt.Print(util.ClearDown)
		tableLines, err = t.PrintTasksAsTableWithSelection(i.File())
		oldState, err = i.SetRawMode()
		if err != nil {
			return err
		}
	}
}

func (t *Table) PrintTasksAsTableWithSelection(file *os.File) (int, error) {
	header := []string{"#", "Text", "Created At", "Completed", "Completed At"}
	var rows [][]string
	var maxRow int

	pageNum := (t.CursorPosition / util.MaxTableHeight)
	pageCursorPosition := t.CursorPosition % util.MaxTableHeight

	if len(t.TaskList)/util.MaxTableHeight > pageNum {
		maxRow = pageNum*util.MaxTableHeight + util.MaxTableHeight
	} else {
		maxRow = len(t.TaskList)
	}

	pageTasks := t.TaskList[pageNum*util.MaxTableHeight : maxRow]

	for i, task := range pageTasks {
		var bgColor string
		var cursor string = "   "

		if slices.Contains(t.SelectedTasks, task) && i != pageCursorPosition {
			cursor = "  " + util.SymbRedCross
		} else if slices.Contains(t.SelectedTasks, task) && i == pageCursorPosition {
			bgColor = colors.SetBackgroundColor(240)
			cursor = util.CursorSelectionDeletion
		} else if !slices.Contains(t.SelectedTasks, task) && i == pageCursorPosition {
			bgColor = colors.SetBackgroundColor(240)
			cursor = util.CursorSelection
		}

		var row []string
		taskNumber := i + 1 + (util.MaxTableHeight * pageNum)
		row = append(row, cursor+bgColor+strconv.Itoa(taskNumber)+colors.Reset)
		row = append(row, bgColor+task.Text+colors.Reset)
		row = append(row, bgColor+time.Unix(task.CreatedAt, 0).Format(time.DateTime)+colors.Reset)
		if !task.IsCompleted {
			row = append(row, util.SymbRedCross)
			row = append(row, "")
		} else {
			row = append(row, util.SymbGreenSimleCheck)
			row = append(row, bgColor+time.Unix(task.CompletedAt, 0).Format(time.DateTime)+colors.Reset)
		}
		rows = append(rows, row)
	}

	if len(t.TaskList)/util.MaxTableHeight > pageNum {
		lastRow := []string{"", "...", "", "", ""}
		rows = append(rows, lastRow)
	}

	return t.Renderer.RenderTable(file, header, rows)
}

func (t *Table) PrintTasksAsTable(file *os.File) (int, error) {
	header := []string{"#", "Text", "Created At", "Completed", "Completed At"}
	var rows [][]string
	for i, task := range t.TaskList {
		var row []string
		row = append(row, strconv.Itoa(i+1))
		row = append(row, task.Text)
		row = append(row, time.Unix(task.CreatedAt, 0).Format(time.DateTime))
		if !task.IsCompleted {
			row = append(row, util.SymbRedCross)
			row = append(row, "")
		} else {
			row = append(row, util.SymbGreenSimleCheck)
			row = append(row, time.Unix(task.CompletedAt, 0).Format(time.DateTime))
		}
		rows = append(rows, row)
	}

	return t.Renderer.RenderTable(file, header, rows)
}
