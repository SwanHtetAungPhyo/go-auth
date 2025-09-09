-- name: GoAuthRegister :one
INSERT INTO goauth_user (
    email,
    hash_password,
    name,
    role_name,
    metadata
) VALUES (
             @email,
             @hash_password,
             @name,
          @role_name,
             @metadata
         ) RETURNING *;

-- name: GetUserByEmail :one
SELECT
    u.*,
    ARRAY_AGG(
    CASE
        WHEN a.id IS NOT NULL THEN
            ROW(a.id, a.provider, a.provider_id, a.created_at)::goauth_account
        ELSE NULL
        END
             ) FILTER (WHERE a.id IS NOT NULL) as accounts
FROM goauth_user u
         LEFT JOIN goauth_account a ON u.id = a.user_id
WHERE u.email = @email
GROUP BY u.id;

-- name: GetUserByID :one
SELECT
    u.*,
    ARRAY_AGG(
    CASE
        WHEN a.id IS NOT NULL THEN
            ROW(a.id, a.provider, a.provider_id, a.created_at)::goauth_account
        ELSE NULL
        END
             ) FILTER (WHERE a.id IS NOT NULL) as accounts
FROM goauth_user u
         LEFT JOIN goauth_account a ON u.id = a.user_id
WHERE u.id = @user_id
GROUP BY u.id;

-- name: UpdateUser :one
UPDATE goauth_user
SET
    name = COALESCE(@name, name),
    image = COALESCE(@image, image),
    metadata = COALESCE(@metadata, metadata),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserEmailVerified :exec
UPDATE goauth_user
SET email_verified = @email_verified, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserTwoFactor :exec
UPDATE goauth_user
SET
    two_factor_enabled = @two_factor_enabled,
    two_factor_secret = @two_factor_secret,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM goauth_user WHERE id = $1;

-- sql/queries/sessions.sql
-- name: CreateSession :one
INSERT INTO goauth_session (
    user_id,
    token,
    expires_at,
    user_agent,
    ip_address
) VALUES (
             @user_id,
             @token,
             @expires_at,
             @user_agent,
             @ip_address
         ) RETURNING *;

-- name: GetSession :one
SELECT * FROM goauth_session
WHERE token = @token AND expires_at > NOW();

-- name: GetSessionByID :one
SELECT * FROM goauth_session
WHERE id = $1 AND expires_at > NOW();

-- name: DeleteSession :exec
DELETE FROM goauth_session WHERE id = $1;

-- name: DeleteUserSessions :exec
DELETE FROM goauth_session WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM goauth_session WHERE expires_at <= NOW();

-- sql/queries/accounts.sql
-- name: CreateAccount :one
INSERT INTO goauth_account (
    user_id,
    provider,
    provider_id
) VALUES (
             @user_id,
             @provider,
             @provider_id
         ) RETURNING *;

-- name: GetAccountByProvider :one
SELECT a.*, u.* FROM goauth_account a
                         JOIN goauth_user u ON a.user_id = u.id
WHERE a.provider = @provider AND a.provider_id = @provider_id;

-- name: DeleteAccount :exec
DELETE FROM goauth_account
WHERE user_id = @user_id AND provider = @provider;

-- sql/queries/password_reset.sql
-- name: CreatePasswordResetToken :one
INSERT INTO goauth_password_reset (
    user_id,
    token,
    expires_at
) VALUES (
             @user_id,
             @token,
             @expires_at
         ) RETURNING *;

-- name: GetPasswordResetToken :one
SELECT * FROM goauth_password_reset
WHERE token = @token AND expires_at > NOW();

-- name: DeletePasswordResetToken :exec
DELETE FROM goauth_password_reset WHERE token = @token;

-- name: DeleteUserPasswordResetTokens :exec
DELETE FROM goauth_password_reset WHERE user_id = @user_id;

-- name: DeleteExpiredPasswordResetTokens :exec
DELETE FROM goauth_password_reset WHERE expires_at <= NOW();

-- sql/queries/email_verification.sql
-- name: CreateEmailVerificationToken :one
INSERT INTO goauth_email_verification (
    user_id,
    token,
    expires_at
) VALUES (
             @user_id,
             @token,
             @expires_at
         ) RETURNING *;

-- name: GetEmailVerificationToken :one
SELECT * FROM goauth_email_verification
WHERE token = @token AND expires_at > NOW();

-- name: DeleteEmailVerificationToken :exec
DELETE FROM goauth_email_verification WHERE token = @token;

-- name: DeleteUserEmailVerificationTokens :exec
DELETE FROM goauth_email_verification WHERE user_id = @user_id;

-- name: DeleteExpiredEmailVerificationTokens :exec
DELETE FROM goauth_email_verification WHERE expires_at <= NOW();