package handler

import "testing"

func TestValidateCountParams(t *testing.T) {
	tests := []struct {
		name  string
		value []string
		want  bool
	}{
		{name: "valid parameters", value: []string{"1", "2", "3", "4"}, want: true},
		{name: "empty parameters", value: []string{}, want: false},
		{name: "many parameters", value: []string{"1", "2", "3", "4", "5"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := validateCountParams(tt.value); res != tt.want {
				t.Errorf("validateCountParams() = %v, want %v", res, tt.want)
			}
		})
	}
}

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

func TestValidate404(t *testing.T) {
	tests := []struct {
		name  string
		value []string
		want  bool
	}{
		{name: "valid parameters", value: []string{"1", "2", "3", "4"}, want: true},
		{name: "short parameters", value: []string{}, want: false},
		{name: "many parameters", value: []string{"1", "2", "3", "4", "5"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if res := validate404(tt.value); res != tt.want {
				t.Errorf("validate404() = %v, want %v", res, tt.want)
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
			if res := valiteValueMetrics(tt.value); res != tt.want {
				t.Errorf("valiteValueMetrics() = %v, want %v", res, tt.want)
			}
		})
	}
}
