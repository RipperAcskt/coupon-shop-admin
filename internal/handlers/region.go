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

type RegionService interface {
	CreateRegion(ctx context.Context, region entities.Region) error
	GetRegions(ctx context.Context) ([]entities.Region, error)
	UpdateRegion(ctx context.Context, id string, region entities.Region) error
	DeleteRegion(ctx context.Context, id string) error
}

func (handlers Handlers) createRegion(context *gin.Context) {
	region := entities.Region{}
	err := context.ShouldBind(&region)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.CreateRegion(context, region)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create region failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create region failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "region is successfully created",
	})
}

func (handlers Handlers) getRegions(context *gin.Context) {
	regions, err := handlers.svc.GetRegions(context)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get regions failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get regions failed")
		return
	}

	context.JSON(http.StatusOK, regions)
}

func (handlers Handlers) updateRegion(context *gin.Context) {
	id := context.Param("id")

	region := entities.Region{}

	err := context.ShouldBind(&region)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.UpdateRegion(context, id, region)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("update region failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update region failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "region is successfully updated",
	})
}

func (handlers Handlers) deleteRegion(context *gin.Context) {
	id := context.Param("id")
	err := handlers.svc.DeleteRegion(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("delete region failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete region failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "region is successfully deleted",
	})
}
