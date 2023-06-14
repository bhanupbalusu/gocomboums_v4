package handler

import (
	"net/http"
	"strconv"

	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"github.com/bhanupbalusu/gocomboums_v4/internal/service"
	"github.com/gin-gonic/gin"
)

type RolePermissions struct {
	RoleID        uint64   `json:"role_id"`
	PermissionIDs []uint64 `json:"permission_ids"`
}

type PermissionHandler struct {
	PermissionService  *service.PermissionService
	TransactionHandler *TransactionHandler
}

func NewPermissionHandler(permissionService *service.PermissionService, txHandler *TransactionHandler) *PermissionHandler {
	return &PermissionHandler{
		PermissionService:  permissionService,
		TransactionHandler: txHandler,
	}
}

// CreatePermission handles the request to create a new permission.
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	permission, err := h.PermissionService.CreatePermission(permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

// UpdatePermission handles the request to update an existing permission.
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	permission, err := h.PermissionService.UpdatePermission(permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

// DeletePermission handles the request to delete an existing permission.
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	permissionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err := h.PermissionService.DeletePermission(permissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Permission deleted successfully."})
}

func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.PermissionService.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, permissions)
}

func (h *PermissionHandler) GetPermissionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission id"})
		return
	}

	permission, err := h.PermissionService.GetPermissionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// AssignPermissionToRole handles the request to assign a permission to a role.
func (h *PermissionHandler) AssignPermissionToRole(c *gin.Context) {
	var rolePermission model.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.PermissionService.AssignPermissionToRole(rolePermission.RoleID, rolePermission.PermissionID)
	if err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to role successfully."})
}

// RemovePermissionFromRole handles the request to remove a permission from a role.
func (h *PermissionHandler) RemovePermissionFromRole(c *gin.Context) {
	var rolePermission model.RolePermission
	if err := c.ShouldBindJSON(&rolePermission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.PermissionService.RemovePermissionFromRole(rolePermission.RoleID, rolePermission.PermissionID)
	if err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission removed from role successfully."})
}

// AddMultiplePermissionsToRole handles the request to add multiple permissions to a role.
func (h *PermissionHandler) AddMultiplePermissionsToRole(c *gin.Context) {
	var rolePermissions RolePermissions
	if err := c.ShouldBindJSON(&rolePermissions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.PermissionService.AddMultiplePermissionsToRole(rolePermissions.RoleID, rolePermissions.PermissionIDs)
	if err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions added to role successfully."})
}

// RemoveMultiplePermissionsFromRole handles the request to remove multiple permissions from a role.
func (h *PermissionHandler) RemoveMultiplePermissionsFromRole(c *gin.Context) {
	var rolePermissions RolePermissions
	if err := c.ShouldBindJSON(&rolePermissions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.StartTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.PermissionService.RemoveMultiplePermissionsFromRole(rolePermissions.RoleID, rolePermissions.PermissionIDs)
	if err != nil {
		h.TransactionHandler.RollbackTransaction()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.TransactionHandler.CommitTransaction(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions removed from role successfully."})
}
