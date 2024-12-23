package service

import (
	"context"
	"takeout/common"
	"takeout/common/e"
	"takeout/common/enum"
	"takeout/common/utils"
	"takeout/global"
	"takeout/internal/api/request"
	"takeout/internal/api/response"
	"takeout/internal/model"
	"takeout/repository"
)

type IEmployeeService interface {
	Login(context.Context, request.EmployeeLogin) (*response.EmployeeLogin, error)
	Logout(ctx context.Context) error
	//EditPassword(context.Context, request.EmployeeEditPassword) error
	CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error
	PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error)
	//SetStatus(ctx context.Context, id uint64, status int) error
	//UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error
	//GetById(ctx context.Context, id uint64) (model.Employee, error)
}

func (ei *EmployeeImpl) Login(ctx context.Context, employeeLogin request.EmployeeLogin) (*response.EmployeeLogin, error) {
	employee, err := ei.repo.GetByUserName(ctx, employeeLogin.UserName)
	if err != nil || employee == nil {
		return nil, e.Error_ACCOUNT_NOT_FOUND
	}
	password := utils.MD5V(employeeLogin.Password, "", 0)
	if password != employee.Password {
		return nil, e.Error_PASSWORD_ERROR
	}
	if employee.Status == enum.DISABLE {
		return nil, e.Error_ACCOUNT_LOCKED
	}

	// 返回jwt
	jwtConfig := global.Config.Jwt.Admin
	token, err := utils.GenerateToken(employee.Id, jwtConfig.Name, jwtConfig.Secret)
	if err != nil {
		return nil, err
	}

	resp := &response.EmployeeLogin{
		Id:       employee.Id,
		Name:     employee.Name,
		Token:    token,
		UserName: employee.Username,
	}
	return resp, err
}

func (ei *EmployeeImpl) Logout(ctx context.Context) (err error) {
	// 前端删除jwt头数据
	return
}

func (ei *EmployeeImpl) CreateEmployee(ctx context.Context, employee request.EmployeeDTO) (err error) {
	entity := model.Employee{
		Id:       employee.Id,
		Username: employee.UserName,
		Name:     employee.Name,
		Phone:    employee.Phone,
		Sex:      employee.Sex,
		IdNumber: employee.IdNumber,
		Status:   enum.ENABLE,
		Password: utils.MD5V("123456", "", 0),
	}
	err = ei.repo.Insert(ctx, entity)
	return err
}

func (ei *EmployeeImpl) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
	//// 分页查询
	//pageResult, err := ei.repo.PageQuery(ctx, dto)
	//// 屏蔽敏感信息
	//if employees, ok := pageResult.Records.([]model.Employee); ok {
	//	// 替换敏感信息
	//	for key, _ := range employees {
	//		employees[key].Password = "****"
	//		employees[key].IdNumber = "****"
	//		employees[key].Phone = "****"
	//	}
	//	// 重新赋值
	//	pageResult.Records = employees
	//}
	//
	//return pageResult, err
	return nil, nil
}

type EmployeeImpl struct {
	repo repository.EmployeeRepo
}

func NewEmployeeService(repo repository.EmployeeRepo) IEmployeeService {
	return &EmployeeImpl{repo: repo}
}
