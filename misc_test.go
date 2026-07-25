package lo_test

import (
	"reflect"
	"testing"

	upstream "github.com/samber/lo"
	forge "goforge.dev/gplodash"
)

// Property laws for the type-manipulation, condition, and func batches.

func TestTypeManipulationProp(t *testing.T) {
	checkProp(t, "FromPtrOr", func(hasVal bool, v, fb int) bool {
		var p *int
		if hasVal {
			p = &v
		}
		return forge.FromPtrOr(p, fb) == upstream.FromPtrOr(p, fb)
	})
	checkProp(t, "EmptyableToPtr", func(v int) bool {
		fp := forge.EmptyableToPtr(v)
		up := upstream.EmptyableToPtr(v)
		if (fp == nil) != (up == nil) {
			return false
		}
		return fp == nil || *fp == *up
	})
	checkProp(t, "IsEmpty/IsNotEmpty", func(v int) bool {
		return forge.IsEmpty(v) == upstream.IsEmpty(v) &&
			forge.IsNotEmpty(v) == upstream.IsNotEmpty(v)
	})
	checkProp(t, "Coalesce", func(vs []int) bool {
		fv, fok := forge.Coalesce(vs...)
		uv, uok := upstream.Coalesce(vs...)
		return fv == uv && fok == uok
	})
	checkProp(t, "CoalesceSliceOrEmpty", func(a, b []int) bool {
		return reflect.DeepEqual(
			forge.CoalesceSliceOrEmpty(a, b),
			upstream.CoalesceSliceOrEmpty(a, b),
		)
	})
	checkProp(t, "FromAnySlice", func(xs []int) bool {
		anys := forge.ToAnySlice(xs)
		fo, fok := forge.FromAnySlice[int](anys)
		uo, uok := upstream.FromAnySlice[int](anys)
		return fok == uok && reflect.DeepEqual(fo, uo)
	})
}

func TestConditionProp(t *testing.T) {
	checkProp(t, "Ternary", func(c bool, a, b int) bool {
		return forge.Ternary(c, a, b) == upstream.Ternary(c, a, b)
	})
	checkProp(t, "Switch", func(pred int, r1, r2, def int) bool {
		f := forge.Switch[int, int](pred).Case(1, r1).Case(2, r2).Default(def)
		u := upstream.Switch[int, int](pred).Case(1, r1).Case(2, r2).Default(def)
		return f == u
	})
	checkProp(t, "If/ElseIf/Else", func(c1, c2 bool, a, b, c int) bool {
		f := forge.If(c1, a).ElseIf(c2, b).Else(c)
		u := upstream.If(c1, a).ElseIf(c2, b).Else(c)
		return f == u
	})
}

func TestFuncProp(t *testing.T) {
	checkProp(t, "Partial", func(a, b int) bool {
		add := func(x, y int) int { return x + y }
		return forge.Partial(add, a)(b) == upstream.Partial(add, a)(b)
	})
	checkProp(t, "Partial2", func(a, b, c int) bool {
		f := func(x, y, z int) int { return x*100 + y*10 + z }
		return forge.Partial2(f, a)(b, c) == upstream.Partial2(f, a)(b, c)
	})
}
