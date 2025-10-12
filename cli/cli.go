package cli

import (
	"cliTodoist/colors"
	"cliTodoist/internal/util"
	"cliTodoist/storage"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// main cli structure that holds and run a programm
type CLI struct {
	Input util.Input
}

func (cli *CLI) validateArgs() {
	// fmt.Fprintln(os.Stdout, os.Args[0])
	bg := colors.SetBackgroundColor(100)
	fmt.Fprintf(os.Stdout, bg+"Value: %s \n", os.Args[0])
	// fmt.Printf("\033[%dA", 7)
}

func (cli *CLI) Run() error {
	// reader := bufio.NewReader(os.Stdin)
	// scanner := bufio.NewScanner(os.Stdin)

	// Set terminal to raw mode

	db, err := storage.NewDB()
	if err != nil {
		return err
	}

	menu := NewMenu(colors.Bold +
		colors.Yellow +
		"Please select type of operation you want to do with Todoist?\n" +
		"\rUse ↑/↓ to navigate, Enter to select, Esc to exit:\n" +
		colors.Reset)

	menu.AddItem("Add task", StringToID("Add task"), nil)
	menu.AddItem("Update task", StringToID("Update task"), nil)
	menu.AddItem("Delete task", StringToID("Delete task"), nil)
	menu.AddItem("List all task", StringToID("List all task"), nil)

	// Channel to signal interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		<-sigChan
		// Restore terminal state before exiting
		fmt.Print(util.ShowCursor)
		os.Exit(1)
	}()

	for {
		err := ShowMenu(menu, db, cli.Input)
		if err != nil {
			return err
		}
	}
}

func ShowMenu(m *Menu, db *storage.DB, input util.Input) error {

	oldState, err := input.SetRawMode()
	if err != nil {
		return err
	}

	itemID, err := m.Display(input)
	if err != nil {
		return err
	}

	for _, v := range m.MenuItems {
		if itemID == v.ID && v.SubMenu != nil {
			input.RestoreMode(oldState)
			ShowMenu(v.SubMenu, db, input)
		}
	}

	input.RestoreMode(oldState)

	switch itemID {
	case StringToID("Add task"):
		for {
			response, err := db.AddItem(input)
			if err != nil {
				return err
			}
			if response != "y" {
				break
			}
		}
	case StringToID("Update task"):
		for {
			response, err := db.UpdateItem(input)
			if err != nil {
				return err
			}
			if response != "y" {
				break
			}
		}
	case StringToID("Delete task"):
		for {
			response, err := db.DeleteItem(input)
			if err != nil {
				return err
			}
			if response != "y" {
				break
			}
		}
	case StringToID("List all task"):
		db.ListAllItems(input)
	case "exit":
		return errors.New("Exit error")
	}

	return nil
}

func StringToID(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func ClearScreenCmd() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ClearScreen() {
	fmt.Print(util.ClearAndHideCursor)
}
