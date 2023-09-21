package handlers

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, sub entities.Subscription) error
}

func (handlers Handlers) createSubscription(context *gin.Context) {
	subscription := entities.NewSubscription()
	err := context.ShouldBindJSON(&subscription)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.CreateSubscription(context, subscription)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create subscription failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create subscription failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "subscription is successfully created",
	})
}
