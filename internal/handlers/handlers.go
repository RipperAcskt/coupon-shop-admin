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
	RegionService
	CategoryService
	LinkService
	UserService
}

func SetRequestHandlers(service Service, cfg config.Config) (*gin.Engine, error) {
	router := gin.Default()
	handlers := NewAdminHandlers(service, cfg)
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(CORSMiddleware())
	admin := router.Group("/admin")

	admin.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "hello message")
	})
	organization := admin.Group("/organization")
	{
		organization.Use(handlers.VerifyToken())

		organization.POST("", handlers.createOrganization)
		organization.GET("", handlers.getOrganizations)
		organization.GET("/:id", handlers.getOrganization)
		organization.DELETE("/:id", handlers.deleteOrganization)
		organization.PUT("/:id", handlers.updateOrganization)

	}
	members := organization.Group("/members")
	{
		members.Use(handlers.VerifyToken())

		members.POST("/:id", handlers.addMembers)
		members.DELETE("/:id", handlers.deleteMembers)
		members.PUT("/:id", handlers.updateMembers)

	}
	subscription := admin.Group("/subscription")
	{
		subscription.Use(handlers.VerifyToken())

		subscription.POST("", handlers.createSubscription)
		subscription.GET("", handlers.getSubscriptions)
		subscription.GET("/:id", handlers.getSubscription)
		subscription.PUT("/:id", handlers.updateSubscription)
		subscription.DELETE("/:id", handlers.deleteSubscription)
	}

	coupon := admin.Group("/coupon")
	{
		coupon.Use(handlers.VerifyToken())

		coupon.POST("", handlers.createCoupon)
		coupon.GET("", handlers.getCoupons)
		coupon.GET("/:id", handlers.getCoupon)
		coupon.GET("filter/:region", handlers.getCouponsByRegion)
		coupon.PUT("/:id", handlers.updateCoupon)
		coupon.DELETE("/:id", handlers.deleteCoupon)
	}

	auth := admin.Group("/auth")
	{
		auth.POST("/sing-in", handlers.SingIn)
		auth.GET("/refresh", handlers.Refresh)
		auth.GET("/logout", handlers.Logout)
	}

	region := admin.Group("/region")
	{
		region.Use(handlers.VerifyToken())

		region.POST("", handlers.createRegion)
		region.GET("", handlers.getRegions)
		region.PUT("/:id", handlers.updateRegion)
		region.DELETE("/:id", handlers.deleteRegion)
	}

	category := admin.Group("/category")
	{
		category.Use(handlers.VerifyToken())

		category.POST("", handlers.createCategory)
		category.GET("", handlers.getCategories)
		category.GET("/:id", handlers.getCategory)
		category.PUT("/:id", handlers.updateCategory)
		category.DELETE("/:id", handlers.deleteCategory)
	}

	link := admin.Group("/link")
	{
		link.Use(handlers.VerifyToken())

		link.POST("", handlers.createLink)
		link.GET("", handlers.getLinks)
		link.GET("/:region", handlers.getLink)
		link.PUT("/:id", handlers.updateLink)
		link.DELETE("/:id", handlers.deleteLink)
	}

	user := admin.Group("/user")
	{
		user.Use(handlers.VerifyToken())

		user.POST("", handlers.createUser)
		user.GET("", handlers.getUsers)
		user.GET("/:id", handlers.getUser)
		user.PUT("/:id", handlers.updateUser)
		user.DELETE("/:id", handlers.deleteUser)
	}

	admin.GET("/store/:id", handlers.getContent)
	admin.GET("/store/organization/:id", handlers.getContent)
	return router, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:9999")
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
