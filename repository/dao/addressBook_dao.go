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

func (a AddressBookDao) CreateAddress(ctx context.Context, addressBook model.AddressBook) error {
	return a.db.WithContext(ctx).Create(&addressBook).Error
}

func (a AddressBookDao) GetAddressById(ctx context.Context, id uint64) (model.AddressBook, error) {
	var (
		addressBook model.AddressBook
		err         error
	)
	if err = a.db.WithContext(ctx).Where("id=?", id).First(&addressBook).Error; err != nil && err != gorm.ErrRecordNotFound {
		return model.AddressBook{}, err
	}
	return addressBook, nil
}

func (a AddressBookDao) UpdateAddressById(ctx context.Context, addressBook model.AddressBook) error {
	return a.db.WithContext(ctx).Updates(&addressBook).Error
}

func (a AddressBookDao) DeleteById(ctx context.Context, id uint64) error {
	return a.db.WithContext(ctx).Where("id=?", id).Delete(&model.AddressBook{}).Error
}

func (a AddressBookDao) GetCurAddressBook(ctx context.Context, id uint64) (model.AddressBook, error) {
	//TODO implement me
	panic("implement me")
}

func (a AddressBookDao) GetDefaultAddressBook(ctx context.Context, id uint64) (model.AddressBook, error) {
	//TODO implement me
	panic("implement me")
}

func (a AddressBookDao) SetDefaultAddressBook(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}
