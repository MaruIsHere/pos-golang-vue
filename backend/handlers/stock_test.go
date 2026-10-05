package handlers

import "testing"

func TestIsValidStockQuantity(t *testing.T) {
	tests := []struct {
		name     string
		quantity float64
		unit     string
		valid    bool
	}{
		{name: "whole pieces", quantity: 3, unit: "pcs", valid: true},
		{name: "fractional pieces", quantity: 1.5, unit: "pcs", valid: false},
		{name: "fractional grams", quantity: 12.5, unit: "gram", valid: true},
		{name: "fractional liters", quantity: 0.25, unit: "liter", valid: true},
		{name: "three decimal places", quantity: 0.001, unit: "liter", valid: true},
		{name: "more than three decimal places", quantity: 0.0001, unit: "gram", valid: false},
		{name: "negative quantity", quantity: -1, unit: "gram", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidStockQuantity(tt.quantity, tt.unit); got != tt.valid {
				t.Errorf("isValidStockQuantity(%v, %q) = %v, want %v", tt.quantity, tt.unit, got, tt.valid)
			}
		})
	}
}
