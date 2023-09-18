-- name: CreateVerifyEmail :one
insert into verify_emails (
                           username,
                           email,
                           secret_code
) values (
$1, $2, $3
) returning *;

-- name: UpdateVerifyEmail :one
update verify_emails
set is_used = TRUE
where id = @id and secret_code = @secret_code and is_used = FALSE and expired_at > now()
returning *;