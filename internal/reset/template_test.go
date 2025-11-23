package reset

import (
	"go/ast"
	"testing"
)

func TestZeroValue(t *testing.T) {
	type args struct {
		expr       ast.Expr
		structName string
		fieldName  string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "int field",
			args: args{
				expr:       &ast.Ident{Name: "int"},
				structName: "TestStruct",
				fieldName:  "IntegerField",
			},
			want: "TestStruct.IntegerField = 0",
		},
		{
			name: "float64 field",
			args: args{
				expr:       &ast.Ident{Name: "float64"},
				structName: "TestStruct",
				fieldName:  "Float64Field",
			},
			want: "TestStruct.Float64Field = 0",
		},
		{
			name: "string field",
			args: args{
				expr:       &ast.Ident{Name: "string"},
				structName: "TestStruct",
				fieldName:  "StringField",
			},
			want: `TestStruct.StringField = ""`,
		},
		{
			name: "bool field",
			args: args{
				expr:       &ast.Ident{Name: "bool"},
				structName: "TestStruct",
				fieldName:  "BoolField",
			},
			want: "TestStruct.BoolField = false",
		},
		{
			name: "slice field",
			args: args{
				expr:       &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
				structName: "TestStruct",
				fieldName:  "SliceField",
			},
			want: "TestStruct.SliceField = TestStruct.SliceField[:0]",
		},
		{
			name: "map field",
			args: args{
				expr:       &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}},
				structName: "TestStruct",
				fieldName:  "MapField",
			},
			want: "clear(TestStruct.MapField)",
		},
		{
			name: "pointer to primitive",
			args: args{
				expr:       &ast.StarExpr{X: &ast.Ident{Name: "string"}},
				structName: "TestStruct",
				fieldName:  "StringPtrField",
			},
			want: "if TestStruct.StringPtrField != nil {\n    *TestStruct.StringPtrField = \"\"\n}",
		},
		{
			name: "struct field",
			args: args{
				expr:       &ast.Ident{Name: "OtherStruct"},
				structName: "TestStruct",
				fieldName:  "StructField",
			},
			want: "TestStruct.StructField = OtherStruct{}",
		},
		{
			name: "inline struct",
			args: args{
				expr:       &ast.StructType{},
				structName: "TestStruct",
				fieldName:  "InlineStructField",
			},
			want: "TestStruct.InlineStructField = {}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zeroValue(tt.args.expr, tt.args.structName, tt.args.fieldName); got != tt.want {
				t.Errorf("zeroValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
