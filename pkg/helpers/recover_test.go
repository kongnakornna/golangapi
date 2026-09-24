package helpers

import (
	"fmt"
	"strings"
	"testing"
)

func TestRecoverPanicReturnsError(t *testing.T) {
	err := Recover(func() error { panic("boom") }, nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected error to contain %q, got %s", "boom", err)
	}
}

func TestRecoverPanicWritesAndReturnsError(t *testing.T) {
	var err error
	got := Recover(func() error { panic(fmt.Errorf("kaboom")) }, &err)
	if err == nil {
		t.Fatalf("expected err to be written, got nil")
	}
	if got == nil {
		t.Fatalf("expected return error, got nil")
	}
	if got != err {
		t.Fatalf("expected return error to match written error, got %s vs %s", got, err)
	}
}

func TestRecoverDeferStyle(t *testing.T) {
	var err error
	func() {
		defer Recover(nil, &err)
		panic("deferred boom")
	}()
	if err == nil {
		t.Fatalf("expected deferred panic to be captured, got nil")
	}
	if !strings.Contains(err.Error(), "deferred boom") {
		t.Fatalf("expected error to contain %q, got %s", "deferred boom", err)
	}
}

func TestRecoverNoPanicReturnsNil(t *testing.T) {
	if got := Recover(nil, nil); got != nil {
		t.Fatalf("expected nil, got %s", got)
	}
	if got := Recover(func() error { return nil }, nil); got != nil {
		t.Fatalf("expected nil, got %s", got)
	}
}

func TestRecoverReturnsFnResult(t *testing.T) {
	sentinel := fmt.Errorf("sentinel")
	got := Recover(func() error { return sentinel }, nil)
	if got != sentinel {
		t.Fatalf("expected %s, got %s", sentinel, got)
	}
}
