package handler

import "testing"

func TestGetParamsURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{name: "valid read get params", value: "one/two/fhree/", want: []string{"one", "two", "fhree"}},
		{name: "big read get params", value: "o/ne/two/fh/ree/", want: []string{"o", "ne", "two", "fh", "ree"}},
		{name: "plenty ///", value: "one/////two//fhree////", want: []string{"one", "two", "fhree"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := getParamsURL(tt.value)
			if len(res) != len(tt.want) {
				t.Errorf("getParamsURL() = %v, want %v", res, tt.want)
			}

			for i := 0; i < len(res); i++ {
				if res[i] != tt.want[i] {
					t.Errorf("getParamsURL = %v, want %v", res[i], tt.want[i])
					return
				}
			}
		})
	}
}

func TestValiteValueMetrics(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "int", value: "123", want: true},
		{name: "string", value: "ffffff", want: false},
		{name: "float", value: "55.0", want: true},
		{name: "empty", value: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := validateValueMetrics(tt.value); res != tt.want {
				t.Errorf("validateValueMetrics() = %v, want %v", res, tt.want)
			}
		})
	}
}

func TestValidateTypeMetrics(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "gauge", value: "gauge", want: true},
		{name: "counter", value: "counter", want: true},
		{name: "random string", value: "trololo", want: false},
		{name: "empty", value: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := validateTypeMetrics(tt.value); res != tt.want {
				t.Errorf("validateTypeMetrics() = %v, want %v", res, tt.want)
			}
		})
	}
}
