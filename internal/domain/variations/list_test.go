package variations

import (
	"testing"
)

func TestGetVariation(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"linear", "linear", false},
		{"swirl", "swirl", false},
		{"sinusoidal", "sinusoidal", false},
		{"unknown", "unknown", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := GetVariation(tt.input)

			if (err != nil) != tt.expectError {
				t.Errorf("GetVariation(%s) error = %v, expectError = %t", tt.input, err, tt.expectError)
			}
			if !tt.expectError && fn == nil {
				t.Errorf("GetVariation(%s) returned nil function", tt.input)
			}
		})
	}
}
