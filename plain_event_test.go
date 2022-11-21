package fiberEvents

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestPlainEventStream(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	log := NewStream(os.Stdout, EventLevelLog)
	log.Debug("my debug message", "my debug payload")
	w.Close()

	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	expected := fmt.Sprintf("%s\n", `{"event_type":"debug","level":1,"message":"my debug message","payload":"my debug payload"}`)
	if string(out) != expected {
		spew.Dump(string(out), expected)
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}
