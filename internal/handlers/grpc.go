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
	service.OrganizationService
	service.MembersService
	service.CategoryService
	service.RegionService
	service.LinkService
	adminpb.UnimplementedAdminServiceServer
}

type AdminService interface {
	GetCoupons(ctx context.Context) ([]entities.Coupon, error)
	GetSubscriptions(ctx context.Context) ([]entities.Subscription, error)
}

func (s Server) GetLinksGRPC(ctx context.Context, in *adminpb.Region) (*adminpb.Links, error) {
	links, err := s.LinkService.GetLinkByRegion(ctx, in.Region)
	if err != nil {
		return nil, err
	}

	var Response = &adminpb.Links{}

	var link = &adminpb.Link{
		Id:     links.Id,
		Name:   links.Name,
		Link:   links.Link,
		Region: links.Region,
	}
	Response.Links = link

	return Response, nil
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

func (s Server) UpdateCoupon(ctx context.Context, in *adminpb.Coupon) (*adminpb.UpdateMembersResponse, error) {
	err := s.CouponService.UpdateCoupon(ctx, in.ID, entities.Coupon{
		Name:        in.Name,
		Description: in.Description,
		Price:       int(in.Price),
		Percent:     int(in.Percent),
		Region:      in.Region,
		Category:    in.Category,
		Subcategory: &in.Subcategory,
	})
	if err != nil {
		return nil, err
	}

	var Response = &adminpb.UpdateMembersResponse{Message: "updated"}

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
			Region:      v.Region,
			Category:    v.Category,
		}
		if v.Subcategory != nil {
			coupon.Subcategory = *v.Subcategory
		}
		Response.Coupons[i] = coupon
	}
	return Response, nil
}

func (s Server) GetCouponsSearchGRPC(ctx context.Context, in *adminpb.Search) (*adminpb.GetCouponsResponse, error) {
	coupons, err := s.CouponService.GetCouponsSearch(ctx, in.S)
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
			Region:      v.Region,
			Category:    v.Category,
		}
		if v.Subcategory != nil {
			coupon.Subcategory = *v.Subcategory
		}
		Response.Coupons[i] = coupon
	}
	return Response, nil
}

func (s Server) GetCategoriesGRPC(ctx context.Context, in *adminpb.Empty) (*adminpb.GetCategoryResponse, error) {
	categories, err := s.CategoryService.GetCategories(ctx, false)
	if err != nil {
		return nil, err
	}

	subcategories, err := s.CategoryService.GetCategories(ctx, true)
	if err != nil {
		return nil, err
	}

	var Response = &adminpb.GetCategoryResponse{}
	for _, cat := range categories {
		var categoriesResp adminpb.CategoryResponse
		categoriesResp.ID = cat.Id
		categoriesResp.Name = cat.Name
		for _, sub := range subcategories {
			if sub.Subcategory == cat.Name {
				categoriesResp.Subcategories = append(categoriesResp.Subcategories, &adminpb.SubcategoryResponse{
					ID:   sub.Id,
					Name: sub.Name,
				})
			}
		}
		Response.Categories = append(Response.Categories, &categoriesResp)
	}
	return Response, nil
}

func (s Server) GetRegionsGRPC(ctx context.Context, in *adminpb.Empty) (*adminpb.RegionResponse, error) {
	regions, err := s.RegionService.GetRegions(ctx)
	if err != nil {
		return nil, err
	}

	var Response = &adminpb.RegionResponse{}
	for _, v := range regions {
		Response.Regions = append(Response.Regions, &adminpb.Region{
			Region: v.Name,
			Tg:     v.Tg,
			Vk:     v.Vk,
		})
	}
	return Response, nil
}

func (s Server) GetCouponsByRegionGRPC(ctx context.Context, region *adminpb.Region) (*adminpb.GetCouponsResponse, error) {
	coupons, err := s.CouponService.GetCouponsByRegion(ctx, region.Region)
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
			Region:      v.Region,
			Category:    v.Category,
		}
		if v.Subcategory != nil {
			coupon.Subcategory = *v.Subcategory
		}

		Response.Coupons[i] = coupon
	}
	return Response, nil
}

