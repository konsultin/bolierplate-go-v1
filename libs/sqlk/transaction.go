package sqlk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TxContext wraps sqlx.Tx with context
type TxContext struct {
	tx  *sqlx.Tx
	ctx context.Context
}

// Tx returns the underlying sqlx.Tx
func (tc *TxContext) Tx() *sqlx.Tx {
	return tc.tx
}

// Ctx returns the context
func (tc *TxContext) Ctx() context.Context {
	return tc.ctx
}

// Commit commits the transaction
func (tc *TxContext) Commit() error {
	return tc.tx.Commit()
}

// Rollback rolls back the transaction
func (tc *TxContext) Rollback() error {
	return tc.tx.Rollback()
}

// BeginTx starts a new transaction with the given context and options
func BeginTx(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions) (*TxContext, error) {
	tx, err := db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("sqlk: failed to begin transaction: %w", err)
	}

	return &TxContext{
		tx:  tx,
		ctx: ctx,
	}, nil
}

// TxFunc is a function that executes within a transaction
type TxFunc func(*TxContext) error

// WithTransaction executes the given function within a transaction
// Automatically commits on success, rolls back on error or panic
func WithTransaction(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions, fn TxFunc) (err error) {
	// Begin transaction
	txCtx, err := BeginTx(ctx, db, opts)
	if err != nil {
		return err
	}

	// Handle panic and errors
	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			_ = txCtx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			// Rollback on error
			_ = txCtx.Rollback()
		} else {
			// Commit on success
			err = txCtx.Commit()
		}
	}()

	// Execute function
	err = fn(txCtx)
	return err
}

// CommitOrRollback commits the transaction if err is nil, otherwise rolls back
// This is useful for manual transaction control
func CommitOrRollback(txCtx *TxContext, err error) error {
	if err != nil {
		if rbErr := txCtx.Rollback(); rbErr != nil {
			return fmt.Errorf("sqlk: failed to rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if commitErr := txCtx.Commit(); commitErr != nil {
		return fmt.Errorf("sqlk: failed to commit transaction: %w", commitErr)
	}

	return nil
}
