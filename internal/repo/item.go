package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lmnq/test-thai/database/postgres"
	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
)

type ItemRepo struct {
	*postgres.Postgres
}

func NewItemRepo(pg *postgres.Postgres) *ItemRepo {
	return &ItemRepo{pg}
}

func (r *ItemRepo) Create(ctx context.Context, name string) (int, error) {
	var res int
	q := "INSERT INTO tbl_items (item_name) VALUES ($1) RETURNING id"
	err := r.Pool.QueryRow(ctx, q, name).Scan(&res)
	if isUniqueConstraintError(err) {
		return 0, errs.ErrUniqueConstraint
	}
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *ItemRepo) Get(ctx context.Context, id int) (*model.Item, error) {
	var item model.Item
	q := `SELECT 
			id,
			item_name,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_items 
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.Pool.QueryRow(ctx, q, id).Scan(
		&item.ID,
		&item.ItemName,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepo) GetIDByName(ctx context.Context, name string) (int, error) {
	var itemID int
	q := `SELECT id FROM tbl_items WHERE item_name = $1 AND deleted_at IS NULL`
	err := r.Pool.QueryRow(ctx, q, name).Scan(&itemID)
	if err == pgx.ErrNoRows {
		return 0, errs.ErrNotFound
	}
	if err != nil {
		return 0, err
	}

	return itemID, nil
}

func (r *ItemRepo) GetAll(ctx context.Context) ([]*model.Item, error) {
	var items []*model.Item
	q := `SELECT 
			id,
			item_name,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_items
		WHERE deleted_at IS NULL
	`
	rows, err := r.Pool.Query(ctx, q)
	if err == pgx.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.ItemName,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *ItemRepo) Update(ctx context.Context, id int, name string) error {
	// updated_at can be automatically updated by postgres trigger function
	q := `UPDATE tbl_items 
		SET item_name = $1,
		updated_at = now()
		WHERE id = $2
	`
	result, err := r.Pool.Exec(ctx, q, name, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *ItemRepo) Delete(ctx context.Context, id int) error {
	q := `UPDATE tbl_items 
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
