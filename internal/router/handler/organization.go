package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/router/handler/response"
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
	Id string `json:"id"`
	Name string `json:"name" binding:"required"`
	ParentId string `json:"parentId" binding:"required"`
	Type int32  `json:"type"`   // 0 单位、 1 工资表
	SalaryType int32  `json:"salaryType"` // 0:工资 1:福利 2: 退休    工资类型
	EmployeeType int32 `json:"employeeType"` // 员工类型： 0: 公务员  1:事业 2: 企业
	Children []*Organization `json:"children"`
}

func (org *Organization) toService() *service.Organization{
	return &service.Organization{
		Id: org.Id,
		Name: org.Name,
		ParentId: org.ParentId,
		Type: org.Type,
		SalaryType: org.SalaryType,
		EmployeeType: org.EmployeeType,
	}
}

func (org *Organization) fromService(so  *service.Organization){
	org.Id= so.Id
	org.Name= so.Name
	org.ParentId= so.ParentId
	org.Type= so.Type
	org.SalaryType= so.SalaryType
	org.EmployeeType= so.EmployeeType
}
// @Summary 获取组织列表
// @Description 获取组织列表
// @Tags 组织管理
// @Accept application/json
// @Success 200 {object} Organization
// @Router /v1/auth/organization [get]
func (r *OrganizationHandler) ListOrganization(c *gin.Context) {
	ctx := c.Request.Context()
	orgs,err := r.org.ListOrganization(ctx)
	if err != nil {
		log.Println("ListOrganization error: ", err)
		c.JSON(http.StatusBadRequest,err)
		return
	}

	c.JSON(http.StatusOK, orgs)
}

// @Summary 添加组织
// @Description 添加组织、工资表
// @Tags 组织管理
// @Accept application/json
// @Param Organization body Organization true ""
// @Success 200 {object} response.SuccessMessage
// @Router /v1/auth/organization [post]
func (r *OrganizationHandler) AddOrganization(c *gin.Context) {
	org := &Organization{}
	err := c.ShouldBindJSON(org)
	if err != nil {
		c.JSON(http.StatusBadRequest,err)
		return
	}

	ctx := c.Request.Context()
	err = r.org.AddOrganization(ctx, &service.Organization{
		ParentId: org.ParentId,
		Name:org.Name,
		SalaryType: org.SalaryType,   // 0:工资 1:福利 2: 退休    工资类型
		EmployeeType: org.EmployeeType, // 员工类型： 0: 公务员  1:事业 2: 企业
		Type: org.Type,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest,err)
		return
	}

	response.Success(c,nil)
}


// @Summary 添加组织
// @Description 添加组织、工资表
// @Tags 组织管理
// @Accept application/json
// @Param Organization body Organization true ""
// @Success 200 {object} response.SuccessMessage
// @Router /v1/auth/organization [update]
func (r *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	org := &Organization{}
	err := c.ShouldBindJSON(org)
	if err != nil {
		c.JSON(http.StatusBadRequest,err)
		return
	}

	if org.Id == "" {
		c.JSON(http.StatusBadRequest, "数据不存在")
		return
	}

	ctx := c.Request.Context()
	err = r.org.UpdateOrganization(ctx, &service.Organization{
		Id: org.Id,
		ParentId: org.ParentId,
		Name:org.Name,
		SalaryType: org.SalaryType,   // 0:工资 1:福利 2: 退休    工资类型
		EmployeeType: org.EmployeeType, // 员工类型： 0: 公务员  1:事业 2: 企业
		Type: org.Type,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest,err)
		return
	}

	response.Success(c,nil)
}


// @Summary 添加组织
// @Description 添加组织、工资表
// @Tags 组织管理
// @Accept application/json
// @Param Organization body Organization true ""
// @Success 200 {object} response.SuccessMessage
// @Router /v1/auth/organization/{id} [get]
func (r *OrganizationHandler) GetOrganization(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest,"参数不完整")
		return
	}

	ctx := c.Request.Context()
	sorg, err := r.org.GetOrganization(ctx, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, sorg)
}
