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

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, sub entities.Subscription) error
	GetSubscriptions(ctx context.Context) ([]entities.Subscription, error)
	GetSubscription(ctx context.Context, id string) (entities.Subscription, error)
	UpdateSubscription(ctx context.Context, id string, subscription entities.Subscription) error
	DeleteSubscription(ctx context.Context, id string) error
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
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

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

func (handlers Handlers) getSubscriptions(context *gin.Context) {
	subs, err := handlers.svc.GetSubscriptions(context)
	if err != nil {
		if errors.Is(err, entities.ErrNoAnySubscription) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get subscriptions failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get subscriptions failed")
		return
	}
	context.JSON(http.StatusCreated, subs)
}

func (handlers Handlers) getSubscription(context *gin.Context) {
	id := context.Param("id")
	subs, err := handlers.svc.GetSubscription(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get subscription failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get subscription failed")
		return
	}
	context.JSON(http.StatusCreated, subs)
}

func (handlers Handlers) updateSubscription(context *gin.Context) {
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

	id := context.Param("id")
	err = handlers.svc.UpdateSubscription(context, id, subscription)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create subscription failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update subscription failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "subscription is successfully updated",
	})
}

func (handlers Handlers) deleteSubscription(context *gin.Context) {
	id := context.Param("id")
	err := handlers.svc.DeleteSubscription(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create subscription failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete subscription failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "subscription is successfully deleted",
	})
}
