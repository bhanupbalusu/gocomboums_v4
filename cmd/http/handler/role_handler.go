package handler

import (
	"net/http"
	"strconv"

	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"github.com/bhanupbalusu/gocomboums_v4/internal/service"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	RoleService        *service.RoleService
	TransactionHandler *TransactionHandler
}

func NewRoleHandler(roleService *service.RoleService, txHandler *TransactionHandler) *RoleHandler {
	return &RoleHandler{
		RoleService:        roleService,
		TransactionHandler: txHandler,
	}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.TransactionHandler.StartTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			_ = h.TransactionHandler.RollbackTransaction()
		}
	}()

	newRole, err := h.RoleService.CreateRole(&role)
	if err != nil {
		_ = h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.TransactionHandler.CommitTransaction()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newRole)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedRole, err := h.RoleService.UpdateRole(&role)
	if err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRole)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.RoleService.DeleteRole(roleID); err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Role deleted"})
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.RoleService.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.RoleService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) GetRolesByUserID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roles, err := h.RoleService.GetRolesByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *RoleHandler) UserHasRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleName := c.Param("roleName")

	hasRole, err := h.RoleService.UserHasRole(userID, roleName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hasRole": hasRole})
}

func (h *RoleHandler) AddUserRole(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"userID"`
		RoleID uint64 `json:"roleID"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.RoleService.AddUserRole(req.UserID, req.RoleID); err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Role added to user"})
}

func (h *RoleHandler) RemoveUserRole(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"userID"`
		RoleID uint64 `json:"roleID"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.RoleService.RemoveUserRole(req.UserID, req.RoleID); err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Role removed from user"})
}
