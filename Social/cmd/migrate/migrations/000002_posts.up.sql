CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY, 
    content text NOT NULL, 
    title varchar(255) NOT NULL, 
    user_id bigint NOT NULL,
    tags text[] NOT NULL, 
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT now()
)


-- type Post struct {
--     ID int64 `json:"id"` // postgres generates this — you need it back
--     Content string `json:"content"`
--     Title string `json:"title"`
--     UserId int64 `json:"user_id"`
--     Tags []string `json:"tags"`
--     CreatedAt string `json:"created_at"` // postgres generates this — you need it back
--     UpdatedAt string `json:"updated_at"` // postgres generates this — you need it back
-- }