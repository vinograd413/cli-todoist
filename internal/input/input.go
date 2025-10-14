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
	Fd() int
	File() *os.File
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
