package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/common"
	"takeout/internal/api/admin/request"
	"takeout/internal/model"
)

type EmployeeDao struct {
	db *gorm.DB
}

func NewEmployeeDao(db *gorm.DB) *EmployeeDao {
	return &EmployeeDao{db: db}
}

func (e *EmployeeDao) Insert(ctx context.Context, entity model.Employee) error {
	return e.db.WithContext(ctx).Create(&entity).Error
}

func (e *EmployeeDao) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
	var (
		err          error
		employeeList []model.Employee
		result       common.PageResult
	)
	query := e.db.WithContext(ctx).Model(&model.Employee{})
	if dto.Name != "" {
		query = query.Where("name LIKE ?", "%"+dto.Name+"%")
	}
	if err = query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	// 分页操作
	err = query.Scopes(result.Paginate(&dto.Page, &dto.PageSize)).Find(&employeeList).Error
	result.Records = employeeList
	return &result, err
}

func (e *EmployeeDao) GetByUserName(ctx context.Context, userName string) (*model.Employee, error) {
	var employee model.Employee
	err := e.db.WithContext(ctx).Where("username=?", userName).First(&employee).Error
	return &employee, err
}

func (e *EmployeeDao) GetById(ctx context.Context, id uint64) (*model.Employee, error) {
	var employee model.Employee
	err := e.db.WithContext(ctx).Where("id=?", id).First(&employee).Error
	return &employee, err
}

func (e *EmployeeDao) UpdateStatus(ctx context.Context, employee model.Employee) error {
	return e.db.WithContext(ctx).Model(&model.Employee{}).Where("id=?", employee.Id).Update("status", employee.Status).Error
}

func (e *EmployeeDao) Update(ctx context.Context, employee model.Employee) error {
	err := e.db.WithContext(ctx).Model(&employee).Select("password").Updates(employee).Error
	return err
}
