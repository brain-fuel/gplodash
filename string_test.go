package lo_test

import (
	"reflect"
	"testing"
	"testing/quick"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

// Property laws: the pure string batch must equal pinned samber/lo v1.53.0 for
// arbitrary inputs (testing/quick generates and shrinks counterexamples).

func TestStringWordsProp(t *testing.T) {
	checkProp(t, "RuneLength", func(s string) bool {
		return forge.RuneLength(s) == upstream.RuneLength(s)
	})
	checkProp(t, "Words", func(s string) bool {
		return reflect.DeepEqual(forge.Words(s), upstream.Words(s))
	})
	checkProp(t, "KebabCase", func(s string) bool {
		return forge.KebabCase(s) == upstream.KebabCase(s)
	})
	checkProp(t, "SnakeCase", func(s string) bool {
		return forge.SnakeCase(s) == upstream.SnakeCase(s)
	})
	checkProp(t, "Ellipsis", func(s string, n int8) bool {
		return forge.Ellipsis(s, int(n)) == upstream.Ellipsis(s, int(n))
	})
}

func TestChunkStringProp(t *testing.T) {
	// size is forced positive; the panic-on-nonpositive contract is checked
	// separately below.
	checkProp(t, "ChunkString", func(s string, size uint8) bool {
		n := int(size)%16 + 1
		return reflect.DeepEqual(forge.ChunkString(s, n), upstream.ChunkString(s, n))
	})
	if err := quick.Check(func(s string, size uint8) bool {
		n := int(size)%16 + 1
		return reflect.DeepEqual(forge.ChunkString(s, n), upstream.ChunkString(s, n))
	}, &quick.Config{MaxCount: 2000}); err != nil {
		t.Error(err)
	}
}

func TestChunkStringPanicsParity(t *testing.T) {
	for _, size := range []int{0, -1, -5} {
		assertPanicsSame(t, size)
	}
}

func assertPanicsSame(t *testing.T, size int) {
	t.Helper()
	forgePanic := didPanic(func() { forge.ChunkString("abc", size) })
	upPanic := didPanic(func() { upstream.ChunkString("abc", size) })
	if forgePanic != upPanic {
		t.Errorf("ChunkString(size=%d): forge panic=%v, upstream panic=%v", size, forgePanic, upPanic)
	}
}

func didPanic(f func()) (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	f()
	return false
}
