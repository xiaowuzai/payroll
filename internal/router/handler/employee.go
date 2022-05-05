package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaowuzai/payroll/internal/pkg/logger"
	"github.com/xiaowuzai/payroll/internal/pkg/requestid"
	"github.com/xiaowuzai/payroll/internal/pkg/response"
	"github.com/xiaowuzai/payroll/internal/service"
	"time"
)

type EmployeeHandler struct {
	emp    *service.EmployeeService
	logger *logger.Logger
}

func NewEmployeeHandler(emp *service.EmployeeService, logger *logger.Logger) *EmployeeHandler {
	return &EmployeeHandler{
		emp:    emp,
		logger: logger,
	}
}

type Employee struct {
	Id           string         `json:"id" binding:"omitempty,uuid"`
	Name         string         `json:"name" binding:"required,lte=20"`   // 姓名
	IdCard       string         `json:"idCard" binding:"required,len=18"` // 身份证号码
	Telephone    string         `json:"telephone"`                 // 手机号码
	Duty         string         `json:"duty"`                      // 职务
	Post         string         `json:"post"`                      // 岗位
	Level        string         `json:"level"`                     // 级别
	OfferTime    int64          `json:"offerTime"`
	RetireTime   int64          `json:"retireTime"`
	Number       int            `json:"number"`                    // 员工编号
	Sex          int32          `json:"sex" binding:"oneof=1 2"`                       // 性别： 1: 男、2: 女
	Status       int32          `json:"status"  binding:"oneof=1 2 3"` // 1: 在职, 2:离职、 3: 退休
	BaseSalary   int32          `json:"baseSalary"`                // 基本工资
	Identity     int32          `json:"identity" binding:"oneof=1 2 3"`                  // 身份类别： 1:公务员、 2: 事业、3: 企业
	EmployeeBankInfos []*EmployeeBankInfo `json:"employeeBankInfos" binding:"required"`
}

type EmployeeBankInfo struct {
	Id             string `json:"id" binding:"omitempty,uuid"`
	BankId         string `json:"bankId" binding:"uuid"`
	CardNumber     string `json:"cardNumber" binding:"required"`
	OrganizationId string `json:"organizationId" binding:"uuid"`
}

func (e *Employee) toService() *service.Employee {
	spi := make([]*service.EmployeeBankInfo, 0, len(e.EmployeeBankInfos))
	for _, v := range e.EmployeeBankInfos {
		spi = append(spi, &service.EmployeeBankInfo{
			EmployeeId  : e.Id,
			BankId    : v.BankId,
			CardNumber  : v.CardNumber,
			OrganizationId : v.OrganizationId,
		})
	}

	return &service.Employee{
		Id:         e.Id,
		Name:       e.Name,
		IdCard:     e.IdCard,
		Telephone:  e.Telephone,
		Duty:       e.Duty,
		Post:       e.Post,
		Level:      e.Level,
		OfferTime:  time.Unix(e.OfferTime, 0),
		RetireTime: time.Unix(e.RetireTime, 0),
		Number:     e.Number,
		Sex:        e.Sex,
		Status:     e.Status,
		BaseSalary: e.BaseSalary,
		Identity:   e.Identity,
		EmployeeBankInfos: spi,
	}
}

func (e *Employee) fromService(se *service.Employee) {
	e.Id = se.Id
	e.Name = se.Name
	e.IdCard = se.IdCard
	e.Telephone = se.Telephone
	e.Duty = se.Duty
	e.Post = se.Post
	e.Level = se.Level
	e.OfferTime = se.OfferTime.Unix()
	e.RetireTime = se.RetireTime.Unix()
	e.Number = se.Number
	e.Sex = se.Sex
	e.Status = se.Status
	e.BaseSalary = se.BaseSalary
	e.Identity = se.Identity

	if se.EmployeeBankInfos != nil {
		for _, v := range se.EmployeeBankInfos {
			e.EmployeeBankInfos = append(e.EmployeeBankInfos, &EmployeeBankInfo{
				Id:             v.Id,
				//EmployeeId:     v.EmployeeId,
				BankId:         v.BankId,
				CardNumber:     v.CardNumber,
				OrganizationId: v.OrganizationId,
			})
		}
	}
}

func (eh *EmployeeHandler) AddEmployee(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := eh.logger.WithRequestId(ctx)
	log.Info("AddEmployee function called")

	req := &Employee{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Error("AddEmployee ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	se := req.toService()
	log.Infof("AddEmployee toService %+v\n", *se)
	err = eh.emp.AddEmployee(ctx, se)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("AddEmployee function success")
	response.Success(c, nil)
}

func (eh *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := eh.logger.WithRequestId(ctx)
	log.Info("DeleteEmployee function called")

	req := &RequestId{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Error("DeleteEmployee ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}
	if req.Id == "" {
		log.Error("DeleteEmployee ShouldBindJSON error: ", ErrIdEmpty)
		response.ParamsError(c, ErrIdEmpty)
		return
	}

	log.Infof("DeleteEmployee toService id = %s\n", req.Id)
	err = eh.emp.DeleteEmployee(ctx, req.Id)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("DeleteEmployee function success")
	response.Success(c, nil)
}

func (eh *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := eh.logger.WithRequestId(ctx)
	log.Info("UpdateEmployee function called")

	req := &Employee{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		log.Error("UpdateEmployee ShouldBindJSON error: ", err.Error())
		response.ParamsError(c, err.Error())
		return
	}

	se := req.toService()
	log.Infof("UpdateEmployee toService %+v\n", *se)
	err = eh.emp.UpdateEmployee(ctx, se)
	if err != nil {
		response.WithError(c, err)
		return
	}

	log.Info("UpdateEmployee function success")
	response.Success(c, nil)
}

func (eh *EmployeeHandler) GetEmployee(c *gin.Context) {
	ctx := requestid.WithRequestId(c)
	log := eh.logger.WithRequestId(ctx)
	log.Info("GetEmployee function called")

	id, has := c.Params.Get("id")
	if !has || id == "" {
		log.Error("GetEmployee id is empty")
		response.ParamsError(c, ErrIdEmpty)
		return
	}

	log.Infof("GetEmployee id = %s\n", id)
	se, err := eh.emp.GetEmployee(ctx, id)
	if err != nil {
		response.WithError(c, err)
		return
	}

	e := &Employee{}
	e.fromService(se)

	log.Info("GetEmployee function success")
	response.Success(c, e)
}

func (eh *EmployeeHandler) ListEmployee(c *gin.Context) {
	ctx := requestid.WithRequestId(c)

	log := eh.logger.WithRequestId(ctx)
	log.Info("ListEmployee function called")

	name := c.Query("name")
	organizationId := c.Query("organizationId")

	log.Infof("ListEmployee name = %s, organizationId = %s\n", name, organizationId)
	ses, err := eh.emp.ListEmployee(ctx, name, organizationId)
	if err != nil {
		response.WithError(c, err)
		return
	}

	es := make([]*Employee, 0, len(ses))
	for _, se := range ses {
		e := &Employee{}
		e.fromService(se)
		es = append(es, e)
	}

	log.Info("ListEmployee function success")
	response.Success(c, es)
}
