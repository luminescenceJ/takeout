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

// Insert inserts a new employee into the database
func (e *EmployeeDao) Insert(ctx context.Context, entity model.Employee) error {
	return e.db.WithContext(ctx).Create(&entity).Error
}

// PageQuery performs paginated query for employees with filter options
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
	// Count total records for pagination
	if err = query.Count(&result.Total).Error; err != nil {
		return nil, err
	}
	// Paginate query
	err = query.Scopes(result.Paginate(&dto.Page, &dto.PageSize)).Find(&employeeList).Error
	if err != nil {
		return nil, err
	}

	// Mask sensitive fields
	for i := range employeeList {
		employeeList[i].Password = "****" // Mask password
	}
	result.Records = employeeList
	return &result, nil
}

// GetByUserName retrieves an employee by username
func (e *EmployeeDao) GetByUserName(ctx context.Context, userName string) (*model.Employee, error) {
	var employee model.Employee
	err := e.db.WithContext(ctx).Where("username = ?", userName).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetById retrieves an employee by ID
func (e *EmployeeDao) GetById(ctx context.Context, id uint64) (*model.Employee, error) {
	var employee model.Employee
	err := e.db.WithContext(ctx).Where("id = ?", id).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// UpdateStatus updates the status of an employee
func (e *EmployeeDao) UpdateStatus(ctx context.Context, employee model.Employee) error {
	return e.db.WithContext(ctx).Model(&model.Employee{}).Where("id = ?", employee.Id).Update("status", employee.Status).Error
}

// Update updates an employee's information
func (e *EmployeeDao) Update(ctx context.Context, employee model.Employee) error {
	// Update only password here, consider updating other fields based on the use case
	err := e.db.WithContext(ctx).Model(&employee).Select("password").Updates(employee).Error
	return err
}

//package dao
//
//import (
//	"context"
//	"gorm.io/gorm"
//	"takeout/common"
//	"takeout/internal/api/admin/request"
//	"takeout/internal/model"
//)
//
//type EmployeeDao struct {
//	db *gorm.DB
//}
//
//func NewEmployeeDao(db *gorm.DB) *EmployeeDao {
//	return &EmployeeDao{db: db}
//}
//
//func (e *EmployeeDao) Insert(ctx context.Context, entity model.Employee) error {
//	return e.db.WithContext(ctx).Create(&entity).Error
//}
//
//func (e *EmployeeDao) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
//	var (
//		err          error
//		employeeList []model.Employee
//		result       common.PageResult
//	)
//	query := e.db.WithContext(ctx).Model(&model.Employee{})
//	if dto.Name != "" {
//		query = query.Where("name LIKE ?", "%"+dto.Name+"%")
//	}
//	if err = query.Count(&result.Total).Error; err != nil {
//		return nil, err
//	}
//	// 分页操作
//	err = query.Scopes(result.Paginate(&dto.Page, &dto.PageSize)).Find(&employeeList).Error
//	result.Records = employeeList
//	return &result, err
//}
//
//func (e *EmployeeDao) GetByUserName(ctx context.Context, userName string) (*model.Employee, error) {
//	var employee model.Employee
//	err := e.db.WithContext(ctx).Where("username=?", userName).First(&employee).Error
//	return &employee, err
//}
//
//func (e *EmployeeDao) GetById(ctx context.Context, id uint64) (*model.Employee, error) {
//	var employee model.Employee
//	err := e.db.WithContext(ctx).Where("id=?", id).First(&employee).Error
//	return &employee, err
//}
//
//func (e *EmployeeDao) UpdateStatus(ctx context.Context, employee model.Employee) error {
//	return e.db.WithContext(ctx).Model(&model.Employee{}).Where("id=?", employee.Id).Update("status", employee.Status).Error
//}
//
//func (e *EmployeeDao) Update(ctx context.Context, employee model.Employee) error {
//	err := e.db.WithContext(ctx).Model(&employee).Select("password").Updates(employee).Error
//	return err
//}
