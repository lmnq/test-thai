package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
	"github.com/lmnq/test-thai/internal/repo"
)

type ItemService struct {
	repo repo.Item
}

func NewItemService(repo repo.Item) *ItemService {
	return &ItemService{repo}
}

func (s *ItemService) Create(ctx context.Context, name string) (int, errs.Error) {
	if name == "" {
		return 0, errs.Error{
			Err:     errors.New("item name is empty"),
			Code:    400,
			Message: fmt.Sprintf("%s: item name is empty", errs.StatusBadRequestMessage),
		}
	}

	id, err := s.repo.Create(ctx, name)
	if err == errs.ErrUniqueConstraint {
		return 0, errs.Error{
			Err:     fmt.Errorf("create item error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: item name already exists", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("create item error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return id, errs.NilError()
}

func (s *ItemService) Get(ctx context.Context, id int) (*model.Item, errs.Error) {
	item, err := s.repo.Get(ctx, id)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get item error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get item error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return item, errs.NilError()
}

func (s *ItemService) GetAll(ctx context.Context) ([]*model.Item, errs.Error) {
	items, err := s.repo.GetAll(ctx)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all items error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all items error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return items, errs.NilError()
}

func (s *ItemService) Update(ctx context.Context, id int, name string) errs.Error {
	if name == "" {
		return errs.Error{
			Err:     errors.New("item name is empty"),
			Code:    400,
			Message: fmt.Sprintf("%s: item name is empty", errs.StatusBadRequestMessage),
		}
	}

	err := s.repo.Update(ctx, id, name)
	if err == errs.ErrUniqueConstraint {
		return errs.Error{
			Err:     fmt.Errorf("update item error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: item name already exists", errs.StatusBadRequestMessage),
		}
	}
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("update item error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("update item error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return errs.NilError()
}

func (s *ItemService) Delete(ctx context.Context, id int) errs.Error {
	err := s.repo.Delete(ctx, id)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("delete item error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("delete item error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return errs.NilError()
}
