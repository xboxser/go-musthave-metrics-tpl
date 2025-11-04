package middleware

import "testing"

func TestIsCompressible(t *testing.T) {
	type args struct {
		contentType string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "valid parameters 1", args: args{contentType: "application/json"}, want: true},
		{name: "valid parameters 2", args: args{contentType: "text/html"}, want: true},
		{name: "no valid parameters", args: args{contentType: "zip"}, want: false},
		{name: "empty parameters", args: args{contentType: ""}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCompressible(tt.args.contentType)
			if result != tt.want {
				t.Errorf("isCompressible(%q) = %v; expected %v", tt.args.contentType, result, tt.want)
			}
		})
	}
}
