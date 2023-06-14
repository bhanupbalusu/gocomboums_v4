package handler

import (
	"github.com/bhanupbalusu/gocomboums_v4/internal/service"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	UserService       *service.UserService
	RoleService       *service.RoleService
	PermissionService *service.PermissionService
	DB                *gorm.DB
	Tx                *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{
		DB: db,
	}
}

func (h *TransactionHandler) StartTransaction() error {
	h.Tx = h.DB.Begin()

	if h.Tx.Error != nil {
		return errors.Wrap(h.Tx.Error, "start transaction failed")
	}
	return nil

}

func (h *TransactionHandler) CommitTransaction() error {
	if h.Tx == nil {
		return errors.New("no transaction found")
	}
	h.Tx = h.Tx.Commit()
	if h.Tx.Error != nil {
		return errors.Wrap(h.Tx.Error, "commit transaction failed")
	}
	return nil

}

func (h *TransactionHandler) RollbackTransaction() error {
	if h.Tx == nil {
		return errors.New("no transaction found")
	}
	h.Tx = h.Tx.Rollback()
	if h.Tx.Error != nil {
		return errors.Wrap(h.Tx.Error, "rollback transaction failed")
	}
	return nil

}
