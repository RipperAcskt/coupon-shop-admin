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

type CouponService interface {
	CreateCoupon(ctx context.Context, coupon entities.Coupon) error
}

func (handlers Handlers) createCoupon(context *gin.Context) {
	media := entities.NewMedia()

	file, _ := context.FormFile("file")
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

	coupon := entities.NewCoupon()
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
