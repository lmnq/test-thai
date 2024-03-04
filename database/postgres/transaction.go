package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Implementation of sql transaction
// not used.


type Connection interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type ctxTxKey struct{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey{}, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(ctxTxKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}

func (p *Postgres) GetConn(ctx context.Context) Connection {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}

	return p.Pool
}

func (p *Postgres) BeginTransaction(ctx context.Context) (context.Context, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return injectTx(ctx, tx), nil
}

func (p *Postgres) CommitTransaction(ctx context.Context) error {
	tx := extractTx(ctx)
	if tx == nil {
		return fmt.Errorf("transaction missing from context")
	}

	defer func() {
		if tx.Conn() != nil {
			tx.Conn().Close(ctx) // Close the connection regardless of commit success
		}
	}()

	err := tx.Commit(ctx)
	if err != nil {
		p.rollback(ctx, tx, err) // Attempt a rollback
		return err               // Return the original error for the caller to handle
	}
	return nil
}

func (p *Postgres) RollbackTransaction(ctx context.Context) error {
	tx := extractTx(ctx)
	if tx == nil {
		return nil // No transaction to roll back
	}

	defer func() {
		if tx.Conn() != nil {
			tx.Conn().Close(ctx)
		}
	}()

	err := tx.Rollback(ctx)
	if err != nil && err != pgx.ErrTxClosed {
		return fmt.Errorf("transaction rollback failed: %w", err)
	}
	return nil
}

func (p *Postgres) rollback(ctx context.Context, tx pgx.Tx, originalErr error) {
	rollbackErr := tx.Rollback(ctx)
	if rollbackErr != nil && rollbackErr != pgx.ErrTxClosed {
		log.Printf("transaction commit failed: %v; rollback also failed: %v", originalErr, rollbackErr)
	}
}
