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

type OrganizationService interface {
	CreateOrganization(ctx context.Context, org entities.Organization) error
	GetOrganizations(ctx context.Context) ([]entities.Organization, error)
	GetOrganization(ctx context.Context, organizationID string) (entities.Organization, error)
	DeleteOrganization(ctx context.Context, id string) error
	UpdateOrganization(ctx context.Context, org entities.Organization) error
}

func (handlers Handlers) createOrganization(context *gin.Context) {

	organization := entities.NewOrganization()
	image := entities.NewImageID(organization.ID)

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
	file.Filename = image.ID
	err = context.SaveUploadedFile(file, "./store/organization/"+image.ID+".jpg")
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("save upload file failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("save upload file failed")
		return
	}

	organization.OrgImage = image
	err = context.ShouldBind(&organization)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.CreateOrganization(context, organization)
	if err != nil {
		if errors.Is(err, entities.ErrOrganizationAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("create organization failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("create organization failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "organization is successfully created",
	})
}

func (handlers Handlers) getOrganizations(context *gin.Context) {
	orgs, err := handlers.svc.GetOrganizations(context)
	if err != nil {
		if errors.Is(err, entities.ErrNoAnyOrganization) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get organization failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get organization failed")
		return
	}
	context.JSON(http.StatusCreated, orgs)
}

func (handlers Handlers) deleteOrganization(context *gin.Context) {
	id := context.Param("id")
	err := handlers.svc.DeleteOrganization(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrOrganizationnDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("delete organization failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("delete organization failed")
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "organization is successfully deleted",
	})
}

func (handlers Handlers) getOrganization(context *gin.Context) {
	id := context.Param("id")
	org, err := handlers.svc.GetOrganization(context, id)
	if err != nil {
		if errors.Is(err, entities.ErrOrganizationnDoesNotExist) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("get organization failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get organization failed")
		return
	}
	context.JSON(http.StatusOK, org)
}

func (handlers Handlers) updateOrganization(context *gin.Context) {

	organization := entities.NewOrganization()
	image := entities.NewImageID(organization.ID)

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
	file.Filename = image.ID
	err = context.SaveUploadedFile(file, "./store/organization/"+image.ID+".jpg")
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("save upload file failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("save upload file failed")
		return
	}

	organization.OrgImage = image
	err = context.ShouldBind(&organization)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Errorf("should bind json failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("should bind json failed")
		return
	}

	err = handlers.svc.UpdateOrganization(context, organization)
	if err != nil {
		if errors.Is(err, entities.ErrOrganizationAlreadyExists) {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("update organization failed: %w", err).Error(),
		})

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("update organization failed")
		return
	}
	context.JSON(http.StatusCreated, gin.H{
		"message": "organization is successfully updated",
	})
}
