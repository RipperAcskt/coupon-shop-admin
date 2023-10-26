package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type MemberService interface {
	AddMembers(ctx context.Context, members []entities.Member, organizationID string) error
	DeleteMembers(ctx context.Context, membersToDelete []entities.Member, organizationID string) error
	UpdateMembers(ctx context.Context, members []entities.Member, organizationID string) error
}

func (handlers Handlers) addMembers(context *gin.Context) {
	id := context.Param("id")
	members := make([]entities.Member, 0)
	err := context.ShouldBindJSON(&members)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}
	err = handlers.svc.AddMembers(context, members, id)
	if err != nil {
		if errors.Is(err, entities.ErrMembersAlreadyAdded) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("adding members failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("adding members failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "members are successfully added",
	})
}

func (handlers Handlers) deleteMembers(context *gin.Context) {
	id := context.Param("id")
	members := make([]entities.Member, 0)
	err := context.ShouldBindJSON(&members)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}
	err = handlers.svc.DeleteMembers(context, members, id)
	if err != nil {
		if errors.Is(err, entities.ErrMembersAlreadyAdded) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("deleting members failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("deleting members failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "members are successfully delited",
	})
}

func (handlers Handlers) updateMembers(context *gin.Context) {
	id := context.Param("id")
	members := make([]entities.Member, 0)
	err := context.ShouldBindJSON(&members)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}
	err = handlers.svc.UpdateMembers(context, members, id)
	if err != nil {
		if errors.Is(err, entities.ErrMembersAlreadyAdded) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("updating members failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("updating members failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "members are successfully updated",
	})
}
