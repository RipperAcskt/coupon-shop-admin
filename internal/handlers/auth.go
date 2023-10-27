package handlers

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/RipperAcskt/coupon-shop-admin/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	SingIn(authInfo entities.Auth) (*entities.Token, error)
	Verify(token string) error
	Logout(token string, expired time.Duration) error
	CheckToken(userId string) bool
}

func (handlers Handlers) SingIn(c *gin.Context) {

	var auth entities.Auth

	if err := c.BindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println(auth)

	token, err := handlers.svc.SingIn(auth)
	if err != nil {
		if errors.Is(err, entities.ErrWrongLoginOrPassword) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("sing in failed")
		return
	}

	exp := int((time.Duration(handlers.cfg.RefreshTokenExp) * time.Hour * 24).Seconds())
	c.SetCookie("refresh_token", token.RT, exp, "/admin/auth", "parcus.shop", false, true)
	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Access,
	})
}

func (handlers Handlers) VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Split(c.GetHeader("Authorization"), " ")
		if len(token) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Errorf("access token required").Error(),
			})
			return
		}
		accessToken := token[1]

		err := handlers.svc.Verify(accessToken)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": err.Error(),
				})
				return
			}
			if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": fmt.Errorf("wrong signature").Error(),
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("wrong token").Error(),
			})
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("verify failed")
			return
		}

		if !handlers.svc.CheckToken(accessToken) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "you already logged out",
			})
			return
		}

		c.Next()

	}
}

func (handlers Handlers) Refresh(c *gin.Context) {

	refresh, err := c.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("bad refresh token").Error(),
			})
			return
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get cookies failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = handlers.svc.Verify(refresh)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": fmt.Errorf("wrong signature").Error(),
			})
			return
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("verify failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("verify rt failed: %w", err).Error(),
		})
		return
	}

	params := service.TokenParams{
		Type:            service.Admin,
		Hs256Secret:     handlers.cfg.Hs256Secret,
		AccessTokenExp:  handlers.cfg.AccessTokenExp,
		RefreshTokenExp: handlers.cfg.RefreshTokenExp,
	}

	token, err := service.NewToken(params)
	if err != nil {
		if errors.Is(err, service.ErrUnknownType) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("new token failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	exp := int((time.Duration(handlers.cfg.RefreshTokenExp) * time.Hour * 24).Seconds())
	c.SetCookie("refresh_token", token.RT, exp, "/admin/auth", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"access_token": token.Access,
	})
}

func (handlers Handlers) Logout(c *gin.Context) {
	exp := time.Duration(handlers.cfg.AccessTokenExp) * time.Minute
	token := strings.Split(c.GetHeader("Authorization"), " ")
	if len(token) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Errorf("access token required").Error(),
		})
		return
	}
	accessToken := token[1]

	err := handlers.svc.Logout(accessToken, exp)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("logout failed")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.SetCookie("refresh_token", "", time.Now().Second(), "/admin/auth", "", false, true)
	c.Status(http.StatusOK)
}
