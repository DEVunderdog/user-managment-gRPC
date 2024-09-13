// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: jwtkeys.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createJWTKey = `-- name: CreateJWTKey :one
insert into jwtkeys (
    public_key,
    private_key,
    algorithm,
    is_active,
    expires_at
) values (
    $1, $2, $3, $4, $5
) returning id, public_key, private_key, algorithm, is_active, expires_at, created_at, updated_at
`

type CreateJWTKeyParams struct {
	PublicKey  string             `json:"public_key"`
	PrivateKey string             `json:"private_key"`
	Algorithm  string             `json:"algorithm"`
	IsActive   pgtype.Bool        `json:"is_active"`
	ExpiresAt  pgtype.Timestamptz `json:"expires_at"`
}

func (q *Queries) CreateJWTKey(ctx context.Context, arg CreateJWTKeyParams) (Jwtkey, error) {
	row := q.db.QueryRow(ctx, createJWTKey,
		arg.PublicKey,
		arg.PrivateKey,
		arg.Algorithm,
		arg.IsActive,
		arg.ExpiresAt,
	)
	var i Jwtkey
	err := row.Scan(
		&i.ID,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Algorithm,
		&i.IsActive,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getActiveKey = `-- name: GetActiveKey :one
select id, public_key, private_key, algorithm, is_active, expires_at, created_at, updated_at from jwtkeys
where is_active = $1
`

func (q *Queries) GetActiveKey(ctx context.Context, isActive pgtype.Bool) (Jwtkey, error) {
	row := q.db.QueryRow(ctx, getActiveKey, isActive)
	var i Jwtkey
	err := row.Scan(
		&i.ID,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Algorithm,
		&i.IsActive,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateJWTKeys = `-- name: UpdateJWTKeys :one
update jwtkeys
set is_active = $1
where id = $2
returning id, public_key, private_key, algorithm, is_active, expires_at, created_at, updated_at
`

type UpdateJWTKeysParams struct {
	IsActive pgtype.Bool `json:"is_active"`
	ID       int64       `json:"id"`
}

func (q *Queries) UpdateJWTKeys(ctx context.Context, arg UpdateJWTKeysParams) (Jwtkey, error) {
	row := q.db.QueryRow(ctx, updateJWTKeys, arg.IsActive, arg.ID)
	var i Jwtkey
	err := row.Scan(
		&i.ID,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Algorithm,
		&i.IsActive,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
