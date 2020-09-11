package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

const whitespace = " \t\n\r"

func find(str []string, x string) int {
	for i, n := range str {
		if x == n {
			return i
		}
	}
	return len(str)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func apply(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func countPrefix(str, chars string) int {
	for i, c := range str {
		if !strings.ContainsRune(chars, c) {
			return i
		}
	}
	return len(str)
}

func identifyPrefix(s string) (string, string) {
	prefixLen := countPrefix(s, whitespace)
	indentStartPos := strings.LastIndex(s[:prefixLen], "\n") + 1
	return s[:indentStartPos], s[indentStartPos:prefixLen]
}

func trimTrailingWhitespaces(s string) string {
	lines := strings.Split(s, "\n")
	lines = apply(lines, func(line string) string { return strings.TrimRight(line, whitespace) })
	return strings.Join(lines, "\n")
}

func main() {
	flagset := flag.NewFlagSet("main", flag.ExitOnError)
	executable := flagset.String("executable", "clang-format", "clang-format excutable name")
	keepIndentation := flagset.Bool("keep-indentation", true, "keep the original indentation")
	stripLeadingNewlines := flagset.Bool("strip-leading-newlines", false, "strip newlines before the first line of code")

	argsSep := find(os.Args, "--")
	argsToParse := os.Args[1:argsSep]
	argsSep = min(argsSep, len(os.Args)-1)
	clangFormatFlags := os.Args[argsSep+1:]
	flagset.Parse(argsToParse)

	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		log.Fatalf("Error reading the clipboard: %v\n", err)
		os.Exit(1)
	}

	clipboardContent = trimTrailingWhitespaces(clipboardContent)
	prefix, indent := identifyPrefix(clipboardContent)
	clipboardContent = clipboardContent[len(prefix):]

	log.Printf("Formatting %v chars with %v %v\n", len(clipboardContent), *executable, clangFormatFlags)

	cmd := exec.Command(*executable)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error with stdin: %v\n", err)
		os.Exit(1)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, clipboardContent)
	}()

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error with clang-format: %v\n", err)
		switch specErr := err.(type) {
		case *exec.ExitError:
			log.Fatalln(specErr)
		}
	}

	outStr := string(out)
	lines := strings.SplitAfter(outStr, "\n")
	if *keepIndentation {
		lines = apply(lines, func(s string) string {
			if len(s) > 0 {
				return indent + s
			}
			return s
		})
	}

	if *stripLeadingNewlines {
		outStr = strings.Join(lines, "")
	} else {
		outStr = prefix + strings.Join(lines, "")
	}
	err = clipboard.WriteAll(outStr)
	if err != nil {
		log.Fatalf("Error writing the clipboard: %v\n", err)
		os.Exit(1)
	}
}
