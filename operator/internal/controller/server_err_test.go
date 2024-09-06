package controller

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_serverErr_Error(t *testing.T) {
	type fields struct {
		ErrorMsg string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test Error",
			fields: fields{
				ErrorMsg: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &serverErr{
				ErrorMsg: tt.fields.ErrorMsg,
			}
			assert.Equalf(t, tt.want, e.Error(), "Error()")
		})
	}
}
