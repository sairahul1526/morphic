package server

import (
	"github.com/gin-gonic/gin"
	"github.com/raystack/salt/db"
	employeeapi "github.com/sairahul1526/morphic/api/employee"
	userapi "github.com/sairahul1526/morphic/api/user"
	"github.com/sairahul1526/morphic/config"
	employeecore "github.com/sairahul1526/morphic/core/employee"
	usercore "github.com/sairahul1526/morphic/core/user"
	employeerepo "github.com/sairahul1526/morphic/store/repositories/employee"
	userrepo "github.com/sairahul1526/morphic/store/repositories/user"
)

func registerEmployeeV1APIs(db *db.Client, router *gin.RouterGroup) {
	// employee
	employeeRepo := employeerepo.NewRepository(db)
	employeeService := employeecore.NewService(employeeRepo)
	employeeHandlers := employeeapi.NewEmployeeHandlers(employeeService)
	router.POST("/employees", DBTransactionMiddleware(db), employeeHandlers.CreateEmployee)
	router.DELETE("/employees", employeeHandlers.DeleteEmployee)
	router.GET("/employees/summary", employeeHandlers.GetSummary)
	router.GET("/employees/summary/department", employeeHandlers.GetSummaryByDepartment)
	router.GET("/employees/summary/subdepartment", employeeHandlers.GetSummaryByDepartmentAndSubDepartment)
}

func registerUserV1APIs(db *db.Client, router *gin.RouterGroup, cfg config.Config) {
	// user
	userRepo := userrepo.NewRepository(db)
	userService := usercore.NewService(userRepo)
	userHandlers := userapi.NewUserHandlers(userService, cfg)
	router.POST("/users/login", userHandlers.Login)
}
