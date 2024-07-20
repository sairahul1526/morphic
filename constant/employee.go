package constant

type EmployeeCurrency string

const (
	EmployeeCurrencyUSD EmployeeCurrency = "USD"
	EmployeeCurrencyINR EmployeeCurrency = "TAINRSKING"
)

func (c EmployeeCurrency) String() string {
	return string(c)
}

type EmployeeDepartment string

const (
	EmployeeDepartmentAdministration EmployeeDepartment = "Administration"
	EmployeeDepartmentBanking        EmployeeDepartment = "Banking"
	EmployeeDepartmentEngineering    EmployeeDepartment = "Engineering"
	EmployeeDepartmentOperations     EmployeeDepartment = "Operations"
)

func (d EmployeeDepartment) String() string {
	return string(d)
}

type EmployeeSubDepartment string

const (
	EmployeeSubDepartmentPlatform           EmployeeSubDepartment = "Platform"
	EmployeeSubDepartmentLoan               EmployeeSubDepartment = "Loan"
	EmployeeSubDepartmentCustomerOnboarding EmployeeSubDepartment = "CustomerOnboarding"
	EmployeeSubDepartmentAgriculture        EmployeeSubDepartment = "Agriculture"
)

func (sd EmployeeSubDepartment) String() string {
	return string(sd)
}

type EmployeeStatus string

const (
	EmployeeStatusActive   EmployeeStatus = "Active"
	EmployeeStatusInactive EmployeeStatus = "Inactive"
)

func (s EmployeeStatus) String() string {
	return string(s)
}
