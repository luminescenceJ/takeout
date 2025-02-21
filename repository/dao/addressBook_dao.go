package dao

import (
	"context"
	"gorm.io/gorm"
	"takeout/internal/model"
	"takeout/repository"
)

type AddressBookDao struct {
	db *gorm.DB
}

func NewAddressBookDao(db *gorm.DB) repository.AddressBookRepo {
	return &AddressBookDao{db: db}
}

// CreateAddress 新增地址簿
func (a AddressBookDao) CreateAddress(ctx context.Context, addressBook model.AddressBook) error {
	return a.db.WithContext(ctx).Create(&addressBook).Error
}

// GetAddressById 根据id查询地址
func (a AddressBookDao) GetAddressById(ctx context.Context, id uint64) (model.AddressBook, error) {
	var (
		addressBook model.AddressBook
		err         error
	)
	if err = a.db.WithContext(ctx).Table("address_book").Where("id=?", id).First(&addressBook).Error; err != nil && err != gorm.ErrRecordNotFound {
		return model.AddressBook{}, err
	}
	return addressBook, nil
}

// UpdateAddressById 根据id修改地址
func (a AddressBookDao) UpdateAddressById(ctx context.Context, addressBook model.AddressBook) error {
	return a.db.WithContext(ctx).Table("address_book").Updates(&addressBook).Error
}

// DeleteById 根据id删除地址
func (a AddressBookDao) DeleteById(ctx context.Context, id uint64) error {
	return a.db.WithContext(ctx).Table("address_book").Where("id=?", id).Delete(&model.AddressBook{}).Error
}

// GetCurAddressBook 查询登录用户所有地址
func (a AddressBookDao) GetCurAddressBook(ctx context.Context, userId uint64) ([]model.AddressBook, error) {
	var (
		addressBook []model.AddressBook
		err         error
	)
	if err = a.db.WithContext(ctx).Table("address_book").Where("user_id=?", userId).Find(&addressBook).Error; err != nil && err != gorm.ErrRecordNotFound {
		return []model.AddressBook{}, err
	}
	return addressBook, nil
}

// GetDefaultAddressBook 查询默认地址
func (a AddressBookDao) GetDefaultAddressBook(ctx context.Context, userId uint64) (model.AddressBook, error) {
	var (
		addressBook model.AddressBook
		err         error
	)
	if err = a.db.WithContext(ctx).Table("address_book").Where("user_id=? && is_default=?", userId, 1).First(&addressBook).Error; err != nil && err != gorm.ErrRecordNotFound {
		return model.AddressBook{}, err
	}
	return addressBook, nil
}

// SetDefaultAddressBook 设置默认地址
func (a AddressBookDao) SetDefaultAddressBook(ctx context.Context, userId, id uint64) error {
	var (
		err error
	)
	// 设置所有地址为非默认地址
	if err = a.db.WithContext(ctx).Table("address_book").Where("user_id=?", userId).Update("is_default", 0).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	// 将当前地址改为默认地址
	if err = a.db.WithContext(ctx).Table("address_book").Where("id=?", id).Update("is_default", 1).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}
