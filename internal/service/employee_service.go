package service

import (
	"context"
	"takeout/common/e"
	"takeout/common/enum"
	"takeout/common/utils"
	"takeout/global"
	"takeout/internal/api/request"
	"takeout/internal/api/response"
	"takeout/repository"
)

type IEmployeeService interface {
	Login(context.Context, request.EmployeeLogin) (*response.EmployeeLogin, error)
	//Logout(ctx context.Context) error
	//EditPassword(context.Context, request.EmployeeEditPassword) error
	//CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error
	//PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error)
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
	if password != employeeLogin.Password {
		return nil, e.Error_PASSWORD_ERROR
	}
	if employee.Status == enum.DISABLE {
		return nil, e.Error_ACCOUNT_LOCKED
	}

	// 返回jwt
	jwtConfig := global.Config.Jwt.Admin
	token, err := utils

	return nil, err
}

type EmployeeImpl struct {
	repo repository.EmployeeRepo
}

func NewEmployeeService(repo repository.EmployeeRepo) IEmployeeService {
	return &EmployeeImpl{repo: repo}
}
