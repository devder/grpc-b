package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// contains the input parameters of the transfer transaction
type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

// result of the transfer tx
type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// performs a money transfer from one account to another
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Valid: true,
				Bool:  true,
			}})

		return err
	})

	return result, err
}
