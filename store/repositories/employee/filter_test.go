package employee

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestToQuery checks that the ToQuery function generates the correct SQL query and arguments.
func TestToQuery(t *testing.T) {
	trueValue := true
	falseValue := false

	tests := []struct {
		name      string
		filter    SummaryFilter
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			name: "OnContract is true",
			filter: SummaryFilter{
				OnContract: &trueValue,
			},
			wantQuery: " WHERE on_contract = $1",
			wantArgs:  []interface{}{trueValue},
		},
		{
			name: "OnContract is false",
			filter: SummaryFilter{
				OnContract: &falseValue,
			},
			wantQuery: " WHERE on_contract = $1",
			wantArgs:  []interface{}{falseValue},
		},
		{
			name: "OnContract is nil",
			filter: SummaryFilter{
				OnContract: nil,
			},
			wantQuery: "",
			wantArgs:  []interface{}(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args := tt.filter.ToQuery()
			assert.Equal(t, tt.wantQuery, query, "Expected and actual query strings do not match")
			assert.Equal(t, tt.wantArgs, args, "Expected and actual arguments do not match")
		})
	}
}
