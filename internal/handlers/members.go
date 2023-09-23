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
	AddMembers(ctx context.Context, members []entities.Member, id string) error
	GetMembers(ctx context.Context, organizationID string) ([]entities.Member, error)
	DeleteMembers(ctx context.Context, membersToDelete []entities.Member) error
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
