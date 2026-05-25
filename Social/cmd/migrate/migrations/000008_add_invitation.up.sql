CREATE TABLE IF NOT EXISTS user_invitations(
    token bytea,
    user_id bigint NOT NULL
)