func (s Server) GetCouponsByCategoryGRPC(ctx context.Context, category *adminpb.Category) (*adminpb.GetCouponsResponse, error) {
	var coupons []entities.Coupon
	var err error
	if category.Subcategory {
		coupons, err = s.CouponService.GetCouponsBySubcategory(ctx, category.Name)
		if err != nil {
			return nil, err
		}
	} else {
		coupons, err = s.CouponService.GetCouponsByCategory(ctx, category.Name)
		if err != nil {
			return nil, err
		}
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
			Region:      v.Region,
			Category:    v.Category,
		}
		if v.Subcategory != nil {
			coupon.Subcategory = *v.Subcategory
		}

		Response.Coupons[i] = coupon
	}
	return Response, nil
}

func (s Server) GetOrganizationInfo(ctx context.Context, in *adminpb.InfoOrganizationRequest) (*adminpb.InfoOrganizationResponse, error) {
	orgInfo, err := s.OrganizationService.GetOrganization(ctx, in.GetOrgId())
	if err != nil {
		return nil, err
	}
	var Response = &adminpb.InfoOrganizationResponse{
		ID:                orgInfo.ID,
		Name:              orgInfo.Name,
		EmailAdmin:        orgInfo.EmailAdmin,
		Orgn:              orgInfo.ORGN,
		Kpp:               orgInfo.KPP,
		Inn:               orgInfo.INN,
		LevelSubscription: int32(orgInfo.LevelSubscription),
		Address:           orgInfo.Address,
		Members:           make([]*adminpb.MemberInfo, len(orgInfo.Members)),
		ContentUrl:        orgInfo.ContentUrl,
	}
	fmt.Printf("%+v", Response)
	for i, v := range orgInfo.Members {
		Response.Members[i] = &adminpb.MemberInfo{
			Id:         v.ID,
			Email:      v.Email,
			FirstName:  v.FirstName,
			SecondName: v.SecondName,
			OrgID:      v.OrganizationID,
			Role:       string(v.Role),
		}
	}

	return Response, nil
}

func (s Server) UpdateOrganizationInfo(ctx context.Context, in *adminpb.UpdateOrganizationRequest) (*adminpb.UpdateOrganizationResponse, error) {
	fmt.Println("entered GRPC")
	fmt.Println(*in)
	if entities.Role(in.GetRoleUser()) != entities.Owner && entities.Role(in.GetRoleUser()) != entities.Editor {
		return &adminpb.UpdateOrganizationResponse{Message: "Action not allowed for default users"}, nil
	}
	fmt.Println("from if not returns")
	err := s.OrganizationService.UpdateOrganization(ctx, entities.Organization{
		Name:              in.GetName(),
		EmailAdmin:        in.EmailAdmin,
		LevelSubscription: int(in.LevelSubscription),
		ORGN:              in.GetOrgn(),
		KPP:               in.GetKpp(),
		INN:               in.GetInn(),
		Address:           in.GetAddress(),
	}, in.GetID())
	fmt.Println("err", err)
	if err != nil {
		return nil, err
	}
	return &adminpb.UpdateOrganizationResponse{Message: "organization successfully updated"}, nil
}

func (s Server) UpdateMembersInfo(ctx context.Context, in *adminpb.UpdateMembersRequest) (*adminpb.UpdateMembersResponse, error) {
	fmt.Println("entered gRPC")
	if entities.Role(in.GetRoleUser()) != entities.Owner && entities.Role(in.GetRoleUser()) != entities.Editor {
		return &adminpb.UpdateMembersResponse{Message: "Action not allowed for default users"}, nil
	}
	members := make([]entities.Member, 0)
	fmt.Println("Members : ", in.Members)
	for i := range in.Members {
		var member entities.Member
		member.Email = in.Members[i].Email
		member.FirstName = in.Members[i].FirstName
		member.SecondName = in.Members[i].SecondName
		member.Role = entities.Role(in.Members[i].Role)
		members = append(members, member)
	}
	fmt.Println("members mapped : ", members)
	fmt.Println("orgID", in.GetOrganizationID())
	err := s.MembersService.UpdateMembers(ctx, members, in.GetOrganizationID())
	fmt.Println("err", err)
	if err != nil {
		return nil, err
	}
	return &adminpb.UpdateMembersResponse{Message: "members successfully updated"}, nil
}
