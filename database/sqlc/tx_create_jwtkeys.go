package database

import "context"

type CreateJWTKeysTxParams struct {
	CreateJWTKeyParams
	AfterCreate func(jwtKey Jwtkey) error
}

type CreateJWTKeyTxResult struct {
	Jwtkey Jwtkey
}

func (store *SQLStore) CreateJWTKeyTx(ctx context.Context, arg CreateJWTKeysTxParams) (CreateJWTKeyTxResult, error) {
	var result CreateJWTKeyTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Jwtkey, err = q.CreateJWTKey(ctx, arg.CreateJWTKeyParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Jwtkey)
	})

	return result, err
}