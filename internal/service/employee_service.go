package service

import (
	"context"
	"takeout/common"
	"takeout/common/e"
	"takeout/common/enum"
	"takeout/common/utils"
	"takeout/global"
	"takeout/internal/api/admin/request"
	"takeout/internal/api/admin/response"
	"takeout/internal/model"
	"takeout/repository"
)

type IEmployeeService interface {
	Login(context.Context, request.EmployeeLogin) (*response.EmployeeLogin, error)
	Logout(ctx context.Context) error
	EditPassword(context.Context, request.EmployeeEditPassword) error
	CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error
	PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error)
	SetStatus(ctx context.Context, id uint64, status int) error
	UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error
	GetById(ctx context.Context, id uint64) (model.Employee, error)
}

type EmployeeImpl struct {
	repo repository.EmployeeRepo
}

// NewEmployeeService initializes the employee service
func NewEmployeeService(repo repository.EmployeeRepo) IEmployeeService {
	return &EmployeeImpl{repo: repo}
}

func (ei *EmployeeImpl) GetById(ctx context.Context, id uint64) (model.Employee, error) {
	employee, err := ei.repo.GetById(ctx, id)
	if err != nil {
		global.Log.Error("Failed to get employee by ID:", err.Error())
		return model.Employee{}, err
	}

	// Mask sensitive data
	employee.Password = "****"
	return *employee, nil
}

func (ei *EmployeeImpl) Login(ctx context.Context, employeeLogin request.EmployeeLogin) (*response.EmployeeLogin, error) {
	employee, err := ei.repo.GetByUserName(ctx, employeeLogin.UserName)
	if err != nil || employee == nil {
		global.Log.Warn("Login failed: account not found for username:", employeeLogin.UserName)
		return nil, e.Error_ACCOUNT_NOT_FOUND
	}

	// Validate password
	password := utils.MD5V(employeeLogin.Password, "", 0)
	if password != employee.Password {
		global.Log.Warn("Login failed: incorrect password for username:", employeeLogin.UserName)
		return nil, e.Error_PASSWORD_ERROR
	}

	if employee.Status == enum.DISABLE {
		global.Log.Warn("Login failed: account locked for username:", employeeLogin.UserName)
		return nil, e.Error_ACCOUNT_LOCKED
	}

	// Generate JWT token
	jwtConfig := global.Config.Jwt.Admin
	token, err := utils.GenerateToken(employee.Id, jwtConfig.Name, jwtConfig.Secret)
	if err != nil {
		global.Log.Error("Failed to generate JWT token for username:", employeeLogin.UserName, err.Error())
		return nil, err
	}

	resp := &response.EmployeeLogin{
		Id:       employee.Id,
		Name:     employee.Name,
		Token:    token,
		UserName: employee.Username,
	}
	global.Log.Info("Login successful for username:", employeeLogin.UserName)
	return resp, nil
}

func (ei *EmployeeImpl) Logout(ctx context.Context) error {
	// Frontend deletes JWT token
	// No backend action needed
	global.Log.Info("Employee logged out successfully.")
	return nil
}

func (ei *EmployeeImpl) CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error {
	// Creating new employee entity
	entity := model.Employee{
		Id:       employee.Id,
		Username: employee.UserName,
		Name:     employee.Name,
		Phone:    employee.Phone,
		Sex:      employee.Sex,
		IdNumber: employee.IdNumber,
		Status:   enum.ENABLE,
		Password: utils.MD5V("123456", "", 0), // Default password, should be securely changed later
	}

	err := ei.repo.Insert(ctx, entity)
	if err != nil {
		global.Log.Error("Failed to create employee:", err.Error())
		return err
	}

	global.Log.Info("Employee created successfully:", employee.UserName)
	return nil
}

func (ei *EmployeeImpl) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
	// Paginated query
	pageResult, err := ei.repo.PageQuery(ctx, dto)
	if err != nil {
		global.Log.Error("Failed to query employee page:", err.Error())
		return nil, err
	}

	// Mask sensitive data
	if employeeList, ok := pageResult.Records.([]model.Employee); ok {
		for key := range employeeList {
			employeeList[key].Password = "****"
			employeeList[key].IdNumber = "****"
			employeeList[key].Phone = "****"
		}
		pageResult.Records = employeeList
	}
	global.Log.Info("Employee page query successful.")
	return pageResult, nil
}

func (ei *EmployeeImpl) SetStatus(ctx context.Context, id uint64, status int) error {
	entity := model.Employee{
		Id:     id,
		Status: status,
	}
	err := ei.repo.UpdateStatus(ctx, entity)
	if err != nil {
		global.Log.Error("Failed to set employee status:", err.Error())
		return err
	}

	global.Log.Info("Employee status updated successfully:", id, "Status:", status)
	return nil
}

func (ei *EmployeeImpl) UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error {
	entity := model.Employee{
		Id:       dto.Id,
		Username: dto.UserName,
		Name:     dto.Name,
		Phone:    dto.Phone,
		Sex:      dto.Sex,
		IdNumber: dto.IdNumber,
	}
	err := ei.repo.Update(ctx, entity)
	if err != nil {
		global.Log.Error("Failed to update employee:", err.Error())
		return err
	}

	global.Log.Info("Employee updated successfully:", dto.UserName)
	return nil
}

