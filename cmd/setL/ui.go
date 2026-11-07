package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	reader   = bufio.NewReader(os.Stdin)
	DemoMode = false
)

func init() {
	val := strings.ToLower(os.Getenv("DEMO"))
	if val == "1" || val == "true" || val == "yes" {
		DemoMode = true
	}
	for _, arg := range os.Args[1:] {
		if arg == "--demo" || arg == "-demo" || arg == "demo" {
			DemoMode = true
		}
	}
}

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[1;31m"
	ColorGreen  = "\033[1;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[1;34m"
	ColorCyan   = "\033[1;36m"
	ColorWhite  = "\033[1;37m"
	ColorGray   = "\033[0;90m"
)

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Print("\033[H\033[2J")
	}
}

func ReadLine(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func Pause() {
	fmt.Print("\n" + ColorGray + "Нажмите Enter для продолжения..." + ColorReset)
	_, _ = reader.ReadString('\n')
}

