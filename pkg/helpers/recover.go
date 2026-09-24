package helpers

import "fmt"

const error_msg = "panic captured by helpers.Recover: %v"

// Recover converts a panic into an error.
//
// When fn is non-nil, fn is invoked under a recovery guard and its result is
// returned. When fn is nil, Recover recovers panics of the current goroutine
// and is intended to be used as a deferred call, e.g.
// defer helpers.Recover(nil, &err).
//
// When a panic occurs, the resulting error is written to *err when err is
// non-nil and returned otherwise. When no panic occurs, Recover returns a nil
// error (or fn's result when fn is non-nil).
func Recover(fn func() error, err *error) (ret error) {
	if fn == nil {
		if r := recover(); r != nil {
			e := fmt.Errorf(error_msg, r)
			if err != nil {
				*err = e
			}
			if err == nil {
				panic(r)
			}
			return e
		}
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			e := fmt.Errorf(error_msg, r)
			if err != nil {
				*err = e
			}
			ret = e
		}
	}()
	return fn()
}
