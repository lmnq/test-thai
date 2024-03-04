package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lmnq/test-thai/database/postgres"
	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
)

type CategoryRepo struct {
	*postgres.Postgres
}

func NewCategoryRepo(pg *postgres.Postgres) *CategoryRepo {
	return &CategoryRepo{pg}
}

func (r *CategoryRepo) Create(ctx context.Context, name string) (int, error) {
	var res int
	q := "INSERT INTO tbl_categories (category_name) VALUES ($1) RETURNING id"
	err := r.Pool.QueryRow(ctx, q, name).Scan(&res)
	if isUniqueConstraintError(err) {
		return 0, errs.ErrUniqueConstraint
	}
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *CategoryRepo) Get(ctx context.Context, id int) (*model.Category, error) {
	var category model.Category
	q := `SELECT 
			id,
			category_name,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_categories 
		WHERE id = $1
		AND deleted_at IS NULL
	`
	err := r.Pool.QueryRow(ctx, q, id).Scan(
		&category.ID,
		&category.CategoryName,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepo) GetIDByName(ctx context.Context, name string) (int, error) {
	var categoryID int
	q := `SELECT id FROM tbl_categories WHERE category_name = $1 AND deleted_at IS NULL`
	err := r.Pool.QueryRow(ctx, q, name).Scan(&categoryID)
	if err == pgx.ErrNoRows {
		return 0, errs.ErrNotFound
	}
	if err != nil {
		return 0, err
	}

	return categoryID, nil
}

func (r *CategoryRepo) GetAll(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	q := `SELECT 
			id,
			category_name,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_categories
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
		var category model.Category
		err := rows.Scan(
			&category.ID,
			&category.CategoryName,
			&category.CreatedAt,
			&category.UpdatedAt,
			&category.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryRepo) Update(ctx context.Context, id int, name string) error {
	// updated_at will be automatically updated by postgres trigger function
	q := `UPDATE tbl_categories 
		SET category_name = $1,
		updated_at = now()
		WHERE id = $2
		AND deleted_at IS NULL
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

func (r *CategoryRepo) Delete(ctx context.Context, id int) error {

	q := `UPDATE tbl_categories 
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
