package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
	"github.com/lmnq/test-thai/internal/repo"
)

type ItemDetailService struct {
	repo         repo.ItemDetail
	itemRepo     repo.Item
	groupRepo    repo.Group
	categoryRepo repo.Category
}

func NewItemDetailService(
	repo repo.ItemDetail,
	itemRepo repo.Item,
	groupRepo repo.Group,
	categoryRepo repo.Category,
) *ItemDetailService {
	return &ItemDetailService{
		repo:         repo,
		itemRepo:     itemRepo,
		groupRepo:    groupRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ItemDetailService) Create(ctx context.Context,
	itemDetail *model.ItemDetail, itemName string,
) (int, errs.Error) {
	errMsg := ""
	switch {
	case itemName == "":
		errMsg = "item name is empty"
	case itemDetail.GroupID <= 0:
		errMsg = "invalid group id"
	case itemDetail.CategoryID <= 0:
		errMsg = "invalid category id"
	case itemDetail.Cost <= 0:
		errMsg = "cost must be greater than 0"
	case itemDetail.Price <= 0:
		errMsg = "price must be greater than 0"
	case itemDetail.Sort <= 0:
		errMsg = "sort must be greater than 0"
	}
	if errMsg != "" {
		return 0, errs.Error{
			Err:     errors.New(errMsg),
			Code:    400,
			Message: fmt.Sprintf("%s: %s", errs.StatusBadRequestMessage, errMsg),
		}
	}

	groupIDExists, err := s.groupRepo.Exists(ctx, itemDetail.GroupID)
	if !groupIDExists {
		return 0, errs.Error{
			Err:     fmt.Errorf("group does not exist"),
			Code:    400,
			Message: fmt.Sprintf("%s: group does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("check if group exists error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	categoryIDExists, err := s.categoryRepo.Exists(ctx, itemDetail.CategoryID)
	if !categoryIDExists {
		return 0, errs.Error{
			Err:     fmt.Errorf("category does not exist"),
			Code:    400,
			Message: fmt.Sprintf("%s: category does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("check if category exists error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	// create new item with itemName, if does not exist. otherwise use existing item.
	// then create new item detail
	res, err := s.repo.Create(ctx, itemDetail, itemName)
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("create item detail error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return res, errs.NilError()
}

func (s *ItemDetailService) Get(ctx context.Context, id int) (*model.ItemDetailView, errs.Error) {
	itemDetailView, err := s.repo.Get(ctx, id)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get item detail error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get item detail error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return itemDetailView, errs.NilError()
}

func (s *ItemDetailService) GetAllFilter(ctx context.Context, filter *model.ItemDetailFilter) ([]*model.ItemDetailView, errs.Error) {
	itemDetailViews, err := s.repo.GetAllFilter(ctx, filter)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all item detail error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all item detail error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return itemDetailViews, errs.NilError()
}

func (s *ItemDetailService) Update(ctx context.Context, id int, itemName string, itemDetail *model.ItemDetail) errs.Error {
	errMsg := ""
	switch {
	case itemName == "":
		errMsg = "item name is empty"
	case itemDetail.GroupID <= 0:
		errMsg = "invalid group id"
	case itemDetail.CategoryID == 0:
		errMsg = "invalid category id"
	case itemDetail.Cost <= 0:
		errMsg = "cost must be greater than 0"
	case itemDetail.Price <= 0:
		errMsg = "price must be greater than 0"
	case itemDetail.Sort <= 0:
		errMsg = "sort must be greater than 0"
	}
	if errMsg != "" {
		return errs.Error{
			Err:     errors.New(errMsg),
			Code:    400,
			Message: fmt.Sprintf("%s: %s", errs.StatusBadRequestMessage, errMsg),
		}
	}

	exists, err := s.repo.Exists(ctx, id)
	if !exists {
		return errs.Error{
			Err:     fmt.Errorf("item detail does not exist"),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("check if item detail exists error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	groupIDExists, err := s.groupRepo.Exists(ctx, itemDetail.GroupID)
	if !groupIDExists {
		return errs.Error{
			Err:     fmt.Errorf("group does not exist"),
			Code:    400,
			Message: fmt.Sprintf("%s: group does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("check if group exists error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	categoryIDExists, err := s.categoryRepo.Exists(ctx, itemDetail.CategoryID)
	if !categoryIDExists {
		return errs.Error{
			Err:     fmt.Errorf("category does not exist"),
			Code:    400,
			Message: fmt.Sprintf("%s: category does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("check if category exists error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	err = s.repo.Update(ctx, id, itemName, itemDetail)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("update item detail error: %w", err),
			Code:    404,
			Message: fmt.Sprintf("%s: item detail does not exist", errs.StatusNotFoundMessage),
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("update item detail error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return errs.NilError()
}

func (s *ItemDetailService) Delete(ctx context.Context, id int) errs.Error {
	err := s.repo.Delete(ctx, id)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("delete item detail error: %w", err),
			Code:    404,
			Message: fmt.Sprintf("%s: item detail does not exist", errs.StatusNotFoundMessage),
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("delete item detail error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}
	return errs.NilError()
}
