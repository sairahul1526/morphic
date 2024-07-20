package employee

import (
	"strconv"
	"strings"
)

type SummaryFilter struct {
	OnContract *bool
}

func (f *SummaryFilter) ToQuery() (string, []interface{}) {
	var query strings.Builder
	var args []interface{}
	var conditions []string

	i := 1
	if f.OnContract != nil {
		conditions = append(conditions, "on_contract = $"+strconv.Itoa(i))
		args = append(args, *f.OnContract)
		i++
	}

	if len(conditions) > 0 {
		query.WriteString(" WHERE " + strings.Join(conditions, " AND "))
	}

	return query.String(), args
}
