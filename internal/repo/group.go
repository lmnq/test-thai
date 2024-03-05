package repo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lmnq/test-thai/database/postgres"
	"github.com/lmnq/test-thai/internal/errs"
	"github.com/lmnq/test-thai/internal/model"
)

type GroupRepo struct {
	*postgres.Postgres
}

func NewGroupRepo(pg *postgres.Postgres) *GroupRepo {
	return &GroupRepo{pg}
}

func (r *GroupRepo) Create(ctx context.Context, name string) (int, error) {
	var res int
	q := "INSERT INTO tbl_groups (group_name) VALUES ($1) RETURNING id"
	err := r.Pool.QueryRow(ctx, q, name).Scan(&res)
	if isUniqueConstraintError(err) {
		return 0, errs.ErrUniqueConstraint
	}
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (r *GroupRepo) Get(ctx context.Context, id int) (*model.Group, error) {
	var group model.Group
	q := `SELECT 
			id,
			group_name,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_groups 
		WHERE id = $1
		AND deleted_at IS NULL
	`
	err := r.Pool.QueryRow(ctx, q, id).Scan(
		&group.ID,
		&group.GroupName,
		&group.CreatedAt,
		&group.UpdatedAt,
		&group.DeletedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *GroupRepo) Exists(ctx context.Context, id int) (bool, error) {
	var result bool
	q := `SELECT EXISTS (SELECT 1 FROM tbl_groups WHERE id = $1 AND deleted_at IS NULL)`
	err := r.Pool.QueryRow(ctx, q, id).Scan(&result)
	if err != nil {
		return false, err
	}

	return result, err
}

func (r *GroupRepo) GetAll(ctx context.Context) ([]*model.Group, error) {
	var groups []*model.Group
	q := `SELECT 
			id,
			group_name,
			created_at,
			updated_at,
			deleted_at
		FROM tbl_groups
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
		var group model.Group
		err := rows.Scan(
			&group.ID,
			&group.GroupName,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		groups = append(groups, &group)
	}

	return groups, nil
}

func (r *GroupRepo) Update(ctx context.Context, id int, name string) error {
	// updated_at can be automatically updated by postgres trigger function
	q := `UPDATE tbl_groups 
		SET group_name = $1,
		updated_at = now()
		WHERE id = $2
		AND deleted_at IS NULL
	`
	result, err := r.Pool.Exec(ctx, q, name, id)
	if isUniqueConstraintError(err) {
		return errs.ErrUniqueConstraint
	}
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *GroupRepo) Delete(ctx context.Context, id int) error {
	q := `UPDATE tbl_groups 
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
