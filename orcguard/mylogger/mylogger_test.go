package mylogger

import (
	"testing"
)

func TestConstLevel(t *testing.T) {
	t.Logf("%v %T\n", DebugLevel, DebugLevel)
	t.Logf("%v %T\n", FatalLevel, FatalLevel)
}
