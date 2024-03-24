CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX index_users_friend_count ON USERS USING BRIN (FRIEND_COUNT);
CREATE INDEX index_users_friend_created_at ON USERS USING BRIN (CREATED_AT);
CREATE INDEX index_users_name ON USERS USING GIN (NAME gin_trgm_ops);

CREATE INDEX index_friendships_user_addeby_id ON FRIENDSHIPS (USER_ID,ADDED_BY);
CREATE INDEX index_friendships_created_at ON FRIENDSHIPS (CREATED_AT);

CREATE INDEX index_posts_user_id ON POSTS (USER_ID);
CREATE INDEX index_posts_content ON POSTS USING GIN (CONTENT_HTML gin_trgm_ops);
CREATE INDEX index_posts_tag ON POSTS USING GIN (TAGS gin_trgm_ops);
CREATE INDEX index_posts_created_at ON posts USING BRIN (CREATED_AT);

CREATE INDEX index_comments_created_at ON COMMENTS USING BRIN (CREATED_AT);
CREATE INDEX index_comments_user_post_id ON COMMENTS (USER_ID, POST_ID);

