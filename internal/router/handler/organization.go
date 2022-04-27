package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/pkg/requestid"
	"github.com/xiaowuzai/payroll/internal/pkg/response"
	"github.com/xiaowuzai/payroll/internal/service"
	"net/http"
)


type OrganizationHandler struct {
	org    *service.OrganizationService
	logger *logger.Logger
}

func NewOrganizationHandler(org *service.OrganizationService, logger *logger.Logger) *OrganizationHandler {
	return &OrganizationHandler{
		org:    org,
		logger: logger,
	}
}

type Organization struct {
	Id           string          `json:"id"`
	Name         string          `json:"name" binding:"required"`
	ParentId     string          `json:"parentId" binding:"required"`
	SalaryType   string          `json:"salaryType"`   //   工资类型
	Type         int32           `json:"type"`         // 0 单位、 1 工资表
	FeeType      int32           `json:"feeType"`      // 0:工资 1:福利 2: 退休    费用类型
	EmployeeType int32           `json:"employeeType"` // 员工类型： 0: 公务员  1:事业 2: 企业
	Children     []*Organization `json:"children"`
}

func (org *Organization) toService() *service.Organization {
	return &service.Organization{
		Id:           org.Id,
		Name:         org.Name,
		ParentId:     org.ParentId,
		Type:         org.Type,
		SalaryType:   org.SalaryType,
		EmployeeType: org.EmployeeType,
	}
}

func (org *Organization) fromService(so *service.Organization) {
	org.Id = so.Id
	org.Name = so.Name
	org.ParentId = so.ParentId
	org.Type = so.Type
	org.SalaryType = so.SalaryType
	org.EmployeeType = so.EmployeeType
}

// @Summary 获取组织列表
// @Description 获取组织列表
// @Tags 组织管理
// @Accept application/json
// @Success 200 {object} Organization
// @Router /v1/auth/organization [get]
func (oh *OrganizationHandler) ListOrganization(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := oh.logger.WithRequestId(ctx)
	log.Info("ListOrganization function called")

	orgs, err := oh.org.ListOrganization(ctx)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("ListOrganization function success")
	c.JSON(http.StatusOK, orgs)
}

// @Summary 添加组织
// @Description 添加组织、工资表
// @Tags 组织管理
// @Accept application/json
// @Param Organization body Organization true ""
// @Success 200 {object} response.SuccessMessage
// @Router /v1/auth/organization [post]
func (oh *OrganizationHandler) AddOrganization(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := oh.logger.WithRequestId(ctx)
	log.Info("AddOrganization function called")

	org := &Organization{}
	err := c.ShouldBindJSON(org)
	if err != nil {
		log.Error("AddOrganization  ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	err = oh.org.AddOrganization(ctx, &service.Organization{
		ParentId:     org.ParentId,
		Name:         org.Name,
		SalaryType:   org.SalaryType,   // 0:工资 1:福利 2: 退休    工资类型
		EmployeeType: org.EmployeeType, // 员工类型： 0: 公务员  1:事业 2: 企业
		Type:         org.Type,
	})
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("AddOrganization function success")
	response.Success(c, nil)
}

// @Summary 添加组织
// @Description 添加组织、工资表
// @Tags 组织管理
// @Accept application/json
// @Param Organization body Organization true ""
// @Success 200 {object} response.SuccessMessage
// @Router /v1/auth/organization [update]
func (oh *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := oh.logger.WithRequestId(ctx)
	log.Info("UpdateOrganization function called")

	org := &Organization{}
	err := c.ShouldBindJSON(org)
	if err != nil {
		log.Error("UpdateOrganization ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	if org.Id == "" {
		response.ParamsError(c, "id不存在")
		return
	}


	err = oh.org.UpdateOrganization(ctx, &service.Organization{
		Id:           org.Id,
		ParentId:     org.ParentId,
		Name:         org.Name,
		SalaryType:   org.SalaryType,   // 0:工资 1:福利 2: 退休    工资类型
		EmployeeType: org.EmployeeType, // 员工类型： 0: 公务员  1:事业 2: 企业
		Type:         org.Type,
	})
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("UpdateOrganization function called")
	response.Success(c, nil)
}

func (oh *OrganizationHandler) GetOrganization(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := oh.logger.WithRequestId(ctx)
	log.Info("GetOrganization function called")

	id := c.Param("id")
	if id == "" {
		log.Error("GetOrganization  id is empty")
		response.ParamsError(c, ErrIdEmpty)
		return
	}

	sorg, err := oh.org.GetOrganization(ctx, id)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("GetOrganization function success")
	response.Success(c, sorg)
}

func (oh *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := oh.logger.WithRequestId(ctx)
	log.Info("GetOrganization function called")

	req := &RequestId{}
	err :=  c.ShouldBindJSON(req)
	if err != nil {
		log.Error("DeleteOrganization ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	err = oh.org.DeleteOrganization(ctx, req.Id)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("DeleteOrganization function success")
	response.Success(c, nil)
}
