package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/service"
	"log"
	"net/http"
)

type OrganizationHandler struct {
	org *service.OrganizationService
}

func NewOrganizationHandler(org *service.OrganizationService) *OrganizationHandler{
	return &OrganizationHandler{
		org:org,
	}
}

type Organization struct {
	Id int `json:"id"`
	Name string `json:"name"`
	OrganizationSalary OrganizationSalary `json:"organizationSalary"`
	Type int32  `json:"type"`   // 0 单位、 1 工资表
	ParentId int `json:"parentId"`
	Children []*Organization `json:"children"`
}

type OrganizationSalary  struct {
	Id string `json:"id"`  // OrganizationId
	SalaryType int32  `json:"salaryType"` // 0:工资 1:福利 2: 退休    工资类型
	EmployeeType int32 `json:"employeeType"` // 员工类型： 0: 公务员  1:事业 2: 企业
}


// GET
func (r *OrganizationHandler) Organization(c *gin.Context) {
	ctx := c.Request.Context()
	orgs,err := r.org.ListOrganization(ctx)
	if err != nil {
		log.Println("ListOrganization error: ", err)
		c.JSON(http.StatusBadRequest,err)
		return
	}

	c.JSON(http.StatusOK,orgs)
}