func (ei *EmployeeImpl) EditPassword(ctx context.Context, employeeEdit request.EmployeeEditPassword) error {
	employee, err := ei.repo.GetById(ctx, employeeEdit.EmpId)
	if err != nil || employee == nil {
		global.Log.Warn("Failed to find employee for password edit:", employeeEdit.EmpId)
		return e.Error_ACCOUNT_NOT_FOUND
	}

	oldPassword := utils.MD5V(employeeEdit.OldPassword, "", 0)
	if oldPassword != employee.Password {
		global.Log.Warn("Incorrect old password for employee ID:", employeeEdit.EmpId)
		return e.Error_PASSWORD_ERROR
	}

	newPassword := utils.MD5V(employeeEdit.NewPassword, "", 0)
	err = ei.repo.Update(ctx, model.Employee{
		Id:       employee.Id,
		Password: newPassword,
	})
	if err != nil {
		global.Log.Error("Failed to update password for employee ID:", employeeEdit.EmpId)
		return err
	}

	global.Log.Info("Password updated successfully for employee ID:", employeeEdit.EmpId)
	return nil
}

//package service
//
//import (
//	"context"
//	"takeout/common"
//	"takeout/common/e"
//	"takeout/common/enum"
//	"takeout/common/utils"
//	"takeout/global"
//	"takeout/internal/api/admin/request"
//	"takeout/internal/api/admin/response"
//	"takeout/internal/model"
//	"takeout/repository"
//)
//
//type IEmployeeService interface {
//	Login(context.Context, request.EmployeeLogin) (*response.EmployeeLogin, error)
//	Logout(ctx context.Context) error
//	EditPassword(context.Context, request.EmployeeEditPassword) error
//	CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error
//	PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error)
//	SetStatus(ctx context.Context, id uint64, status int) error
//	UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error
//	GetById(ctx context.Context, id uint64) (model.Employee, error)
//}
//
//func (ei *EmployeeImpl) GetById(ctx context.Context, id uint64) (model.Employee, error) {
//	employee, err := ei.repo.GetById(ctx, id)
//	employee.Password = "****"
//	return *employee, err
//}
//
//func (ei *EmployeeImpl) Login(ctx context.Context, employeeLogin request.EmployeeLogin) (*response.EmployeeLogin, error) {
//	employee, err := ei.repo.GetByUserName(ctx, employeeLogin.UserName)
//	if err != nil || employee == nil {
//		return nil, e.Error_ACCOUNT_NOT_FOUND
//	}
//	password := utils.MD5V(employeeLogin.Password, "", 0)
//	if password != employee.Password {
//		return nil, e.Error_PASSWORD_ERROR
//	}
//	if employee.Status == enum.DISABLE {
//		return nil, e.Error_ACCOUNT_LOCKED
//	}
//
//	// 返回jwt
//	jwtConfig := global.Config.Jwt.Admin
//	token, err := utils.GenerateToken(employee.Id, jwtConfig.Name, jwtConfig.Secret)
//	if err != nil {
//		return nil, err
//	}
//
//	resp := &response.EmployeeLogin{
//		Id:       employee.Id,
//		Name:     employee.Name,
//		Token:    token,
//		UserName: employee.Username,
//	}
//	return resp, err
//}
//
//func (ei *EmployeeImpl) Logout(ctx context.Context) (err error) {
//	// 前端删除jwt头数据
//	return
//}
//
//func (ei *EmployeeImpl) CreateEmployee(ctx context.Context, employee request.EmployeeDTO) (err error) {
//	entity := model.Employee{
//		Id:       employee.Id,
//		Username: employee.UserName,
//		Name:     employee.Name,
//		Phone:    employee.Phone,
//		Sex:      employee.Sex,
//		IdNumber: employee.IdNumber,
//		Status:   enum.ENABLE,
//		Password: utils.MD5V("123456", "", 0),
//	}
//	err = ei.repo.Insert(ctx, entity)
//	return err
//}
//
//func (ei *EmployeeImpl) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
//	// 分页查询
//	pageResult, err := ei.repo.PageQuery(ctx, dto)
//
//	// 屏蔽敏感信息
//	if employeeList, ok := pageResult.Records.([]model.Employee); ok {
//		for key, _ := range employeeList {
//			employeeList[key].Password = "****"
//			employeeList[key].IdNumber = "****"
//			employeeList[key].Phone = "****"
//		}
//		pageResult.Records = employeeList
//	}
//	return pageResult, err
//}
//
//func (ei *EmployeeImpl) SetStatus(ctx context.Context, id uint64, status int) error {
//
//	entity := model.Employee{
//		Id:     id,
//		Status: status,
//	}
//	err := ei.repo.UpdateStatus(ctx, entity)
//	return err
//}
//
//func (ei *EmployeeImpl) UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error {
//	entity := model.Employee{
//		Id:       dto.Id,
//		Username: dto.UserName,
//		Name:     dto.Name,
//		Phone:    dto.Phone,
//		Sex:      dto.Sex,
//		IdNumber: dto.IdNumber,
//	}
//	err := ei.repo.Update(ctx, entity)
//	return err
//}
//
//func (ei *EmployeeImpl) EditPassword(ctx context.Context, employeeEdit request.EmployeeEditPassword) error {
//	employee, err := ei.repo.GetById(ctx, employeeEdit.EmpId)
//	if err != nil {
//		return err
//	}
//	if employee == nil {
//		return e.Error_ACCOUNT_NOT_FOUND
//	}
//	oldPassword := utils.MD5V(employeeEdit.OldPassword, "", 0)
//	if oldPassword != employee.Password {
//		return e.Error_PASSWORD_ERROR
//	}
//	newPassword := utils.MD5V(employee.Password, "", 0)
//	err = ei.repo.Update(ctx, model.Employee{
//		Id:       employee.Id,
//		Password: newPassword,
//	})
//	return err
//}
//
//type EmployeeImpl struct {
//	repo repository.EmployeeRepo
//}
//
//func NewEmployeeService(repo repository.EmployeeRepo) IEmployeeService {
//	return &EmployeeImpl{repo: repo}
//}
