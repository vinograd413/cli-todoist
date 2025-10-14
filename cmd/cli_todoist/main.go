package main

import (
	"cliTodoist/cli"
	"cliTodoist/internal/input"
	"fmt"
	"os"
	"syscall"
)

func main() {
	t := os.NewFile(uintptr(syscall.Stdin), "/dev/tty")

	defer func() {
		t.Close()
	}()

	input := input.NewTerminalInput(t)
	cli := cli.CLI{Input: input}
	err := cli.Run()
	if err != nil {
		fmt.Println("An unexpected error occurred:", err)
		return
	}
}
