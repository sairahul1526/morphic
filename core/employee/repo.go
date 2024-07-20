package employee

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sairahul1526/morphic/entities"
)

type Repository interface {
	Create(context.Context, *sqlx.Tx, entities.Employee) (entities.Employee, error)
	Delete(context.Context, []string) error
	GetSummary(context.Context, entities.SummaryFilter) (entities.Summary, error)
	GetSummaryByDepartment(context.Context, entities.SummaryFilter) (map[string]entities.Summary, error)
	GetSummaryByDepartmentAndSubDepartment(context.Context, entities.SummaryFilter) (map[string]map[string]entities.Summary, error)
}
