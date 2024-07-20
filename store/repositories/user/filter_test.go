package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestToQuery checks that the ToQuery function generates the correct SQL query and arguments.
func TestToQuery(t *testing.T) {
	tests := []struct {
		name      string
		filter    UserFilter
		wantQuery string
		wantArgs  []interface{}
	}{
		{
			name: "Username and Password",
			filter: UserFilter{
				Username: "testuser",
				Password: "password123",
			},
			wantQuery: " WHERE username = $1 AND password = $2",
			wantArgs:  []interface{}{"testuser", "password123"},
		},
		{
			name: "Only Username",
			filter: UserFilter{
				Username: "testuser",
				Password: "",
			},
			wantQuery: " WHERE username = $1",
			wantArgs:  []interface{}{"testuser"},
		},
		{
			name: "Only Password",
			filter: UserFilter{
				Username: "",
				Password: "password123",
			},
			wantQuery: " WHERE password = $1",
			wantArgs:  []interface{}{"password123"},
		},
		{
			name: "Empty fields",
			filter: UserFilter{
				Username: "",
				Password: "",
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
