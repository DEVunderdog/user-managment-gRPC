-- name: CreateSession :one
insert into sessions (
    user_id,
    access_token,
    refresh_token,
    access_token_expires_at,
    refresh_token_expires_at,
    is_active,
    ip,
    user_agent
) values (
    $1, $2, $3, $4, $5, $6, $7, $8
) returning *;

-- name: GetUserSessions :many
select * from sessions
where
    user_id = sqlc.arg('user_id')
    AND
    is_active = sqlc.arg('is_active')
limit sqlc.arg('limit')
offset sqlc.arg('offset');

-- name: UpdateSession :one
update sessions
set
    logged_out = current_timestamp,
    updated_at = current_timestamp
where
    id = sqlc.arg('id')
returning *;

-- name: DeleteSessions :exec
delete from sessions
where id = sqlc.arg('id');