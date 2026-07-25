package lo_test

import (
	"errors"
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

func TestErrorsBatch(t *testing.T) {
	sentinel := errors.New("boom")

	checkProp(t, "Try/no-error", func(shouldErr bool) bool {
		f := func() error {
			if shouldErr {
				return sentinel
			}
			return nil
		}
		return forge.Try(f) == upstream.Try(f)
	})
	checkProp(t, "TryOr", func(shouldErr bool, val, fb int) bool {
		f := func() (int, error) {
			if shouldErr {
				return 0, sentinel
			}
			return val, nil
		}
		fr, fok := forge.TryOr(f, fb)
		ur, uok := upstream.TryOr(f, fb)
		return fr == ur && fok == uok
	})
	checkProp(t, "Validate", func(ok bool, n int) bool {
		fe := forge.Validate(ok, "n=%d", n)
		ue := upstream.Validate(ok, "n=%d", n)
		if (fe == nil) != (ue == nil) {
			return false
		}
		return fe == nil || fe.Error() == ue.Error()
	})

	// Must: both panic on error, both return the value otherwise.
	if !didPanic(func() { forge.Must(0, sentinel) }) || !didPanic(func() { upstream.Must(0, sentinel) }) {
		t.Error("Must should panic on error in both")
	}
	if forge.Must(42, nil) != upstream.Must(42, nil) {
		t.Error("Must value mismatch")
	}
	if didPanic(func() { forge.Must(42, nil) }) {
		t.Error("Must should not panic on nil error")
	}
}
