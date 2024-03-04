package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lmnq/test-thai/database/postgres"
	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
)

type ItemDetailRepo struct {
	*postgres.Postgres
}

func NewItemDetailRepo(pg *postgres.Postgres) *ItemDetailRepo {
	return &ItemDetailRepo{pg}
}

func (r *ItemDetailRepo) Create(ctx context.Context, itemDetail *model.ItemDetail, itemName string) (int, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	defer func() {
		if tx.Conn() != nil {
			tx.Conn().Close(ctx)
		}
	}()

	// check if item name already exists, create item if not
	var itemID int
	q := `INSERT INTO tbl_items (item_name)
		VALUES ($1)
		ON CONFLICT (item_name)
		WHERE deleted_at IS NULL
		DO NOTHING
		RETURNING id
	`
	err = tx.QueryRow(ctx, q, itemName).Scan(&itemID)
	if err != nil {
		return 0, err
	}
	itemDetail.ItemID = itemID

	// create item detail
	var res int
	q = `INSERT INTO tbl_item_details
		(item_id, category_id, group_id, cost, price, sort)
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`
	err = tx.QueryRow(ctx, q,
		itemID,
		itemDetail.CategoryID,
		itemDetail.GroupID,
		itemDetail.Cost,
		itemDetail.Price,
		itemDetail.Sort,
	).Scan(&res)
	if isUniqueConstraintError(err) {
		return 0, errs.ErrUniqueConstraint
	}
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return res, nil
}

func (r *ItemDetailRepo) Get(ctx context.Context, id int) (*model.ItemDetailView, error) {
	var itemDetailView model.ItemDetailView
	// join tbl_items and tbl_categories and tbl_groups to get itemDetailView
	q := `SELECT 
			id,
			item_id,
			i.item_name,
			category_id,
			c.category_name,
			group_id,
			g.group_name,
			cost,
			price,
			sort,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_item_details
		JOIN tbl_items AS i ON tbl_item_details.item_id = tbl_items.id
		JOIN tbl_categories AS c ON tbl_item_details.category_id = tbl_categories.id
		JOIN tbl_groups AS g ON tbl_item_details.group_id = tbl_groups.id
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.Pool.QueryRow(ctx, q, id).Scan(
		&itemDetailView.ID,
		&itemDetailView.ItemID,
		&itemDetailView.ItemName,
		&itemDetailView.CategoryID,
		&itemDetailView.CategoryName,
		&itemDetailView.GroupID,
		&itemDetailView.GroupName,
		&itemDetailView.Cost,
		&itemDetailView.Price,
		&itemDetailView.Sort,
		&itemDetailView.CreatedAt,
		&itemDetailView.UpdatedAt,
		&itemDetailView.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &itemDetailView, nil
}

func (r *ItemDetailRepo) GetAllFilter(ctx context.Context, filter *model.ItemDetailFilter) ([]*model.ItemDetailView, error) {
	var itemDetailViews []*model.ItemDetailView
	// join tbl_items and tbl_categories and tbl_groups to get itemDetailView
	q := `SELECT 
			id,
			item_id,
			i.item_name,
			category_id,
			c.category_name,
			group_id,
			g.group_name,
			cost,
			price,
			sort,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_item_details
		JOIN tbl_items AS i ON tbl_item_details.item_id = tbl_items.id
		JOIN tbl_categories AS c ON tbl_item_details.category_id = tbl_categories.id
		JOIN tbl_groups AS g ON tbl_item_details.group_id = tbl_groups.id
		WHERE 1=1 AND deleted_at IS NULL
	`
	var queryParams []interface{}
	if filter.ID != nil {
		q += " AND id = $1"
		queryParams = append(queryParams, filter.ID)
	}
	if filter.ItemName != nil {
		q += " AND i.item_name = $2"
		queryParams = append(queryParams, filter.ItemName)
	}
	if filter.CategoryName != nil {
		q += " AND c.category_name = $3"
		queryParams = append(queryParams, filter.CategoryName)
	}
	if filter.GroupName != nil {
		q += " AND g.group_name = $4"
		queryParams = append(queryParams, filter.GroupName)
	}

	rows, err := r.Pool.Query(ctx, q, queryParams...)
	if err == pgx.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var itemDetailView model.ItemDetailView
		err := rows.Scan(
			&itemDetailView.ID,
			&itemDetailView.ItemID,
			&itemDetailView.ItemName,
			&itemDetailView.CategoryID,
			&itemDetailView.CategoryName,
			&itemDetailView.GroupID,
			&itemDetailView.GroupName,
			&itemDetailView.Cost,
			&itemDetailView.Price,
			&itemDetailView.Sort,
			&itemDetailView.CreatedAt,
			&itemDetailView.UpdatedAt,
			&itemDetailView.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		itemDetailViews = append(itemDetailViews, &itemDetailView)
	}

	return itemDetailViews, nil
}

func (r *ItemDetailRepo) Update(ctx context.Context, id int, itemDetail *model.ItemDetail) error {
	// updated_at will be automatically updated by postgres trigger function
	q := `UPDATE tbl_item_details 
		SET 
			item_id = $1,
			category_id = $2,
			group_id = $3,
			cost = $4,
			price = $5,
			sort = $6,
			updated_at = now()
		WHERE id = $7
		AND deleted_at IS NULL
	`
	result, err := r.Pool.Exec(ctx, q,
		itemDetail.ItemID,
		itemDetail.CategoryID,
		itemDetail.GroupID,
		itemDetail.Cost,
		itemDetail.Price,
		itemDetail.Sort,
		id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *ItemDetailRepo) Delete(ctx context.Context, id int) error {
	q := `UPDATE tbl_item_details 
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`
	result, err := r.Pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}
