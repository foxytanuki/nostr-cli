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
