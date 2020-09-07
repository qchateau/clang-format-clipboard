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

func Find(str []string, x string) int {
	for i, n := range str {
		if x == n {
			return i
		}
	}
	return len(str)
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func CountPrefix(str, chars string) int {
	for i, c := range str {
		if !strings.ContainsRune(chars, c) {
			return i
		}
	}
	return len(str)
}

func IdentifyPrefix(s string) (string, string) {
	prefix_len := CountPrefix(s, whitespace)
	indent_start_pos := strings.LastIndex(s[:prefix_len], "\n") + 1
	return s[:indent_start_pos], s[indent_start_pos:prefix_len]
}

func TrimTrailingWhitespaces(s string) string {
	lines := strings.Split(s, "\n")
	lines = Map(lines, func(line string) string { return strings.TrimRight(line, whitespace) })
	return strings.Join(lines, "\n")
}

func main() {
	flagset := flag.NewFlagSet("main", flag.ContinueOnError)
	executable := flagset.String("executable", "clang-format", "clang-format excutable name")

	args_sep := Find(os.Args, "--")
	args_to_parse := os.Args[1:args_sep]
	args_sep = Min(args_sep, len(os.Args)-1)
	clang_format_flags := os.Args[args_sep+1:]
	flagset.Parse(args_to_parse)

	clipboard_content, err := clipboard.ReadAll()
	if err != nil {
		log.Fatalf("Error reading the clipboard: %v\n", err)
		os.Exit(1)
	}

	clipboard_content = TrimTrailingWhitespaces(clipboard_content)
	prefix, indent := IdentifyPrefix(clipboard_content)
	clipboard_content = clipboard_content[len(prefix):]

	log.Printf("Formatting %v chars with %v %v\n", len(clipboard_content), *executable, clang_format_flags)
	log.Printf("Indent is '%v'\n", indent)

	cmd := exec.Command(*executable)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Error with stdin: %v\n", err)
		os.Exit(1)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, clipboard_content)
	}()

	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error with clang-format: %v\n", err)
		switch spec_err := err.(type) {
		case *exec.ExitError:
			log.Fatalln(spec_err)
		}
	}

	out_str := string(out)
	log.Printf("Formatted:\n%v\n", out_str)
	lines := strings.SplitAfter(out_str, "\n")
	lines = Map(lines, func(s string) string {
		if len(s) > 0 {
			return indent + s
		} else {
			return s
		}
	})
	log.Printf("Indented:\n%v\n", strings.Join(lines, ""))

	err = clipboard.WriteAll(prefix + strings.Join(lines, ""))
	if err != nil {
		log.Fatalf("Error writing the clipboard: %v\n", err)
		os.Exit(1)
	}
}
