package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	OrganizationService
}

func NewAdminHandlers(organizationService OrganizationService) *Handlers {
	return &Handlers{OrganizationService: organizationService}
}

type OrganizationService interface {
	CreateOrganization()
}

func SetRequestHandlers(service OrganizationService) (*gin.Engine, error) {
	router := gin.New()
	organizationHandlers := NewAdminHandlers(service)
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(CORSMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello message")
	})
	organization := router.Group("/organization")
	{
		organization.POST("/organization", organizationHandlers.createOrganization)
		organization.GET("/organization", organizationHandlers.getOrganizations)
		organization.DELETE("/organization", organizationHandlers.deleteOrganization)
	}
	members := organization.Group("/members")
	{
		members.POST("/members", organizationHandlers.addOrganizationMembers)
		members.GET("/members", organizationHandlers.getOrganizationMembers)
		members.DELETE("/members", organizationHandlers.deleteOrganizationMembers)

	}
	return router, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		fmt.Println(c.Request.Method)

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
