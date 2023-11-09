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

type UserService interface {
	CreateUser(ctx context.Context, user entities.UserEntity) error
	GetUsers(ctx context.Context) ([]entities.UserEntity, error)
	GetUser(ctx context.Context, id string) (entities.UserEntity, error)
	UpdateUser(ctx context.Context, id string, user entities.UserEntity) error
	DeleteUser(ctx context.Context, id string) error
}

func (handlers Handlers) createUser(context *gin.Context) {
	user := entities.UserEntity{}
	err := context.ShouldBind(&user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.CreateUser(context, user)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create user failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create user failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "user is successfully created",
	})
}

func (handlers Handlers) getUsers(context *gin.Context) {
	users, err := handlers.svc.GetUsers(context)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get users failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get users failed")
		return
	}

	context.JSON(http.StatusOK, users)
}

func (handlers Handlers) getUser(context *gin.Context) {
	id := context.Param("id")
	user, err := handlers.svc.GetUser(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrOrganizationnDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get user failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get user failed")
		return
	}
	context.JSON(http.StatusOK, user)
}

func (handlers Handlers) updateUser(context *gin.Context) {
	id := context.Param("id")

	user := entities.UserEntity{}

	err := context.ShouldBind(&user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.UpdateUser(context, id, user)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("update user failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update user failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "user is successfully updated",
	})
}

func (handlers Handlers) deleteUser(context *gin.Context) {
	id := context.Param("id")
	err := handlers.svc.DeleteUser(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("delete user failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete user failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "user is successfully deleted",
	})
}
