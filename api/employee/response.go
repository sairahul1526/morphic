package employee

import (
	"time"

	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
)

type ReadResponse struct {
	ID            string                         `json:"id" example:"47435e2b-d8c4-41ff-9de9-2be3bfc92276"`
	Name          string                         `json:"name" example:"John Doe"`
	Currency      constant.EmployeeCurrency      `json:"currency" example:"USD"`
	Salary        int                            `json:"salary" example:"100000"`
	Department    constant.EmployeeDepartment    `json:"department" example:"Engineering"`
	SubDepartment constant.EmployeeSubDepartment `json:"sub_department" example:"Platform"`
	OnContract    bool                           `json:"on_contract" example:"true"`
	Status        constant.EmployeeStatus        `json:"status" example:"Active"`
	CreatedAt     time.Time                      `json:"created_at" example:"2024-05-27T06:50:20.056549Z"`
	CreatedBy     string                         `json:"created_by" example:"e31ab6f8-d359-4c6a-83c6-bfa32229bb01"`
	UpdatedAt     time.Time                      `json:"updated_at" example:"2024-05-27T06:50:20.056549Z"`
	UpdatedBy     string                         `json:"updated_by" example:"e31ab6f8-d359-4c6a-83c6-bfa32229bb01"`
}

func NewReadResponse(employee entities.Employee) ReadResponse {
	return ReadResponse{
		ID:            employee.ID,
		Name:          employee.Name,
		Currency:      employee.Currency,
		Salary:        employee.Salary,
		Department:    employee.Department,
		SubDepartment: employee.SubDepartment,
		OnContract:    employee.OnContract,
		Status:        employee.Status,
		CreatedAt:     employee.CreatedAt,
		CreatedBy:     employee.CreatedBy,
		UpdatedAt:     employee.UpdatedAt,
		UpdatedBy:     employee.UpdatedBy,
	}
}

type Summary struct {
	Mean float64 `json:"mean"`
	Max  float64 `json:"max"`
	Min  float64 `json:"min"`
}

type DepartmentSummaryResponse struct {
	Summaries []DepartmentSummaries `json:"summaries"`
}

type DepartmentSummaries struct {
	Department string                 `json:"department"`
	Summary    *Summary               `json:"summary,omitempty"`
	Summaries  []SubDepartmentSummary `json:"summaries,omitempty"`
}

type SubDepartmentSummary struct {
	SubDepartment string  `json:"sub_department"`
	Summary       Summary `json:"summary"`
}

func NewSummary(summary entities.Summary) Summary {
	return Summary{
		Mean: summary.Mean,
		Min:  summary.Min,
		Max:  summary.Max,
	}
}

func NewSummaryByDepartment(result map[string]entities.Summary) DepartmentSummaryResponse {
	response := DepartmentSummaryResponse{}
	for department, summary := range result {
		result := NewSummary(summary)
		response.Summaries = append(response.Summaries, DepartmentSummaries{Department: department, Summary: &result})
	}
	return response
}

func NewSummaryByDepartmentAndSubDepartment(result map[string]map[string]entities.Summary) DepartmentSummaryResponse {
	response := DepartmentSummaryResponse{}
	for department, subDepartments := range result {
		departmentSummary := DepartmentSummaries{Department: department}
		for subDepartment, summary := range subDepartments {
			departmentSummary.Summaries = append(departmentSummary.Summaries, SubDepartmentSummary{SubDepartment: subDepartment, Summary: NewSummary(summary)})
		}
		response.Summaries = append(response.Summaries, departmentSummary)
	}
	return response
}
