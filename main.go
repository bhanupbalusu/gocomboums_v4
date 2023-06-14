package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bhanupbalusu/gocomboums_v4/cmd/http/handler"
	"github.com/bhanupbalusu/gocomboums_v4/internal/model"
	"github.com/bhanupbalusu/gocomboums_v4/internal/repository"
	"github.com/bhanupbalusu/gocomboums_v4/internal/service"
)

func main() {
	dsn := "postgresql://postgres:BCombo1Postgres!@localhost:5433/db_ums_v4?sslmode=disable" // Replace with your connection string
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Role{}, &model.UserRole{})
	db.AutoMigrate(&model.Permission{}, &model.RolePermission{})

	r := gin.Default()

	// Create User Service and User Handler
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Create Role Service and Role Handler
	roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService(roleRepo)
	txHandler1 := handler.NewTransactionHandler(db)
	roleHandler := handler.NewRoleHandler(roleService, txHandler1)

	// Create Permission Service and Permission Handler
	permissionRepo := repository.NewPermissionRepository(db)
	permissionService := service.NewPermissionService(permissionRepo)
	txHandler2 := handler.NewTransactionHandler(db)
	permissionHandler := handler.NewPermissionHandler(permissionService, txHandler2)

	// Create a group for routes which don't require authentication
	publicRoutes := r.Group("/")
	{
		publicRoutes.POST("/register", userHandler.RegisterUser)
		publicRoutes.POST("/login", userHandler.LoginUser)
	}

	// Create a group for routes which require authentication
	privateRoutes := r.Group("/")
	privateRoutes.Use(func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Parse the token
		claims, err := service.VerifyAndExtractClaims(token, service.RetrieveKeyInFile())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Store the claims in the context for later use
		c.Set("claims", claims)

		c.Next()
	})
	{
		privateRoutes.GET("/users", userHandler.ListUsers)
		privateRoutes.GET("/users/search", userHandler.SearchUsers)
		privateRoutes.GET("/users/count", userHandler.CountUsers)
		privateRoutes.GET("/user/name/:username", userHandler.GetUserByUsername)
		privateRoutes.GET("/user/:id", userHandler.GetUserByID)
		privateRoutes.PUT("/user", userHandler.UpdateUser)
		privateRoutes.DELETE("/user/:id", userHandler.DeleteUser)

		roles := privateRoutes.Group("/roles")
		{
			roles.POST("/", roleHandler.CreateRole)
			roles.PUT("/", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			roles.GET("/:id", roleHandler.GetRoleByID)
			roles.GET("/", roleHandler.GetAllRoles)
		}

		userRoles := privateRoutes.Group("/user-roles")
		{
			userRoles.POST("/", roleHandler.AddUserRole)
			userRoles.DELETE("/", roleHandler.RemoveUserRole)
		}

		users := privateRoutes.Group("/users")
		{
			users.GET("/user/:userID/has-role/:roleName", roleHandler.UserHasRole)
			users.GET("/user/:userID/roles", roleHandler.GetRolesByUserID)
		}

		// Permission related routes
		permissionGroup := privateRoutes.Group("/permission")
		{
			permissionGroup.POST("", permissionHandler.CreatePermission)
			permissionGroup.PUT("", permissionHandler.UpdatePermission)
			permissionGroup.DELETE("/:id", permissionHandler.DeletePermission)
			permissionGroup.GET("", permissionHandler.GetAllPermissions)
			permissionGroup.GET("/:id", permissionHandler.GetPermissionByID)
			permissionGroup.POST("/assign", permissionHandler.AssignPermissionToRole)
			permissionGroup.POST("/remove", permissionHandler.RemovePermissionFromRole)
			permissionGroup.POST("/assign/multiple", permissionHandler.AddMultiplePermissionsToRole)
			permissionGroup.POST("/remove/multiple", permissionHandler.RemoveMultiplePermissionsFromRole)
		}
	}

	r.Run() // listen and serve on 0.0.0.0:8080

}
