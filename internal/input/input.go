package input

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Input interface {
	// ReadLine reads a line of text (for prompts, etc.)
	ReadLine(prompt string) (string, error)
	// ReadKey reads a single key (for menu navigation)
	ReadKey() (byte, error)
	// SetRawMode set terminal into raw mode
	SetRawMode() (*term.State, error)
	// RestoreMode return terminal into standard mode
	RestoreMode(*term.State) error
	// Fd return file descriptor reference
	Fd() int
	//File return file descriptor
	File() *os.File
	//Get height of terminal window
	GetHeight() (int, error)
	//Get widht of terminal window
	GetWidth() (int, error)
}

type TerminalInput struct {
	file *os.File
}

func NewTerminalInput(file *os.File) *TerminalInput {
	return &TerminalInput{file: file}
}

func (ti *TerminalInput) ReadLine(prompt string) (string, error) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(ti.file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

func (ti *TerminalInput) ReadKey() (byte, error) {
	b := make([]byte, 3)
	n, err := ti.file.Read(b)
	if err != nil {
		return 0, err
	}
	if n == 3 {
		return b[2], nil
	}
	return b[0], nil
}

func (ti *TerminalInput) Fd() int {
	return int(ti.file.Fd())
}

func (ti *TerminalInput) File() *os.File {
	return ti.file
}

func (ti *TerminalInput) SetRawMode() (*term.State, error) {
	fd := int(ti.file.Fd())
	return term.MakeRaw(fd)
}

func (ti *TerminalInput) RestoreMode(state *term.State) error {
	fd := int(ti.file.Fd())
	return term.Restore(fd, state)
}

func (ti TerminalInput) GetHeight() (int, error) {
	_, h, err := term.GetSize(ti.Fd())
	if err != nil {
		return 0, err
	}
	return h, nil
}

func (ti TerminalInput) GetWidth() (int, error) {
	w, _, err := term.GetSize(ti.Fd())
	if err != nil {
		return 0, err
	}
	return w, nil
}
