package cli

import (
	"cliTodoist/internal/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/renderer"
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
	Input input.Input
}

func (cli *CLI) Run() error {
	db, err := storage.NewDB(&renderer.TableWriterRenderer{})
	if err != nil {
		return err
	}

	menu := NewMenu(colors.Bold +
		colors.Yellow +
		"Please select type of operation you want to do with Todoist?\n" +
		colors.Reset + util.NavigationPrompt +
		colors.Reset)

	deleteMenu := NewMenu(colors.Bold +
		colors.Yellow +
		"Select type of interaction\n" +
		colors.Reset + util.NavigationPrompt +
		colors.Reset)
	deleteMenu.AddItem("Normal", StringToID("Delete task Normal"), nil)
	deleteMenu.AddItem("Interactive", StringToID("Delete task Interactive"), nil)

	updateMenu := NewMenu(colors.Bold +
		colors.Yellow + "What do you want to update?\n" +
		colors.Reset + util.NavigationPrompt +
		colors.Reset)

	updateMenuHeader := NewMenu(colors.Bold +
		colors.Yellow +
		"Select type of interaction\n" +
		colors.Reset + util.NavigationPrompt +
		colors.Reset)
	updateMenuHeader.AddItem("Normal", StringToID("Update task Normal"), nil)
	updateMenuHeader.AddItem("Interactive", StringToID("Update task Interactive"), nil)

	updateMenu.AddItem("Header", StringToID("Update header"), updateMenuHeader)
	updateMenu.AddItem("Status", StringToID("Update status"), nil)

	menu.AddItem("Add task", StringToID("Add task"), nil)
	menu.AddItem("Update task", StringToID("Update task"), updateMenu)
	menu.AddItem("Delete task", StringToID("Delete task"), deleteMenu)
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

func ShowMenu(m *Menu, db *storage.DB, input input.Input) error {

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
			itemID = ""
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
			if !response {
				break
			}
		}
	case StringToID("Update task Normal"):
		for {
			response, err := db.UpdateItem(input)
			if err != nil {
				return err
			}

			if !response {
				break
			}
		}
	case StringToID("Update task Interactive"):
		for {
			response, err := db.UpdateItemInteractive(input)
			if err != nil {
				return err
			}

			if !response {
				break
			}
		}
	case StringToID("Update status"):
		for {
			response, err := db.UpdateStatus(input)
			if err != nil {
				return err
			}
			if !response {
				break
			}
		}
	case StringToID("Delete task Normal"):
		for {
			response, err := db.DeleteItem(input)
			if err != nil {
				return err
			}
			if !response {
				break
			}
		}
	case StringToID("Delete task Interactive"):
		for {
			response, err := db.DeleteItemInteractive(input)
			if err != nil {
				return err
			}
			if !response {
				break
			}
		}
	case StringToID("List all task"):
		db.ShowAllItems(input)
	case "exit":
		return errors.New("Exit application")
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
