package user

import (
	"strconv"
	"strings"
)

type UserFilter struct {
	Username string
	Password string
}

func (f UserFilter) ToQuery() (string, []interface{}) {
	var query strings.Builder
	var args []interface{}
	var conditions []string

	i := 1
	if len(f.Username) > 0 {
		conditions = append(conditions, "username = $"+strconv.Itoa(i))
		args = append(args, f.Username)
		i++
	}

	if len(f.Password) > 0 {
		conditions = append(conditions, "password = $"+strconv.Itoa(i))
		args = append(args, f.Password)
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	return query.String(), args
}
