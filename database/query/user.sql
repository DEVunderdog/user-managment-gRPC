-- name: CreateUser :one
insert into users (
    email,
    hashed_password,
    email_verified,
    admin
) values (
    $1, $2, $3, $4
) returning *;

-- name: GetUserByID :one
select * from users
where
    id = sqlc.arg('id');

-- name: CheckForExistingUser :one
select exists(select 1 from users where email = sqlc.arg('email'));

-- name: GetUserByEmail :one
select * from users
where
    email = sqlc.arg('email')
limit 1;

-- name: UpdateUser :one
update users
set
    email = coalesce(sqlc.narg(email), email),
    hashed_password = coalesce(sqlc.narg(hashed_password), hashed_password),
    email_verified = coalesce(sqlc.narg(email_verified), email_verified),
    admin = coalesce(sqlc.narg(admin), admin),
    updated_at = current_timestamp
where
    id = sqlc.arg(id)
returning *;

-- name: DeleteUser :exec
delete from users
where
    email = sqlc.arg('email');