package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArray(t *testing.T) {

	tests := []struct {
		name   string
		entity any

		wantVals []any
		wantErr  error
	}{
		{
			name:   "[]int",
			entity: [3]int{1, 2, 3},

			wantVals: []any{1, 2, 3},
		},
		{
			name:   "slice",
			entity: []int{1, 2, 3},

			wantVals: []any{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vals, err := IterateArrayOrSlice(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantVals, vals)
		})
	}
}

func TestIterateMap(t *testing.T) {

	tests := []struct {
		name   string
		entity any

		wantKeys []any
		wantVals []any
		wantErr  error
	}{
		{
			name: "map",
			entity: map[string]string{
				"A": "a",
				"B": "b",
				"C": "c",
			},
			wantKeys: []any{"A", "B", "C"},
			wantVals: []any{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, values, err := IterateMap(tt.entity)

			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.EqualValues(t, tt.wantKeys, keys)
			assert.EqualValues(t, tt.wantVals, values)
		})
	}
}
