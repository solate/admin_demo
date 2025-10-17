package ent

import (
	"context"
	"fmt"

	"admin_backend/pkg/ent/generated"

	"github.com/zeromicro/go-zero/core/logx"
)

// WithTx 在事务中执行fn
func WithTx(ctx context.Context, client *Client, fn func(tx *generated.Tx) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			logx.WithContext(ctx).Errorf("事务执行过程中发生panic: %v", v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
