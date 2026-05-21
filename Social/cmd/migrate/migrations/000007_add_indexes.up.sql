-- STEP 1 — install the extension
-- pg_trgm = trigram extension
-- trigram breaks text into groups of 3 characters
-- "hello" → "hel", "ell", "llo"
-- allows fuzzy search and LIKE search to be fast
CREATE EXTENSION IF NOT EXISTS pg_trgm;


-- STEP 2 — index on comments content
-- GIN = Generalized Inverted Index
-- good for text search and arrays
-- gin_trgm_ops = use trigram method
-- now searching comments by content is fast
-- SELECT * FROM comments WHERE content LIKE '%hello%' ← fast now
CREATE INDEX idx_comments_content ON comments USING gin (content gin_trgm_ops);

-- STEP 3 — index on posts title
-- same as above but for post titles
-- searching posts by title is now fast
-- SELECT * FROM posts WHERE title LIKE '%hello%' ← fast now
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);


-- STEP 4 — index on posts tags
-- tags is an array column
-- GIN index works great for arrays
-- SELECT * FROM posts WHERE 'go' = ANY(tags) ← fast now
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);


-- STEP 5 — index on users username
-- simple btree index (default)
-- searching users by username is fast
-- SELECT * FROM users WHERE username = 'sanskar' ← fast now
-- also makes login faster
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

-- STEP 6 — index on posts user_id
-- when you get all posts by a user
-- SELECT * FROM posts WHERE user_id = 1 ← fast now
-- very common query in social media app
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);

-- STEP 7 — index on comments post_id
-- when you get all comments for a post
-- SELECT * FROM comments WHERE post_id = 1 ← fast now
-- you do this every time you load a post
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);