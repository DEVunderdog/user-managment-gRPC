-- name: CreateVerificationCode :one
insert into verification_codes (
    user_id,
    code,
    expires_at
) values (
    $1, $2, $3
) returning *;

-- name: UpdateVerificationCodeStatus :one
update verification_codes
set
    is_used = sqlc.arg('is_used'),
    updated_at = current_timestamp
where id = $1
returning *;