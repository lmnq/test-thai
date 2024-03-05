package repo

import (
	"context"
	"fmt"
	"log"

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
		DO UPDATE SET item_name = excluded.item_name
		RETURNING id
	`
	err = tx.QueryRow(ctx, q, itemName).Scan(&itemID)
	if err != nil {
		log.Println("AAAA", err)
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
	// if isUniqueConstraintError(err) {
	// 	return 0, errs.ErrUniqueConstraint
	// }
	if err != nil {
		log.Println("BBBB", err)
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
			itd.id,
			itd.item_id,
			i.item_name,
			itd.category_id,
			c.category_name,
			itd.group_id,
			g.group_name,
			itd.cost,
			itd.price,
			itd.sort,
			itd.created_at,
			itd.updated_at,
			itd.deleted_at
		FROM tbl_item_details AS itd
		JOIN tbl_items AS i ON itd.item_id = i.id
		JOIN tbl_categories AS c ON itd.category_id = c.id
		JOIN tbl_groups AS g ON itd.group_id = g.id
		WHERE itd.id = $1 AND itd.deleted_at IS NULL
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

func (r *ItemDetailRepo) Exists(ctx context.Context, id int) (bool, error) {
	var result bool
	q := `SELECT EXISTS (SELECT 1 FROM tbl_item_details WHERE id = $1 AND deleted_at IS NULL)`
	err := r.Pool.QueryRow(ctx, q, id).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, err
}

func (r *ItemDetailRepo) GetAllFilter(ctx context.Context, filter *model.ItemDetailFilter) ([]*model.ItemDetailView, error) {
	var itemDetailViews []*model.ItemDetailView
	// join tbl_items and tbl_categories and tbl_groups to get itemDetailView
	q := `SELECT 
			itd.id,
			itd.item_id,
			i.item_name,
			itd.category_id,
			c.category_name,
			itd.group_id,
			g.group_name,
			itd.cost,
			itd.price,
			itd.sort,
			itd.created_at,
			itd.updated_at,
			itd.deleted_at
		FROM tbl_item_details AS itd
		JOIN tbl_items AS i ON itd.item_id = i.id
		JOIN tbl_categories AS c ON itd.category_id = c.id
		JOIN tbl_groups AS g ON itd.group_id = g.id
		WHERE 1=1 AND itd.deleted_at IS NULL
	`
	var queryParams []interface{}
	if filter.ID != nil {
		queryParams = append(queryParams, filter.ID)
		q += fmt.Sprintf(" AND itd.id = $%d", len(queryParams))
	}
	if filter.ItemName != nil {
		queryParams = append(queryParams, filter.ItemName)
		q += fmt.Sprintf(" AND i.item_name = $%d", len(queryParams))
	}
	if filter.CategoryName != nil {
		queryParams = append(queryParams, filter.CategoryName)
		q += fmt.Sprintf(" AND c.category_name = $%d", len(queryParams))
	}
	if filter.GroupName != nil {
		queryParams = append(queryParams, filter.GroupName)
		q += fmt.Sprintf(" AND g.group_name = $%d", len(queryParams))
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

func (r *ItemDetailRepo) Update(ctx context.Context, id int, itemName string, itemDetail *model.ItemDetail) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
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

	// get item_id by current item_detail.id, then update items.item_name
	var itemID int
	q := `SELECT item_id FROM tbl_item_details WHERE id = $1 AND deleted_at IS NULL`
	err = tx.QueryRow(ctx, q, id).Scan(&itemID)
	if err != nil {
		return err
	}

	q = `UPDATE tbl_items SET item_name = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, q, itemName, itemID)
	if err != nil {
		return err
	}

	q = `UPDATE tbl_item_details
		SET 
			category_id = $1,
			group_id = $2,
			cost = $3,
			price = $4,
			sort = $5,
			updated_at = now()
		WHERE id = $6
		AND deleted_at IS NULL
	`
	_, err = tx.Exec(ctx, q,
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

	if err := tx.Commit(ctx); err != nil {
		return err
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
