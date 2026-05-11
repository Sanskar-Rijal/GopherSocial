-- emails are case insensitive sanskar@gmail.com is same as SanSkar@gmail.com so we want to store it correctly
CREATE EXTENSION IF NOT EXISTS citext;


CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password bytea NOT NULL, 
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

    -- ID int64 `json:"id"`
    -- Username string `json:"name"`
    -- Email string `json:"email"`
    -- Password string `json:"-"` //we will not return password to the user
    -- CreatedAt string `json:"created_at"