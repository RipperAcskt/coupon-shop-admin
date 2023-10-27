package main

import (
	"context"
	"github.com/RipperAcskt/coupon-shop-admin/internal/handlers"
	"github.com/RipperAcskt/coupon-shop-admin/internal/service"
	adminpb "github.com/RipperAcskt/coupon-shop-admin/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartGrpcServer(ctx context.Context, subService service.SubscriptionService, couponService service.CouponService, organizationService service.OrganizationService, membersService service.MembersService) error {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalln("Failed to listen grpc server: ", err)
	}
	server := handlers.Server{
		MembersService:                  membersService,
		CouponService:                   couponService,
		SubscriptionService:             subService,
		OrganizationService:             organizationService,
		UnimplementedAdminServiceServer: adminpb.UnimplementedAdminServiceServer{},
	}

	s := grpc.NewServer()
	adminpb.RegisterAdminServiceServer(s, server)
	log.Println("serving gRPC on localhost:8080")

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		errCh <- s.Serve(lis)
	}()
	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		s.GracefulStop()
	}
	return nil
}
