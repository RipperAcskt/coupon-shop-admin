package service

import (
	"context"
	"github.com/google/uuid"

	"github.com/RipperAcskt/coupon-shop-admin/internal/entities"
)

type CategoryService struct {
	repo CategoryRepoInterface
}

type CategoryRepoInterface interface {
	CreateCategory(ctx context.Context, category entities.Category) error
	GetCategories(ctx context.Context) ([]entities.Category, error)
	GetCategory(ctx context.Context, id string) (entities.Category, error)
	UpdateCategory(ctx context.Context, id, name string) error
	DeleteCategory(ctx context.Context, id string) error

	CreateSubcategory(ctx context.Context, category entities.Category) error
	GetSubcategories(ctx context.Context) ([]entities.Category, error)
	GetSubcategory(ctx context.Context, id string) (entities.Category, error)
	UpdateSubcategory(ctx context.Context, id, name string) error
	DeleteSubcategory(ctx context.Context, id string) error
}

func NewCategoryService(repo CategoryRepoInterface) CategoryService {
	return CategoryService{repo: repo}
}

func (svc CategoryService) CreateCategory(ctx context.Context, category entities.Category) error {
	category.Id = uuid.NewString()
	if category.Subcategory != "" {
		err := svc.repo.CreateSubcategory(ctx, category)
		return err
	}
	err := svc.repo.CreateCategory(ctx, category)
	return err
}

func (svc CategoryService) GetCategories(ctx context.Context, subcategory bool) ([]entities.Category, error) {
	if subcategory {
		categories, err := svc.repo.GetSubcategories(ctx)
		return categories, err
	}
	categories, err := svc.repo.GetCategories(ctx)
	return categories, err
}

func (svc CategoryService) GetCategory(ctx context.Context, id string, subcategory bool) (entities.Category, error) {
	if subcategory {
		category, err := svc.repo.GetSubcategory(ctx, id)
		return category, err
	}
	category, err := svc.repo.GetCategory(ctx, id)
	return category, err
}

func (svc CategoryService) UpdateCategory(ctx context.Context, id, name string, subcategory bool) error {
	if subcategory {
		err := svc.repo.UpdateSubcategory(ctx, id, name)
		return err
	}
	err := svc.repo.UpdateCategory(ctx, id, name)
	return err
}

func (svc CategoryService) DeleteCategory(ctx context.Context, id string, subcategory bool) error {
	if subcategory {
		err := svc.repo.DeleteSubcategory(ctx, id)
		return err
	}
	err := svc.repo.DeleteCategory(ctx, id)
	return err
}
