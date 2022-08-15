package luhn

import "testing"

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{"good number", "148825017542369", true},
		{"wrong number", "1234567891011", false},
		{"empty", "", false},
		{"trash", "12qwаруш_+.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.arg); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
