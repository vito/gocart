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

func fatal(message interface{}) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

func indent(level int, str string) string {
	indented := ""

	for _, line := range strings.Split(strings.TrimRight(str, " \n"), "\n") {
		indented = indented + padding(level*2) + line + "\n"
	}

	return indented[0 : len(indented)-1]
}
