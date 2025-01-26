package db

import (
	"context"
	"fmt"
)

// executes a function within a DB transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	// calls BeginTx under the hood
	tx, err := store.connPool.Begin(ctx)

	if err != nil {
		return err
	}

	query := New(tx)
	err = fn(query)
	if err != nil {
		if rollbackError := tx.Rollback(ctx); rollbackError != nil {
			// combine the two errors if the rollback also returns an error
			return fmt.Errorf("tx error: %w, rollback error: %w", err, rollbackError)
		}
		return err
	}

	return tx.Commit(ctx)
}
