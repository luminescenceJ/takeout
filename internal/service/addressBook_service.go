package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"takeout/common/enum"
	"takeout/internal/model"
	"takeout/repository"
)

type IAddressBookService interface {
	CreateAddressBook(ctx *gin.Context, address model.AddressBook) error
	UpdateAddressBook(ctx *gin.Context, address model.AddressBook) error
	DeleteAddressBook(ctx *gin.Context, id uint64) error
	GetAddressBook(ctx *gin.Context, id uint64) (model.AddressBook, error)
	GetCurAddressBook(ctx *gin.Context) (model.AddressBook, error)
	GetDefaultAddressBook(ctx *gin.Context) (model.AddressBook, error)
	SetDefaultAddressBook(ctx *gin.Context, id uint64) error
}

type AddressBookService struct {
	repo repository.AddressBookRepo
}

func NewAddressBookService(repo repository.AddressBookRepo) IAddressBookService {
	return &AddressBookService{
		repo: repo,
	}
}

func (as *AddressBookService) CreateAddressBook(ctx *gin.Context, address model.AddressBook) error {
	return as.repo.CreateAddress(ctx, address)
}
func (as *AddressBookService) UpdateAddressBook(ctx *gin.Context, address model.AddressBook) error {
	return as.repo.UpdateAddressById(ctx, address)
}
func (as *AddressBookService) DeleteAddressBook(ctx *gin.Context, id uint64) error {
	return as.repo.DeleteById(ctx, id)
}
func (as *AddressBookService) GetAddressBook(ctx *gin.Context, id uint64) (model.AddressBook, error) {
	return as.repo.GetAddressById(ctx, id)
}

func (as *AddressBookService) GetCurAddressBook(ctx *gin.Context) (model.AddressBook, error) {
	var (
		id      uint64
		address model.AddressBook
		err     error
	)
	if CurrentId, ok := ctx.Get(enum.CurrentId); ok {
		id = CurrentId.(uint64)
	} else {
		return model.AddressBook{}, errors.New("current id book not found")
	}
	if address, err = as.repo.GetCurAddressBook(ctx, id); err != nil {
		return model.AddressBook{}, err
	}
	return address, nil
}
func (as *AddressBookService) GetDefaultAddressBook(ctx *gin.Context) (model.AddressBook, error) {
	var (
		id      uint64
		address model.AddressBook
		err     error
	)
	if CurrentId, ok := ctx.Get(enum.CurrentId); ok {
		id = CurrentId.(uint64)
	} else {
		return model.AddressBook{}, errors.New("current id book not found")
	}
	if address, err = as.repo.GetDefaultAddressBook(ctx, id); err != nil {
		return model.AddressBook{}, err
	}
	return address, nil
}
func (as *AddressBookService) SetDefaultAddressBook(ctx *gin.Context, id uint64) error {
	return as.repo.SetDefaultAddressBook(ctx, id)
}
