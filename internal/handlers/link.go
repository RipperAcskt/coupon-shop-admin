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

type LinkService interface {
	CreateLink(ctx context.Context, link entities.Link) error
	GetLinks(ctx context.Context) ([]entities.Link, error)
	GetLinkByRegion(ctx context.Context, region string) (entities.Link, error)
	UpdateLink(ctx context.Context, id, link string) error
	DeleteLink(ctx context.Context, id string) error
}

func (handlers Handlers) createLink(context *gin.Context) {
	link := entities.Link{}
	err := context.ShouldBindJSON(&link)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.CreateLink(context, link)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create link failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create link failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "link is successfully created",
	})
}

func (handlers Handlers) getLinks(context *gin.Context) {
	subs, err := handlers.svc.GetLinks(context)
	if err != nil {
		if errors.Is(err, entities.ErrNoAnyCategory) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get links failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get links failed")
		return
	}
	context.JSON(http.StatusOK, subs)
}

func (handlers Handlers) getLink(context *gin.Context) {
	region := context.Param("region")
	subs, err := handlers.svc.GetLinkByRegion(context, region)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get link failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get link failed")
		return
	}
	context.JSON(http.StatusOK, subs)
}

func (handlers Handlers) updateLink(context *gin.Context) {
	link := entities.Link{}
	err := context.ShouldBindJSON(&link)
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

	err = handlers.svc.UpdateLink(context, id, link.Link)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("update link failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update link failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "link is successfully updated",
	})
}

func (handlers Handlers) deleteLink(context *gin.Context) {
	id := context.Param("id")
	err := handlers.svc.DeleteLink(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("delete link failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete link failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "link is successfully deleted",
	})
}
