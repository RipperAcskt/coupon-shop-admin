package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type CouponService interface {
	CreateCoupon(ctx context.Context, coupon entities.Coupon) error
	GetCoupons(ctx context.Context) ([]entities.Coupon, error)
	GetCoupon(ctx context.Context, id string) (entities.Coupon, error)
	UpdateCoupon(ctx context.Context, id string, coupon entities.Coupon) error
	DeleteCoupon(ctx context.Context, id string) error
}

func (handlers Handlers) createCoupon(context *gin.Context) {
	coupon := entities.NewCoupon()

	media := entities.NewMediaId(coupon.ID)

	file, err := context.FormFile("file")
	if err != nil {
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("receiving file failed: %w", err).Error(),
			})

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("receiving file failed")
			return
		}
	}
	file.Filename = media.ID
	err = context.SaveUploadedFile(file, "./store/"+media.ID+".jpg")
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("save upload file failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("save upload file failed")
		return
	}

	coupon.Media = media

	err = context.ShouldBind(&coupon)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.CreateCoupon(context, coupon)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create coupon failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create coupon failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "coupon is successfully created",
	})
}

func (handlers Handlers) getCoupons(context *gin.Context) {
	coupons, err := handlers.svc.GetCoupons(context)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get coupons failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get coupons failed")
		return
	}

	context.JSON(http.StatusOK, coupons)
}

func (handlers Handlers) getCoupon(context *gin.Context) {
	id := context.Param("id")
	coupon, err := handlers.svc.GetCoupon(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get coupon failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get coupon failed")
		return
	}
	context.JSON(http.StatusCreated, coupon)
}

func (handlers Handlers) getContent(context *gin.Context) {
	file := "." + strings.TrimPrefix(context.Request.URL.Path, "/admin")
	context.File(file)
}

func (handlers Handlers) updateCoupon(context *gin.Context) {
	id := context.Param("id")

	coupon := entities.NewCoupon()

	file, err := context.FormFile("file")
	if err == nil {
		media := entities.NewMediaId(id)

		file.Filename = media.ID
		err := context.SaveUploadedFile(file, "./store/"+media.ID+".jpg")
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Errorf("save upload file failed: %w", err).Error(),
			})

			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("save upload file failed")
			return
		}

		coupon.Media = media
	}

	err = context.ShouldBind(&coupon)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.UpdateCoupon(context, id, coupon)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("update coupon failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update coupon failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "coupon is successfully updated",
	})
}

func (handlers Handlers) deleteCoupon(context *gin.Context) {
	id := context.Param("id")
	err := handlers.svc.DeleteCoupon(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrSubscriptionDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("delete subscription failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete coupon failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "coupon is successfully deleted",
	})
}
