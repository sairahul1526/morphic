package employee

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/raystack/salt/db"
	"github.com/sairahul1526/morphic/entities"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/store"
	"go.uber.org/zap"
)

type EmployeeRepository struct {
	db *db.Client
}

func NewRepository(db *db.Client) *EmployeeRepository {
	return &EmployeeRepository{db}
}

func (r *EmployeeRepository) Create(ctx context.Context, tx *sqlx.Tx, employee entities.Employee) (createdEmployee entities.Employee, err error) {
	startTime := time.Now()
	defer store.RecordMetrics("employee", "create", startTime, err)

	var createdEmployeeModel EmployeeModel
	query := `INSERT INTO employees (id, name, currency, salary, department, sub_department, on_contract, status, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), $9, now(), $10) RETURNING *`

	args := []interface{}{employee.ID, employee.Name, employee.Currency, employee.Salary, employee.Department, employee.SubDepartment, employee.OnContract, employee.Status, employee.CreatedBy, employee.UpdatedBy}
	if tx == nil {
		err = r.db.QueryRowxContext(ctx, query, args...).StructScan(&createdEmployeeModel)
	} else {
		err = tx.QueryRowxContext(ctx, query, args...).StructScan(&createdEmployeeModel)
	}

	if err != nil {
		logger.Error("failed to create employee", []zap.Field{
			zap.Any("employee", employee),
			zap.Error(err),
		}...)
		return
	}

	return createdEmployeeModel.ToDomain(), nil
}

func (r *EmployeeRepository) Delete(ctx context.Context, ids []string) (err error) {
	startTime := time.Now()
	defer store.RecordMetrics("employee", "delete", startTime, err)

	query := `DELETE FROM employees WHERE id in (`
	args := []interface{}{}

	// increment i and add to query and also add to args
	for i, id := range ids {
		query += "$" + strconv.Itoa(i+1) + ","
		args = append(args, id)
	}

	// remove last comma
	query = query[:len(query)-1]
	query += ")"

	_, err = r.db.ExecContext(ctx, query, args...)
	if err == sql.ErrNoRows {
		logger.Error("failed to delete employees", []zap.Field{
			zap.Strings("employee_ids", ids),
			zap.Error(err),
		}...)
		return nil
	}
	return err
}

func (r *EmployeeRepository) GetSummary(ctx context.Context, filters entities.SummaryFilter) (entities.Summary, error) {
	startTime := time.Now()
	defer store.RecordMetrics("employee", "get_summary", startTime, nil)

	summaryFilter := SummaryFilter{
		OnContract: filters.OnContract,
	}
	filterQuery, args := summaryFilter.ToQuery()

	query := `SELECT AVG(salary) as mean, MIN(salary) as min, MAX(salary) as max FROM employees` + filterQuery

	var summary Summary
	err := r.db.QueryRowxContext(ctx, query, args...).StructScan(&summary)
	if err != nil {
		logger.Error("failed to fetch salary summary", []zap.Field{
			zap.Error(err),
		}...)
		return entities.Summary{}, err
	}

	return summary.ToDomain(), nil
}

func (r *EmployeeRepository) GetSummaryByDepartment(ctx context.Context, filters entities.SummaryFilter) (map[string]entities.Summary, error) {
	startTime := time.Now()
	defer store.RecordMetrics("employee", "get_summary_by_department", startTime, nil)

	summaryFilter := SummaryFilter{
		OnContract: filters.OnContract,
	}
	filterQuery, args := summaryFilter.ToQuery()

	query := `SELECT department, AVG(salary) as mean, MIN(salary) as min, MAX(salary) as max FROM employees` + filterQuery + ` GROUP BY department`

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		logger.Error("failed to fetch salary summary by department", []zap.Field{
			zap.Error(err),
		}...)
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]entities.Summary)
	for rows.Next() {
		var summary Summary
		var department string
		err = rows.Scan(&department, &summary.Mean, &summary.Min, &summary.Max)
		if err != nil {
			logger.Error("failed to scan salary summary by department", []zap.Field{
				zap.Error(err),
			}...)
			return nil, err
		}
		result[department] = summary.ToDomain()
	}

	return result, nil
}

func (r *EmployeeRepository) GetSummaryByDepartmentAndSubDepartment(ctx context.Context, filters entities.SummaryFilter) (map[string]map[string]entities.Summary, error) {
	startTime := time.Now()
	defer store.RecordMetrics("employee", "get_summary_by_department_and_sub_department", startTime, nil)

	summaryFilter := SummaryFilter{
		OnContract: filters.OnContract,
	}
	filterQuery, args := summaryFilter.ToQuery()

	query := `SELECT department, sub_department, AVG(salary) as mean, MIN(salary) as min, MAX(salary) as max FROM employees` + filterQuery + ` GROUP BY department, sub_department`

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		logger.Error("failed to fetch salary summary by department and sub-department", []zap.Field{
			zap.Error(err),
		}...)
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]map[string]entities.Summary)
	for rows.Next() {
		var summary Summary
		var department, subDepartment string
		err = rows.Scan(&department, &subDepartment, &summary.Mean, &summary.Min, &summary.Max)
		if err != nil {
			logger.Error("failed to scan salary summary by department and sub-department", []zap.Field{
				zap.Error(err),
			}...)
			return nil, err
		}
		if _, ok := result[department]; !ok {
			result[department] = make(map[string]entities.Summary)
		}
		result[department][subDepartment] = summary.ToDomain()
	}

	return result, nil
}
