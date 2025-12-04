package main

import (
	"cliTodoist/cli"
	"cliTodoist/internal/colors"
	"cliTodoist/internal/input"
	"cliTodoist/internal/util"
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

func main() {
	t := os.NewFile(uintptr(syscall.Stdin), "/dev/tty")

	defer func() {
		t.Close()
	}()

	input := input.NewTerminalInput(t)
	h, err := input.GetHeight()
	if h < util.MinWindowHeight {
		err = errors.New("To work with this app your terminal should be minimum " + strconv.Itoa(util.MinWindowHeight) + " lines height!")
	}
	if err != nil {
		fmt.Println(colors.Red+"\rAn error occurred: ", err.Error()+colors.Reset)
		return
	}
	cli := cli.CLI{Input: input}
	errExit := cli.Run()
	if errExit != nil {
		fmt.Println("\r", colors.Green+errExit.Error())
		return
	}

}
