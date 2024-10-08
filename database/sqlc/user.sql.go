// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const checkForExistingUser = `-- name: CheckForExistingUser :one
select exists(select 1 from users where email = $1)
`

func (q *Queries) CheckForExistingUser(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRow(ctx, checkForExistingUser, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createUser = `-- name: CreateUser :one
insert into users (
    email,
    hashed_password,
    email_verified,
    admin
) values (
    $1, $2, $3, $4
) returning id, email, hashed_password, email_verified, created_at, updated_at, admin
`

type CreateUserParams struct {
	Email          string      `json:"email"`
	HashedPassword string      `json:"hashed_password"`
	EmailVerified  bool        `json:"email_verified"`
	Admin          pgtype.Bool `json:"admin"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.HashedPassword,
		arg.EmailVerified,
		arg.Admin,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.EmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Admin,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
delete from users
where
    email = $1
`

func (q *Queries) DeleteUser(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, deleteUser, email)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
select id, email, hashed_password, email_verified, created_at, updated_at, admin from users
where
    email = $1
limit 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.EmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Admin,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
select id, email, hashed_password, email_verified, created_at, updated_at, admin from users
where
    id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.EmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Admin,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
update users
set
    email = coalesce($1, email),
    hashed_password = coalesce($2, hashed_password),
    email_verified = coalesce($3, email_verified),
    admin = coalesce($4, admin),
    updated_at = current_timestamp
where
    id = $5
returning id, email, hashed_password, email_verified, created_at, updated_at, admin
`

type UpdateUserParams struct {
	Email          pgtype.Text `json:"email"`
	HashedPassword pgtype.Text `json:"hashed_password"`
	EmailVerified  pgtype.Bool `json:"email_verified"`
	Admin          pgtype.Bool `json:"admin"`
	ID             int64       `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.Email,
		arg.HashedPassword,
		arg.EmailVerified,
		arg.Admin,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.EmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Admin,
	)
	return i, err
}
