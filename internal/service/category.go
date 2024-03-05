package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
	"github.com/lmnq/test-thai/internal/repo"
)

type CategoryService struct {
	repo repo.Category
}

func NewCategoryService(repo repo.Category) *CategoryService {
	return &CategoryService{repo}
}

func (s *CategoryService) Create(ctx context.Context, name string) (int, errs.Error) {
	if name == "" {
		return 0, errs.Error{
			Err:     errors.New("category name is empty"),
			Code:    400,
			Message: fmt.Sprintf("%s: category name is empty", errs.StatusBadRequestMessage),
		}
	}

	id, err := s.repo.Create(ctx, name)
	if err == errs.ErrUniqueConstraint {
		return 0, errs.Error{
			Err:     fmt.Errorf("create category error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: category name already exists", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("create category error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return id, errs.NilError()
}

func (s *CategoryService) Get(ctx context.Context, id int) (*model.Category, errs.Error) {
	category, err := s.repo.Get(ctx, id)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get category error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get category error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return category, errs.NilError()
}

func (s *CategoryService) GetAll(ctx context.Context) ([]*model.Category, errs.Error) {
	categories, err := s.repo.GetAll(ctx)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all categories error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all categories error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return categories, errs.NilError()
}

func (s *CategoryService) Update(ctx context.Context, id int, name string) errs.Error {
	if name == "" {
		return errs.Error{
			Err:     errors.New("category name is empty"),
			Code:    400,
			Message: fmt.Sprintf("%s: category name is empty", errs.StatusBadRequestMessage),
		}
	}

	err := s.repo.Update(ctx, id, name)
	if err == errs.ErrUniqueConstraint {
		return errs.Error{
			Err:     fmt.Errorf("update category error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: category name already exists", errs.StatusBadRequestMessage),
		}
	}
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("update category error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("update category error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return errs.NilError()
}

func (s *CategoryService) Delete(ctx context.Context, id int) errs.Error {
	err := s.repo.Delete(ctx, id)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("delete category error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("delete category error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return errs.NilError()
}
