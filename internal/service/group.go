package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
	"github.com/lmnq/test-thai/internal/repo"
)

type GroupService struct {
	repo repo.Group
}

func NewGroupService(repo repo.Group) *GroupService {
	return &GroupService{repo}
}

func (s *GroupService) Create(ctx context.Context, name string) (int, errs.Error) {
	if name == "" {
		return 0, errs.Error{
			Err:     errors.New("group name is empty"),
			Code:    400,
			Message: fmt.Sprintf("%s: group name is empty", errs.StatusBadRequestMessage),
		}
	}

	id, err := s.repo.Create(ctx, name)
	if err == errs.ErrUniqueConstraint {
		return 0, errs.Error{
			Err:     fmt.Errorf("create group error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: group name already exists", errs.StatusBadRequestMessage),
		}
	}
	if err != nil {
		return 0, errs.Error{
			Err:     fmt.Errorf("create group error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return id, errs.NilError()
}

func (s *GroupService) Get(ctx context.Context, id int) (*model.Group, errs.Error) {
	group, err := s.repo.Get(ctx, id)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get group error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get group error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return group, errs.NilError()
}

func (s *GroupService) GetAll(ctx context.Context) ([]*model.Group, errs.Error) {
	groups, err := s.repo.GetAll(ctx)
	if err == errs.ErrNotFound {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all groups error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return nil, errs.Error{
			Err:     fmt.Errorf("get all groups error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return groups, errs.NilError()
}

func (s *GroupService) Update(ctx context.Context, id int, name string) errs.Error {
	if name == "" {
		return errs.Error{
			Err:     errors.New("group name is empty"),
			Code:    400,
			Message: fmt.Sprintf("%s: group name is empty", errs.StatusBadRequestMessage),
		}
	}

	err := s.repo.Update(ctx, id, name)
	if err == errs.ErrUniqueConstraint {
		return errs.Error{
			Err:     fmt.Errorf("update group error: %w", err),
			Code:    400,
			Message: fmt.Sprintf("%s: group name already exists", errs.StatusBadRequestMessage),
		}
	}
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("update group error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("update group error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return errs.NilError()
}

func (s *GroupService) Delete(ctx context.Context, id int) errs.Error {
	err := s.repo.Delete(ctx, id)
	if err == errs.ErrNotFound {
		return errs.Error{
			Err:     fmt.Errorf("delete group error: %w", err),
			Code:    404,
			Message: errs.StatusNotFoundMessage,
		}
	}
	if err != nil {
		return errs.Error{
			Err:     fmt.Errorf("delete group error: %w", err),
			Code:    500,
			Message: errs.StatusInternalServerErrorMessage,
		}
	}

	return errs.NilError()
}
