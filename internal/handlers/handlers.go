package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	svc Service
}

func NewAdminHandlers(service Service) *Handlers {
	return &Handlers{svc: service}
}

type Service interface {
	MemberService
	OrganizationService
	SubscriptionService
	CouponService
}

func SetRequestHandlers(service Service) (*gin.Engine, error) {
	router := gin.New()
	handlers := NewAdminHandlers(service)
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(CORSMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello message")
	})
	organization := router.Group("/organization")
	{
		organization.POST("/", handlers.createOrganization)
		organization.GET("/", handlers.getOrganizations)
		organization.GET("/:id")
		organization.DELETE("/:id", handlers.deleteOrganization)
	}
	members := organization.Group("/:id/members")
	{
		members.POST("/", handlers.addOrganizationMembers)
		members.GET("/", handlers.getOrganizationMembers)
		members.DELETE("/", handlers.deleteOrganizationMembers)

	}
	subscription := router.Group("/subscription")
	{
		subscription.POST("/", handlers.createSubscription)
		subscription.GET("/", handlers.getSubscriptions)
		subscription.GET("/:id", handlers.getSubscription)
		subscription.PUT("/:id", handlers.updateSubscription)
		subscription.DELETE("/:id", handlers.deleteSubscription)
	}

	coupon := router.Group("/coupon")
	{
		coupon.POST("/", handlers.createCoupon)
		coupon.GET("/", handlers.getCoupons)
	}

	router.GET("/store/:id", handlers.getContent)
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
