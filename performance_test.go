package lo

import (
	"testing"

	upstream "github.com/samber/lo"
)

var benchmarkInput = Times(1024, func(index int) int { return index })
var benchmarkOutput []int

func selectedSquare(value, _ int) (int, bool) { return value * value, value%2 == 0 }

func BenchmarkCoreFusedInto(b *testing.B) {
	buffer := make([]int, 0, len(benchmarkInput))
	b.ReportAllocs()
	for b.Loop() {
		benchmarkOutput = FilterMapInto(buffer[:0], benchmarkInput, selectedSquare)
	}
}

func BenchmarkUpstreamMapFilterPipeline(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		mapped := upstream.Map(benchmarkInput, func(value, _ int) int { return value * value })
		benchmarkOutput = upstream.Filter(mapped, func(_ int, index int) bool { return benchmarkInput[index]%2 == 0 })
	}
}

func BenchmarkUpstreamFilterMap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkOutput = upstream.FilterMap(benchmarkInput, selectedSquare)
	}
}

func TestFusedIntoAllocationBudget(t *testing.T) {
	buffer := make([]int, 0, len(benchmarkInput))
	core := testing.AllocsPerRun(1000, func() { benchmarkOutput = FilterMapInto(buffer[:0], benchmarkInput, selectedSquare) })
	upstreamPipeline := testing.AllocsPerRun(1000, func() {
		mapped := upstream.Map(benchmarkInput, func(value, _ int) int { return value * value })
		benchmarkOutput = upstream.Filter(mapped, func(_ int, index int) bool { return benchmarkInput[index]%2 == 0 })
	})
	if core != 0 {
		t.Fatalf("core allocations = %.1f, want 0", core)
	}
	if upstreamPipeline < 2 || core*2 > upstreamPipeline {
		t.Fatalf("allocation reduction not at least 50%%: core %.1f, upstream %.1f", core, upstreamPipeline)
	}
}

func TestCompatibilityAllocationParity(t *testing.T) {
	checks := []struct {
		name   string
		ours   func()
		theirs func()
	}{
		{"Map", func() { benchmarkOutput = Map(benchmarkInput, func(v, _ int) int { return v + 1 }) }, func() { benchmarkOutput = upstream.Map(benchmarkInput, func(v, _ int) int { return v + 1 }) }},
		{"Filter", func() { benchmarkOutput = Filter(benchmarkInput, func(v, _ int) bool { return v%2 == 0 }) }, func() { benchmarkOutput = upstream.Filter(benchmarkInput, func(v, _ int) bool { return v%2 == 0 }) }},
		{"FilterMap", func() { benchmarkOutput = FilterMap(benchmarkInput, selectedSquare) }, func() { benchmarkOutput = upstream.FilterMap(benchmarkInput, selectedSquare) }},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			ours := testing.AllocsPerRun(1000, check.ours)
			theirs := testing.AllocsPerRun(1000, check.theirs)
			if ours > theirs {
				t.Fatalf("allocations = %.1f, upstream %.1f", ours, theirs)
			}
		})
	}
}

func BenchmarkCompatibilityMap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkOutput = Map(benchmarkInput, func(value, _ int) int { return value + 1 })
	}
}

func BenchmarkUpstreamMap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchmarkOutput = upstream.Map(benchmarkInput, func(value, _ int) int { return value + 1 })
	}
}
