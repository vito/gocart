package main

import (
	"fmt"
	"os"
	"strings"
)

func bold(str string) string {
	return "\x1b[1m" + str + "\x1b[0m"
}

func red(str string) string {
	return "\x1b[31m" + str + "\x1b[0m"
}

func green(str string) string {
	return "\x1b[32m" + str + "\x1b[0m"
}

func cyan(str string) string {
	return "\x1b[36m" + str + "\x1b[0m"
}

func padding(size int) string {
	return strings.Repeat(" ", size)
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, red(message))
	os.Exit(1)
}
