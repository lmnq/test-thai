package service

import (
	"context"

	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
	"github.com/lmnq/test-thai/internal/repo"
)

type Service struct {
	Item
	Category
	Group
	ItemDetail
}

func New(repo *repo.Repo) *Service {
	return &Service{
		Item:     NewItemService(repo.Item),
		Category: NewCategoryService(repo.Category),
		Group:    NewGroupService(repo.Group),
		ItemDetail: NewItemDetailService(
			repo.ItemDetail, repo.Item, repo.Group, repo.Category,
		),
	}
}

// service interfaces -.
type (
	Item interface {
		Create(ctx context.Context, name string) (int, errs.Error)  // create new item
		Get(ctx context.Context, id int) (*model.Item, errs.Error)  // get item by id
		GetAll(ctx context.Context) ([]*model.Item, errs.Error)     // get all items
		Update(ctx context.Context, id int, name string) errs.Error // update item by id
		Delete(ctx context.Context, id int) errs.Error              // delete item by id
	}

	Category interface {
		Create(ctx context.Context, name string) (int, errs.Error)     // create new category
		Get(ctx context.Context, id int) (*model.Category, errs.Error) // get category by id
		GetAll(ctx context.Context) ([]*model.Category, errs.Error)    // get all categories
		Update(ctx context.Context, id int, name string) errs.Error    // update category by id
		Delete(ctx context.Context, id int) errs.Error                 // delete category by id
	}

	Group interface {
		Create(ctx context.Context, name string) (int, errs.Error)  // create new group
		Get(ctx context.Context, id int) (*model.Group, errs.Error) // get group by id
		GetAll(ctx context.Context) ([]*model.Group, errs.Error)    // get all groups
		Update(ctx context.Context, id int, name string) errs.Error // update group by id
		Delete(ctx context.Context, id int) errs.Error              // delete group by id
	}

	ItemDetail interface {
		Create(ctx context.Context, itemDetail *model.ItemDetail, itemName, groupName, categoryName string) (int, errs.Error)  // create new item (if needed) and new item detail
		Get(ctx context.Context, id int) (*model.ItemDetailView, errs.Error)                                                   // get item detail by id
		GetAllFilter(ctx context.Context, filter *model.ItemDetailFilter) ([]*model.ItemDetailView, errs.Error)                // get item detail list by filter
		Update(ctx context.Context, id int, itemDetail *model.ItemDetail, itemName, groupName, categoryName string) errs.Error // update item detail by id
		Delete(ctx context.Context, id int) errs.Error                                                                         // delete item detail by id
	}
)
