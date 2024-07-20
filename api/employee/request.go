package employee

import (
	"github.com/gin-gonic/gin"
	"github.com/sairahul1526/morphic/constant"
	"github.com/sairahul1526/morphic/entities"
)

type Request struct {
	Name          string                         `json:"name" validate:"required" example:"John Doe"`
	Currency      constant.EmployeeCurrency      `json:"currency" validate:"oneof=USD INR" example:"USD"`
	Salary        int                            `json:"salary" validate:"required" example:"100000"`
	Department    constant.EmployeeDepartment    `json:"department" validate:"oneof=Administration Banking Engineering Operations" example:"Engineering"`
	SubDepartment constant.EmployeeSubDepartment `json:"sub_department" validate:"oneof=Platform Loan CustomerOnboarding Agriculture" example:"Platform"`
	OnContract    bool                           `json:"on_contract" validate:"required" example:"true"`
	Status        constant.EmployeeStatus        `json:"status" validate:"oneof=Active Inactive" example:"Active"`
}

func (r Request) toDomain() entities.Employee {
	return entities.Employee{
		Name:          r.Name,
		Currency:      r.Currency,
		Salary:        r.Salary,
		Department:    r.Department,
		SubDepartment: r.SubDepartment,
		OnContract:    r.OnContract,
		Status:        r.Status,
	}
}

func parseFilters(c *gin.Context) entities.SummaryFilter {
	var filters entities.SummaryFilter

	onContract := c.Query("on_contract")
	if onContract == "true" {
		filters.OnContract = &[]bool{true}[0]
	} else if onContract == "false" {
		filters.OnContract = &[]bool{false}[0]
	}

	return filters
}
