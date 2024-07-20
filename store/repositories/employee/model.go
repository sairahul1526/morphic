package employee

import (
	"time"

	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
)

type EmployeeModel struct {
	ID            string    `db:"id"`
	Name          string    `db:"name"`
	Currency      string    `db:"currency"`
	Salary        int       `db:"salary"`
	Department    string    `db:"department"`
	SubDepartment string    `db:"sub_department"`
	OnContract    bool      `db:"on_contract"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	CreatedBy     string    `db:"created_by"`
	UpdatedAt     time.Time `db:"updated_at"`
	UpdatedBy     string    `db:"updated_by"`
}

func (e *EmployeeModel) ToDomain() entities.Employee {
	return entities.Employee{
		ID:            e.ID,
		Name:          e.Name,
		Currency:      constant.EmployeeCurrency(e.Currency),
		Salary:        e.Salary,
		Department:    constant.EmployeeDepartment(e.Department),
		SubDepartment: constant.EmployeeSubDepartment(e.SubDepartment),
		OnContract:    e.OnContract,
		Status:        constant.EmployeeStatus(e.Status),
		CreatedAt:     e.CreatedAt,
		CreatedBy:     e.CreatedBy,
		UpdatedAt:     e.UpdatedAt,
		UpdatedBy:     e.UpdatedBy,
	}
}

type Summary struct {
	Mean float64 `db:"mean"`
	Min  float64 `db:"min"`
	Max  float64 `db:"max"`
}

func (s *Summary) ToDomain() entities.Summary {
	return entities.Summary{
		Mean: s.Mean,
		Min:  s.Min,
		Max:  s.Max,
	}
}
