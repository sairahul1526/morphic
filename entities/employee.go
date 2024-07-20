package entities

import (
	"time"

	"github.com/sairahul1526/morphic/constant"
)

type Employee struct {
	ID            string
	Name          string
	Currency      constant.EmployeeCurrency
	Salary        int
	Department    constant.EmployeeDepartment
	SubDepartment constant.EmployeeSubDepartment
	OnContract    bool
	Status        constant.EmployeeStatus
	CreatedAt     time.Time
	CreatedBy     string
	UpdatedAt     time.Time
	UpdatedBy     string
}

type SummaryFilter struct {
	OnContract *bool
}

type Summary struct {
	Mean float64
	Min  float64
	Max  float64
}
