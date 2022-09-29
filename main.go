package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

const helpMessage = "\nYou are using a custom version of ls where -a -> -A and vice-versa."

func main() {
	err := RunLsCommand(os.Args[1:], os.Stdout, os.Stdin, os.Stderr)

	// ls takes care of printing an error message (ex: "no such or file directory).
	// We just have to transmit the status code.
	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}
}

func RunLsCommand(args []string, stdout io.ReadWriter, stderr io.Writer, stdin io.Reader) error {
	for _, arg := range args {
		if arg == "--help" {
			buf := &bytes.Buffer{}
			cmd := exec.Command("ls", "--help")
			cmd.Stderr = stderr
			cmd.Stdin = stdin
			cmd.Stdout = buf
			if err := cmd.Run(); err != nil {
				return err
			}

			updatedOut := bytes.Replace(
				buf.Bytes(),
				[]byte("-a, --all                  do not ignore entries starting with .\n  -A, --almost-all           do not list implied . and .."),
				[]byte("-a, --almost-all           do not list implied . and ..\n  -A, --all                  do not ignore entries starting with .\n"),
				1,
			)

			_, err := stdout.Write(updatedOut)
			return err
		}
	}

	return NewLsCommand(args, stdout, stderr, stdin).Run()
}

func NewLsCommand(args []string, stdout io.Writer, stderr io.Writer, stdin io.Reader) *exec.Cmd {
	ls, _ := exec.LookPath("ls")

	// Go does not support lookaheads and probably won't support it.
	// If supports lands, this code should be refactored to use it.
	// https://github.com/golang/go/issues/18868
	for i, arg := range args {
		if arg[0] != '-' || !strings.ContainsAny(arg, "aA") {
			continue
		}

		// We let ls handle any conflicting options.
		args[i] = strings.NewReplacer("a", "A", "A", "a").Replace(arg)
		if arg[0] != '-' || !strings.ContainsAny(arg, "aA") {
			continue
		}

		// We let ls handle any conflicting options.
		args[i] = strings.NewReplacer("a", "A", "A", "a").Replace(arg)
	}

	return &exec.Cmd{
		Path:   ls,
		Args:   append([]string{"ls"}, args...),
		Stdout: stdout,
		Stdin:  stdin,
		Stderr: stderr,
	}
}
