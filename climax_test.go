package climax

import (
	"bytes"
	"os"
	"testing"
)

var output bytes.Buffer

func init() {
	outputDevice = &output
	errorDevice = &output
}

func setArguments(args ...string) {
	os.Args = append([]string{"test"}, args...)
}

func mustPanic(t *testing.T, text string, fn func()) {
	defer func() {
		state := recover()
		if state == nil {
			t.Errorf(`case "%s" didn't panic`, text)
		}
	}()

	fn()
}

func TestNew(t *testing.T) {
	a := New("smth")
	if a.Name != "smth" {
		t.Errorf("actual app name (%s) doesn't match passed (smth)", a.Name)
	}
}
