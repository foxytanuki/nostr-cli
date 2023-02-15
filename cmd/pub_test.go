package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// ref: https://github.com/spf13/cobra/blob/ad6db7f8f6e485f55b1561df9276fa4d8e278bde/command_test.go#L74-L78
func checkStringContains(t *testing.T, got, expected string) {
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}

func TestPub(t *testing.T) {
	type TestCase struct {
		Name        string
		FlagName    string
		FlagValue   string
		ExpectErr   bool
		ExpectedMsg string
	}

	cases := []TestCase{
		{
			Name:      "success",
			FlagName:  "relay",
			FlagValue: "wss://foo.io",
			ExpectErr: false,
		},
		{
			Name:      "should read an argument",
			FlagName:  "relay",
			FlagValue: "",
			ExpectErr: true,
		},
	}

	initPubCmd()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			err := pub(rootCmd, nil)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestPubCmd(t *testing.T) {
	initPubCmd()

	b := bytes.NewBufferString("")
	rootCmd.SetOut(b)

	type TestCase struct {
		Name        string
		Args        []string
		ExpectErr   bool
		ExpectedMsg string
	}

	cases := []TestCase{
		{
			Name:        "should raise an error for an argument",
			Args:        []string{"pub", "--relay"},
			ExpectErr:   true,
			ExpectedMsg: "flag needs an argument: --relay",
		},
		{
			Name:        "should raise an error for insufficient flags",
			Args:        []string{"pub"},
			ExpectErr:   true,
			ExpectedMsg: "required flag(s) \"relay\" not set",
		},
	}

	for _, c := range cases {
		rootCmd.SetArgs(c.Args)
		err := rootCmd.Execute()
		t.Run(c.Name, func(t *testing.T) {
			if c.ExpectErr {
				if err != nil {
					checkStringContains(t, err.Error(), c.ExpectedMsg)
				} else {
					t.Error("error should be raised")
				}
			} else if err != nil {
				t.Error("unexpected error")
			}
		})
		// out, _ := io.ReadAll(b)
		// fmt.Println(out)
	}
}
