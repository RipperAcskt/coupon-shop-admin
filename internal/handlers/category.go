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

type CategoryService interface {
	CreateCategory(ctx context.Context, category entities.Category) error
	GetCategories(ctx context.Context, subcategory bool) ([]entities.Category, error)
	GetCategory(ctx context.Context, id string, subcategory bool) (entities.Category, error)
	UpdateCategory(ctx context.Context, id, name string, subcategory bool) error
	DeleteCategory(ctx context.Context, id string, subcategory bool) error
}

func (handlers Handlers) createCategory(context *gin.Context) {
	category := entities.Category{}
	err := context.ShouldBindJSON(&category)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.CreateCategory(context, category)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create category failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create category failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "category is successfully created",
	})
}

func (handlers Handlers) getCategories(context *gin.Context) {
	var sub bool
	subcategory := context.GetHeader("Subcategory")
	if subcategory == "true" {
		sub = true
	}

	subs, err := handlers.svc.GetCategories(context, sub)
	if err != nil {
		if errors.Is(err, entities.ErrNoAnyCategory) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get categories failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get categories failed")
		return
	}
	context.JSON(http.StatusCreated, subs)
}

func (handlers Handlers) getCategory(context *gin.Context) {
	var sub bool
	subcategory := context.GetHeader("Subcategory")
	if subcategory == "true" {
		sub = true
	}

	id := context.Param("id")
	subs, err := handlers.svc.GetCategory(context, id, sub)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get category failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get category failed")
		return
	}
	context.JSON(http.StatusCreated, subs)
}

func (handlers Handlers) updateCategory(context *gin.Context) {
	category := entities.Category{}
	err := context.ShouldBindJSON(&category)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	var sub bool
	subcategory := context.GetHeader("Subcategory")
	if subcategory == "true" {
		sub = true
	}
	id := context.Param("id")

	err = handlers.svc.UpdateCategory(context, id, category.Name, sub)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("update category failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update category failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "category is successfully updated",
	})
}

func (handlers Handlers) deleteCategory(context *gin.Context) {
	var sub bool
	subcategory := context.GetHeader("Subcategory")
	if subcategory == "true" {
		sub = true
	}

	id := context.Param("id")
	err := handlers.svc.DeleteCategory(context, id, sub)
	if err != nil {
		if errors.Is(err, entities.ErrCategoryDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("delete category failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete category failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "category is successfully deleted",
	})
}
