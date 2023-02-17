package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// ref: https://github.com/spf13/cobra/blob/ad6db7f8f6e485f55b1561df9276fa4d8e278bde/command_test.go#L74-L78
func checkStringContains(t *testing.T, got, expected string) {
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}

// ref: https://github.com/spf13/cobra/blob/ad6db7f8f6e485f55b1561df9276fa4d8e278bde/command_test.go#L32-L35
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

// ref: https://github.com/spf13/cobra/blob/ad6db7f8f6e485f55b1561df9276fa4d8e278bde/command_test.go#L48-L57
func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

var u string

func TestMain(m *testing.M) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u = "ws" + strings.TrimPrefix(s.URL, "http")
	code := m.Run()
	os.Exit(code)
}

/*
===Private Key (Do NOT share to anyone)===
0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fba
nsec1q9t3skr57525l2gpxn0csuvy5x8z68gcpplmj4jn7spxnpxfr7aqn88ve4

===Public Key (Share with your friends!)===
a63e5a0a5747b86ad6fd1bb1a04cce0fa3a718cff3dcd90dca4f3e968eac049d
npub15cl95zjhg7ux44harwc6qnxwp736wxx070wdjrw2fulfdr4vqjwsu7n80u
*/

func TestPub(t *testing.T) {
	type TestCase struct {
		Name        string
		FlagRelay   string
		FlagSk      string
		FlagContent string
		ExpectErr   string
	}

	cases := []TestCase{
		{
			Name:        "success",
			FlagRelay:   u,
			FlagSk:      "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fba",
			FlagContent: "content",
		},
		{
			Name:      "non-exist http endpoint",
			FlagRelay: "https://localhost:3939",
			FlagSk:    "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fba",
			ExpectErr: "error opening websocket to 'wss://localhost:3939': dial tcp [::1]:3939: connect: connection refused",
		},
		{
			Name:      "non-exist websocket endpoint",
			FlagRelay: "wss://localhost:3939",
			FlagSk:    "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fba",
			ExpectErr: "error opening websocket to 'wss://localhost:3939': dial tcp [::1]:3939: connect: connection refused",
		},
		{
			Name:      "invalid format relay",
			FlagRelay: "foo",
			FlagSk:    "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fba",
			ExpectErr: "error opening websocket to 'wss://foo': dial tcp: lookup foo: no such host",
		},
		{
			Name:      "empty sk",
			FlagRelay: u,
			FlagSk:    "", // zero
			ExpectErr: "invalid hash length",
		},
		{
			Name:      "invalid length sk",
			FlagRelay: u,
			FlagSk:    "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fb", // less than 32-bytes
			ExpectErr: "invalid hash length",
		},
		{
			Name:      "invalid length sk",
			FlagRelay: u,
			FlagSk:    "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fbaa", // more than 32-bytes
			ExpectErr: "invalid hash length",
		},
		{
			Name:      "invalid character in sk",
			FlagRelay: u,
			FlagSk:    "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fbx", // invalid character: x
			ExpectErr: "encoding/hex: invalid byte:",
		},
	}

	for _, c := range cases {
		rootCmd := initRootCmd()
		initPubCmd(rootCmd)
		t.Run(c.Name, func(t *testing.T) {
			err := pub(c.FlagRelay, c.FlagSk, c.FlagContent)

			if len(c.ExpectErr) > 0 {
				if err != nil {
					checkStringContains(t, err.Error(), c.ExpectErr)
				} else {
					t.Error("error should be raised")
				}
			} else if err != nil {
				t.Error("unexpected error")
			}
		})
	}
}

func TestPubCmd(t *testing.T) {
	type TestCase struct {
		Name      string
		Args      []string
		ExpectErr string
	}
	cases := []TestCase{
		{
			Name:      "should raise an error for an argument",
			Args:      []string{"pub", "--relay"},
			ExpectErr: "flag needs an argument: --relay",
		},
		{
			Name:      "should raise an error for an argument",
			Args:      []string{"pub", "--secret"},
			ExpectErr: "flag needs an argument: --secret",
		},
		{
			Name:      "should raise an error for an argument",
			Args:      []string{"pub", "--content"},
			ExpectErr: "flag needs an argument: --content",
		},
		{
			Name:      "should raise an error for insufficient flags: --relay,--secret,--content",
			Args:      []string{"pub"},
			ExpectErr: "required flag(s) \"content\", \"relay\", \"secret\" not set",
		},
		{
			Name:      "should raise an error for insufficient flag: --secret",
			Args:      []string{"pub", "--relay", "test", "--content", "text"},
			ExpectErr: "required flag(s) \"secret\" not set",
		},
		{
			Name:      "should raise an error for insufficient flag: --relay",
			Args:      []string{"pub", "--secret", "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fba", "--content", "text"},
			ExpectErr: "required flag(s) \"relay\" not set",
		},
		{
			Name:      "should raise an error for insufficient flag: --content",
			Args:      []string{"pub", "--relay", "test", "--secret", "0157185874f5154fa90134df887184a18e2d1d18087fb95653f4026984c91fba"},
			ExpectErr: "required flag(s) \"content\" not set",
		},
	}

	for _, c := range cases {
		rootCmd := initRootCmd()
		initPubCmd(rootCmd)
		_, err := executeCommand(rootCmd, c.Args...)
		t.Run(c.Name, func(t *testing.T) {
			if len(c.ExpectErr) > 0 {
				if err != nil {
					checkStringContains(t, err.Error(), c.ExpectErr)
				} else {
					t.Error("error should be raised")
				}
			} else if err != nil {
				t.Error("unexpected error")
			}
		})
	}
}
