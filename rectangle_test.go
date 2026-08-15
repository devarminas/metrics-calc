package main

import (
	"math"
	"testing"
)

const relTolerance = 1e-9

func closeEnough(got, want float64) bool {
	if got == want {
		return true
	}
	return math.Abs(got-want) <= relTolerance*math.Max(math.Abs(got), math.Abs(want))
}

func TestNewRectangle(t *testing.T) {
	tests := []struct {
		name          string
		width, height float64
		want          rectangle
	}{
		{"square", 10, 10, rectangle{width: 10, height: 10}},
		{"unit square", 1, 1, rectangle{width: 1, height: 1}},
		{"wider than tall", 3, 4, rectangle{width: 3, height: 4}},
		{"fractional sides", 2.5, 0.5, rectangle{width: 2.5, height: 0.5}},
		{
			"smallest positive sides", math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64,
			rectangle{width: math.SmallestNonzeroFloat64, height: math.SmallestNonzeroFloat64},
		},
		{
			"largest finite sides", math.MaxFloat64, math.MaxFloat64,
			rectangle{width: math.MaxFloat64, height: math.MaxFloat64},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRectangle(tt.width, tt.height)
			if err != nil {
				t.Fatalf("newRectangle(%v, %v) error = %v, want nil", tt.width, tt.height, err)
			}
			if got != tt.want {
				t.Errorf("newRectangle(%v, %v) = %+v, want %+v", tt.width, tt.height, got, tt.want)
			}
		})
	}
}

func TestNewRectangle_invalidDimensions(t *testing.T) {
	tests := []struct {
		name          string
		width, height float64
	}{
		{"zero width", 0, 4},
		{"zero height", 3, 0},
		{"both sides zero", 0, 0},
		{"negative width", -3, 4},
		{"negative height", 3, -4},
		{"both sides negative", -3, -4},
		{"NaN width", math.NaN(), 4},
		{"NaN height", 3, math.NaN()},
		{"positive infinite width", math.Inf(1), 4},
		{"negative infinite height", 3, math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRectangle(tt.width, tt.height)
			if err == nil {
				t.Fatalf("newRectangle(%v, %v) = %+v, want error", tt.width, tt.height, got)
			}
			if got != (rectangle{}) {
				t.Errorf("newRectangle(%v, %v) = %+v, want zero rectangle alongside error",
					tt.width, tt.height, got)
			}
		})
	}
}

func TestArea(t *testing.T) {
	tests := []struct {
		name string
		rec  rectangle
		want float64
	}{
		{"square", rectangle{width: 10, height: 10}, 100},
		{"unit square", rectangle{width: 1, height: 1}, 1},
		{"wider than tall", rectangle{width: 3, height: 4}, 12},
		{"taller than wide", rectangle{width: 4, height: 3}, 12},
		{"zero width", rectangle{width: 0, height: 5}, 0},
		{"zero height", rectangle{width: 5, height: 0}, 0},
		{"fractional side", rectangle{width: 2.5, height: 4}, 10},
		{"both sides fractional", rectangle{width: 0.1, height: 0.2}, 0.02},
		{"very small sides", rectangle{width: 1e-8, height: 1e-8}, 1e-16},
		{"very large sides", rectangle{width: 1e150, height: 1e150}, 1e300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rec.area()
			if !closeEnough(got, tt.want) {
				t.Errorf("area() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArea_overflowsToInf(t *testing.T) {
	r := rectangle{width: math.MaxFloat64, height: 2}

	got := r.area()

	if !math.IsInf(got, 1) {
		t.Errorf("area() = %v, want +Inf", got)
	}
}

func TestArea_doesNotMutateReceiver(t *testing.T) {
	r := rectangle{width: 3, height: 4}

	r.area()

	if r.width != 3 || r.height != 4 {
		t.Errorf("area() mutated receiver: got %+v, want {width:3 height:4}", r)
	}
}

func TestArea_fromNewRectangle(t *testing.T) {
	r, err := newRectangle(2.5, 4)
	if err != nil {
		t.Fatalf("newRectangle(2.5, 4) error = %v", err)
	}

	got := r.area()
	if !closeEnough(got, 10) {
		t.Errorf("area() = %v, want 10", got)
	}
}

func TestPerimeter(t *testing.T) {
	tests := []struct {
		name string
		rec  rectangle
		want float64
	}{
		{"square", rectangle{width: 10, height: 10}, 40},
		{"unit square", rectangle{width: 1, height: 1}, 4},
		{"wider than tall", rectangle{width: 3, height: 4}, 14},
		{"taller than wide", rectangle{width: 4, height: 3}, 14},
		{"zero width", rectangle{width: 0, height: 5}, 10},
		{"zero height", rectangle{width: 5, height: 0}, 10},
		{"fractional side", rectangle{width: 2.5, height: 0.5}, 6},
		{"both sides fractional", rectangle{width: 0.1, height: 0.2}, 0.6},
		{"very small sides", rectangle{width: 1e-8, height: 1e-8}, 4e-8},
		{"very large sides", rectangle{width: 1e150, height: 1e150}, 4e150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rec.perimeter()
			if !closeEnough(got, tt.want) {
				t.Errorf("perimeter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerimeter_overflowsToInf(t *testing.T) {
	r := rectangle{width: math.MaxFloat64, height: math.MaxFloat64}

	got := r.perimeter()

	if !math.IsInf(got, 1) {
		t.Errorf("perimeter() = %v, want +Inf", got)
	}
}
