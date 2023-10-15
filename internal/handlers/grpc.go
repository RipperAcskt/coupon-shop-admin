package handlers

import (
	"context"
	"fmt"
	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
	"github.com/RipperAcskt/coupon-shop-admin/internal/service"
	adminpb "github.com/RipperAcskt/coupon-shop-admin/proto"
)

type Server struct {
	service.CouponService
	service.SubscriptionService
	adminpb.UnimplementedAdminServiceServer
}

type AdminService interface {
	GetCoupons(ctx context.Context) ([]entities.Coupon, error)
	GetSubscriptions(ctx context.Context) ([]entities.Subscription, error)
}

func (s Server) GetSubsGRPC(ctx context.Context, in *adminpb.Empty) (*adminpb.SubscriptionsResponse, error) {
	subs, err := s.SubscriptionService.GetSubscriptions(ctx)
	fmt.Printf("%+v", subs)
	if err != nil {
		return nil, err
	}

	var Response = &adminpb.SubscriptionsResponse{Subs: make([]*adminpb.Subscription, len(subs))}

	for i := range subs {
		var sub = &adminpb.Subscription{
			Name:        subs[i].Name,
			Description: subs[i].Description,
			Price:       int32(subs[i].Price),
			Level:       int32(subs[i].Level),
		}
		Response.Subs[i] = sub
	}

	return Response, nil
}

func (s Server) GetCouponsGRPC(ctx context.Context, in *adminpb.Empty) (*adminpb.GetCouponsResponse, error) {
	coupons, err := s.CouponService.GetCoupons(ctx)
	if err != nil {
		return nil, err
	}
	var Response = &adminpb.GetCouponsResponse{Coupons: make([]*adminpb.Coupon, len(coupons))}
	for i, v := range coupons {
		var media = &adminpb.Media{
			ID:   v.Media.ID,
			Path: v.Media.Path,
		}

		var coupon = &adminpb.Coupon{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Price:       int32(v.Price),
			Level:       int32(v.Level),
			Percent:     int32(v.Percent),
			ContentUrl:  v.ContentUrl,
			Media:       media,
		}

		Response.Coupons[i] = coupon
	}
	return Response, nil
}
