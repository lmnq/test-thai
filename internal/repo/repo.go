package repo

import (
	"context"

	"github.com/lmnq/test-thai/database/postgres"
	"github.com/lmnq/test-thai/internal/model"
)

type Repo struct {
	Item
	Category
	Group
	ItemDetail
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		Item:       NewItemRepo(pg),
		Category:   NewCategoryRepo(pg),
		Group:      NewGroupRepo(pg),
		ItemDetail: NewItemDetailRepo(pg),
	}
}

// repo interfaces -.
type (
	Item interface {
		Create(ctx context.Context, name string) (int, error)      // create new item
		Get(ctx context.Context, id int) (*model.Item, error)      // get item by id
		GetIDByName(ctx context.Context, name string) (int, error) // get item id by name
		GetAll(ctx context.Context) ([]*model.Item, error)         // get all items
		Update(ctx context.Context, id int, name string) error     // update item by id
		Delete(ctx context.Context, id int) error                  // delete item by id
	}

	Category interface {
		Create(ctx context.Context, name string) (int, error)     // create new category
		Get(ctx context.Context, id int) (*model.Category, error) // get category by id
		Exists(ctx context.Context, id int) (bool, error)         // check if category exists
		GetAll(ctx context.Context) ([]*model.Category, error)    // get all categories
		Update(ctx context.Context, id int, name string) error    // update category by id
		Delete(ctx context.Context, id int) error                 // delete category by id
	}

	Group interface {
		Create(ctx context.Context, name string) (int, error)  // create new group
		Get(ctx context.Context, id int) (*model.Group, error) // get group by id
		Exists(ctx context.Context, id int) (bool, error)      // check if group exists
		GetAll(ctx context.Context) ([]*model.Group, error)    // get all groups
		Update(ctx context.Context, id int, name string) error // update group by id
		Delete(ctx context.Context, id int) error              // delete group by id
	}

	ItemDetail interface {
		Create(ctx context.Context, itemDetail *model.ItemDetail, itemName string) (int, error)            // create new item (if needed) and new item detail
		Get(ctx context.Context, id int) (*model.ItemDetailView, error)                                    // get item detail by id
		Exists(ctx context.Context, id int) (bool, error)                                                  // check if item detail exists
		GetAllFilter(ctx context.Context, filter *model.ItemDetailFilter) ([]*model.ItemDetailView, error) // get item detail list by filter
		Update(ctx context.Context, id int, itemName string, itemDetail *model.ItemDetail) error           // update item detail by id
		Delete(ctx context.Context, id int) error                                                          // delete item detail by id
	}
)
