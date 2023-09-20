package handlers

import (
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (handlers Handlers) createOrganization(context *gin.Context) {
	organization := entities.NewOrganization()
	err := context.ShouldBindJSON(&organization)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}
	//TODO : handlers.OrganizationService.CreateOrganization(organization)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "user is successfully created",
	})
}

func (handlers Handlers) getOrganizations(context *gin.Context) {

}

func (handlers Handlers) deleteOrganization(context *gin.Context) {

}

func (handlers Handlers) addOrganizationMembers(context *gin.Context) {

}

func (handlers Handlers) deleteOrganizationMembers(context *gin.Context) {

}

func (handlers Handlers) getOrganizationMembers(context *gin.Context) {

}
