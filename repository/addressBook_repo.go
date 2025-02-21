package repository

import (
	"context"
	"takeout/internal/model"
)

type AddressBookRepo interface {
	CreateAddress(ctx context.Context, addressBook model.AddressBook) error
	GetAddressById(ctx context.Context, id uint64) (model.AddressBook, error)
	UpdateAddressById(ctx context.Context, addressBook model.AddressBook) error
	DeleteById(ctx context.Context, id uint64) error
	GetCurAddressBook(ctx context.Context, userId uint64) ([]model.AddressBook, error)
	GetDefaultAddressBook(ctx context.Context, userId uint64) (model.AddressBook, error)
	SetDefaultAddressBook(ctx context.Context, userId, id uint64) error
}
