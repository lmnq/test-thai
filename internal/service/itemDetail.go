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
	itemDetail *model.ItemDetail,
	itemName, groupName, categoryName string,
) (int, errs.Error) {
	errMsg := ""
	switch {
	case itemName == "":
		errMsg = "item name is empty"
	case groupName == "":
		errMsg = "group name is empty"
	case categoryName == "":
		errMsg = "category name is empty"
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

	groupID, err := s.groupRepo.GetIDByName(ctx, groupName)
	if err == errs.ErrNotFound {
		return 0, errs.Error{
			Err:     fmt.Errorf("get group id by name error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: group does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("get group id by name error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}
	itemDetail.GroupID = groupID

	categoryID, err := s.categoryRepo.GetIDByName(ctx, categoryName)
	if err == errs.ErrNotFound {
		return 0, errs.Error{
			Err:     fmt.Errorf("get category id by name error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: category does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("get category id by name error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}
	itemDetail.CategoryID = categoryID

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

func (s *ItemDetailService) Update(ctx context.Context,
	id int, itemDetail *model.ItemDetail,
	itemName, groupName, categoryName string,
) errs.Error {
	errMsg := ""
	switch {
	case itemName == "":
		errMsg = "item name is empty"
	case groupName == "":
		errMsg = "group name is empty"
	case categoryName == "":
		errMsg = "category name is empty"
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

	itemID, err := s.itemRepo.GetIDByName(ctx, itemName)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("get item id by name error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: item does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("get item id by name error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}
	itemDetail.ItemID = itemID

	groupID, err := s.groupRepo.GetIDByName(ctx, groupName)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("get group id by name error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: group does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("get group id by name error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}
	itemDetail.GroupID = groupID

	categoryID, err := s.categoryRepo.GetIDByName(ctx, categoryName)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("get category id by name error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: category does not exist", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("get category id by name error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}
	itemDetail.CategoryID = categoryID

	err = s.repo.Update(ctx, id, itemDetail)
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
