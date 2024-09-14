-- name: CreateJWTKey :one
insert into jwtkeys (
    public_key,
    private_key,
    algorithm,
    is_active,
    expires_at
) values (
    $1, $2, $3, $4, $5
) returning *;

-- name: CountJWTKeys :one
SELECT COUNT(*) FROM jwtkeys;

-- name: GetActiveKey :one
select * from jwtkeys
where is_active = sqlc.arg('is_active')
ORDER BY created_at DESC
LIMIT 1;


-- name: UpdateJWTKeys :one
update jwtkeys
set is_active = sqlc.arg('is_active')
where id = sqlc.arg('id')
returning *;