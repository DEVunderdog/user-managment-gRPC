package database

import "context"

type CreateSessionTxParams struct {
	CreateSessionParams
	AfterCreate func(session Session) error
}

type CreateSessionTxResult struct {
	Session Session
}

func (store *SQLStore) CreateSessionTx(ctx context.Context, arg CreateSessionTxParams) (CreateSessionTxResult, error) {
	var result CreateSessionTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Session, err = q.CreateSession(ctx, arg.CreateSessionParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Session)
	})

	return result, err
}