package main

import (
	"fmt"
	"math"
)

type rectangle struct {
	width  float64
	height float64
}

func newRectangle(width, height float64) (rectangle, error) {
	if err := validSide("width", width); err != nil {
		return rectangle{}, err
	}
	if err := validSide("height", height); err != nil {
		return rectangle{}, err
	}
	return rectangle{width: width, height: height}, nil
}

func validSide(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%s must be a finite number, got %v", name, v)
	}
	if v <= 0 {
		return fmt.Errorf("%s must be positive, got %v", name, v)
	}
	return nil
}

func (r rectangle) area() float64 {
	return r.width * r.height
}

func (r rectangle) perimeter() float64 {
	return 2 * (r.width + r.height)
}
