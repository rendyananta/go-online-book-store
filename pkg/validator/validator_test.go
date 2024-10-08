package validator

import (
	"testing"
)

type TestValidation struct {
	Name string `validate:"required"`
}

func TestStruct(t *testing.T) {
	type args struct {
		any any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "struct validation",
			args: args{
				any: TestValidation{Name: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Struct(tt.args.any); (err != nil) != tt.wantErr {
				t.Errorf("Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
