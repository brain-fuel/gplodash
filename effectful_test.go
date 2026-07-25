package lo_test

import (
	"testing"
	"time"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

// The channel/concurrency/retry/time helpers are effectful (goroutines, timers)
// so these are behavioral smoke tests against upstream rather than property
// laws: same observable result on a deterministic scenario.

func TestConcurrencyAsync(t *testing.T) {
	f := forge.Async(func() int { return 42 })
	u := upstream.Async(func() int { return 42 })
	if <-f != 42 || <-u != 42 {
		t.Error("Async result mismatch")
	}
	// Async0..Async6 arity family exists and completes.
	<-forge.Async0(func() {})
	r := forge.Async2(func() (int, string) { return 1, "a" })
	got := <-r
	if got.A != 1 || got.B != "a" {
		t.Error("Async2 tuple mismatch")
	}
}

func TestRetryAttempt(t *testing.T) {
	fail := func(i int) error {
		if i < 2 {
			return errDummy
		}
		return nil
	}
	fi, ferr := forge.Attempt(5, fail)
	ui, uerr := upstream.Attempt(5, fail)
	if ui != fi || (uerr == nil) != (ferr == nil) {
		t.Errorf("Attempt parity: forge (%d,%v) vs upstream (%d,%v)", fi, ferr, ui, uerr)
	}
	if ferr != nil {
		t.Errorf("Attempt should succeed by iter 2: %v", ferr)
	}
}

func TestTimeDuration(t *testing.T) {
	d := forge.Duration(func() { time.Sleep(2 * time.Millisecond) })
	if d < time.Millisecond {
		t.Errorf("Duration too small: %v", d)
	}
	// Duration1 threads a value through.
	got, dur := forge.Duration1(func() int { return 7 })
	if got != 7 || dur < 0 {
		t.Errorf("Duration1: got=%d dur=%v", got, dur)
	}
}

func TestChannelSliceToChannel(t *testing.T) {
	ch := forge.SliceToChannel(2, []int{1, 2, 3})
	got := forge.ChannelToSlice(ch)
	if len(got) != 3 || got[0] != 1 || got[2] != 3 {
		t.Errorf("SliceToChannel/ChannelToSlice roundtrip: %v", got)
	}
}

var errDummy = &dummyErr{}

type dummyErr struct{}

func (*dummyErr) Error() string { return "dummy" }
