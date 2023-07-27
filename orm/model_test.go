package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseModel(t *testing.T) {

	tests := []struct {
		name      string
		entity    any
		wantmodel *model
		wantErr   error
	}{
		{
			name: "test model",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := parseModel(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantmodel, m)

		})
	}
}
