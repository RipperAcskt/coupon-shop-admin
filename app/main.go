package main

import (
	"errors"
	"net/http"

	"github.com/RipperAcskt/coupon-shop-admin/internal/handlers"
	"github.com/RipperAcskt/coupon-shop-admin/internal/repository"
	"github.com/RipperAcskt/coupon-shop-admin/internal/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("load env failed")
	}

	repo, err := repository.New()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("repository new failed")
	}
	defer func() {
		err := repo.Close()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("repo close failed")
		}
	}()

	svc := service.New(repo)
	handlersEngine, err := handlers.SetRequestHandlers(svc)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("set request handlers failed")
	}

	srv := &Server{}

	go func() {
		if err := srv.Run(handlersEngine); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("server run failed")
			return
		}
	}()

	if err := srv.WaitForShutDown(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("server shut down failed")
	}
}
