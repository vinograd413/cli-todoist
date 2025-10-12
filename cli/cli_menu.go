package cli

import (
	"cliTodoist/colors"
	"cliTodoist/internal/util"
	"errors"
	"fmt"
)

var (
	ErrNoMenuItems = errors.New("Menu has no items to display")
)

type Menu struct {
	Prompt       string
	CursorPos    int
	ScrollOffset int
	MenuItems    []*MenuItem
}

type MenuItem struct {
	Text    string
	ID      string
	SubMenu *Menu
}

func NewMenu(prompt string) *Menu {
	return &Menu{Prompt: prompt, MenuItems: make([]*MenuItem, 0)}
}

func (m *Menu) AddItem(text string, id string, subMenu *Menu) *MenuItem {
	var menuItem *MenuItem
	if subMenu != nil {
		menuItem = &MenuItem{Text: text, ID: id, SubMenu: subMenu}
	} else {
		menuItem = &MenuItem{Text: text, ID: id}
	}

	m.MenuItems = append(m.MenuItems, menuItem)
	return menuItem
}

func (m *Menu) Display(input util.Input) (string, error) {
	if len(m.MenuItems) == 0 {
		return "", ErrNoMenuItems
	}

	defer func() {
		fmt.Print(util.ShowCursor)
	}()
	ClearScreen()
	fmt.Print(m.Prompt)
	m.renderMenu(false)

	for {
		keyInput, err := input.ReadKey()
		if err != nil {
			return "", err
		}
		switch keyInput {
		case util.MoveUp:
			if m.CursorPos > 0 {
				m.CursorPos--
			}
		case util.MoveDown:
			if m.CursorPos < len(m.MenuItems)-1 {
				m.CursorPos++
			}
		case util.Enter:
			selectedItem := m.MenuItems[m.CursorPos]
			return selectedItem.ID, nil
		case util.Escape:
			return "exit", nil
		}
		m.renderMenu(true)
	}
}

func (m *Menu) renderMenu(redraw bool) {
	menuSize := len(m.MenuItems)
	if redraw {
		fmt.Printf(util.CursorUpFormat, menuSize)
	}
	cursor := "  "
	select_cursor := colors.Yellow + "> " + colors.Reset

	for i, v := range m.MenuItems {
		if i == m.CursorPos {
			fmt.Printf("\r%s %s\n", select_cursor, colors.SetBackgroundColor(240)+v.Text+colors.Reset)
		} else {
			fmt.Printf("\r%s %s\n", cursor, v.Text)
		}
	}
}
