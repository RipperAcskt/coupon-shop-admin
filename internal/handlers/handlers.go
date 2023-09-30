package handlers

import (
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	svc Service
	cfg config.Config
}

func NewAdminHandlers(service Service, cfg config.Config) *Handlers {
	return &Handlers{svc: service, cfg: cfg}
}

type Service interface {
	MemberService
	OrganizationService
	SubscriptionService
	CouponService
	AuthService
}

func SetRequestHandlers(service Service, cfg config.Config) (*gin.Engine, error) {
	router := gin.New()
	handlers := NewAdminHandlers(service, cfg)
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(CORSMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello message")
	})
	organization := router.Group("/organization")
	{
		organization.Use(handlers.VerifyToken())

		organization.POST("/", handlers.createOrganization)
		organization.GET("/", handlers.getOrganizations)
		organization.GET("/:id", handlers.getOrganization)
		organization.DELETE("/:id", handlers.deleteOrganization)
	}
	members := organization.Group("/members")
	{
		members.Use(handlers.VerifyToken())

		members.POST("/:id", handlers.addMembers)
		//members.GET("/", handlers.ge)
		//members.DELETE("/", handlers.deleteOrganizationMembers)

	}
	subscription := router.Group("/subscription")
	{
		subscription.Use(handlers.VerifyToken())

		subscription.POST("/", handlers.createSubscription)
		subscription.GET("/", handlers.getSubscriptions)
		subscription.GET("/:id", handlers.getSubscription)
		subscription.PUT("/:id", handlers.updateSubscription)
		subscription.DELETE("/:id", handlers.deleteSubscription)
	}

	coupon := router.Group("/coupon")
	{
		coupon.Use(handlers.VerifyToken())

		coupon.POST("/", handlers.createCoupon)
		coupon.GET("/", handlers.getCoupons)
		coupon.GET("/:id", handlers.getCoupon)
		coupon.PUT("/:id", handlers.updateCoupon)
		coupon.DELETE("/:id", handlers.deleteCoupon)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sing-in", handlers.SingIn)
		auth.GET("/refresh", handlers.Refresh)
		auth.GET("/logout", handlers.Logout)
	}

	router.GET("/store/:id", handlers.VerifyToken(), handlers.getContent)
	return router, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
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
