package main

import (
	"strings"
	"testing"
)

func TestNewUUID(t *testing.T) {
	a := newUUID()
	b := newUUID()

	if a == "" {
		t.Error("UUID generated empty")
		t.Fail()
	}

	if strings.TrimSpace(a) == "" {
		t.Error("UUID generated whitespaces")
		t.Fail()
	}

	if a == b {
		t.Error("UUID generated unique values")
		t.Fail()
	}
}
