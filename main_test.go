package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

type Testcase struct {
	name string
	test func(*testing.T, io.Reader)
	args []string
}

func TestRun(t *testing.T) {
	tt := []Testcase{
		{
			name: "test help flag hijacking",
			args: []string{"--help"},
			test: func(t *testing.T, r io.Reader) {
				// yes of course it was done on purpose, why would you think otherwise?
				buf := make([]byte, 69)

				if _, err := io.ReadFull(r, buf); err != nil || bytes.HasSuffix(buf, []byte(helpMessage)) {
					t.Fatalf("Help message is not rendering correctly: %s (err: %s)", buf, err)
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("Can not fake standard io.")
			}

			if err = RunLsCommand(tc.args, w, nil, w); err != nil {
				t.Fatalf("Error running command: %s", err)
			}

			tc.test(t, r)
		})
	}
}

func TestNewLsCommand(t *testing.T) {
	tt := [][2]string{
		{"-âaAA31aaäädgdaA", "-âAaa31AAäädgdAa"},
		{"-a", "-A"},
		{"-A", "-a"},
		{"-l", "-l"},
		{"adir", "adir"},
		{"dir", "dir"},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf("%s to -> %s", tc[0], tc[1]), func(t *testing.T) {

			cmd := NewLsCommand([]string{tc[0]}, nil, nil, nil)

			if cmd.Args[1] != tc[1] {
				t.Errorf("Invalid case switching (%s should become %s, got %s)", tc[0], tc[1], cmd.Args[1])
			}
		})
	}
}
