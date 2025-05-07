package faker

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSimpleFaker_Generate(t *testing.T) {
	type args struct {
		typeName string
		params   map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		expect  func(output string) bool
	}{
		{
			name: "full_name",
			args: args{
				typeName: "full_name",
			},
			expect: func(output string) bool {
				return output[0:1] == `"` && strings.Contains(output, " ")
			},
		},
		{
			name: "surname",
			args: args{
				typeName: "surname",
			},
			expect: func(output string) bool {
				return len(output) >= 5
			},
		},
		{
			name: "first_name",
			args: args{
				typeName: "first_name",
			},
			expect: func(output string) bool {
				return len(output) >= 5
			},
		},
		{
			name: "title",
			args: args{
				typeName: "title",
			},
			expect: func(output string) bool {
				return len(output) <= 6
			},
		},
		{
			name: "phone",
			args: args{
				typeName: "phone",
			},
			expect: func(output string) bool {
				return len(output) == 12
			},
		},
		{
			name: "address",
			args: args{
				typeName: "address",
			},
			expect: func(output string) bool {
				return strings.Contains(output, " ")
			},
		},
		{
			name: "zip_code",
			args: args{
				typeName: "zip_code",
			},
			expect: func(output string) bool {
				return len(output) == 7
			},
		},
		{
			name: "url",
			args: args{
				typeName: "url",
			},
			expect: func(output string) bool {
				return strings.Contains(output, "://")
			},
		},
		{
			name: "mac_address",
			args: args{
				typeName: "mac_address",
			},
			expect: func(output string) bool {
				return strings.Contains(output, ":")
			},
		},
		{
			name: "price",
			args: args{
				typeName: "price",
			},
			expect: func(output string) bool {
				return strings.Contains(output, ".")
			},
		},
		{
			name: "negative_price",
			args: args{
				typeName: "negative_price",
			},
			expect: func(output string) bool {
				return strings.Contains(output, "-")
			},
		},
		{
			name: "credit_card",
			args: args{
				typeName: "credit_card",
			},
			expect: func(output string) bool {
				return len(output) >= 14
			},
		},
		{
			name: "bool",
			args: args{
				typeName: "bool",
			},
			expect: func(output string) bool {
				return output == `"True"` || output == `"False"`
			},
		},
		{
			name: "yes_no",
			args: args{
				typeName: "yes_no",
			},
			expect: func(output string) bool {
				return output == `"Yes"` || output == `"No"`
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assemble
			gen, err := GetGenerator(tt.args.typeName, tt.args.typeName, tt.args.params)
			assert.NoError(t, err)

			// Act
			val, err := gen.Generate()

			// Assert
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.expect != nil {
				assert.True(t, tt.expect(val.(string)), "expectation failed with %v (len=%d)", val, len(val.(string)))
			}
		})
	}
}
