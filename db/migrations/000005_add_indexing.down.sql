DROP INDEX index_users_friend_count;
DROP INDEX index_users_friend_created_at;
DROP INDEX index_users_name;

DROP INDEX index_friendships_user_addeby_id;
DROP INDEX index_friendships_created_at;

DROP INDEX index_posts_user_id;
DROP INDEX index_posts_content;
DROP INDEX index_posts_tag;
DROP INDEX index_posts_created_at;

DROP INDEX index_comments_created_at;
DROP INDEX index_comments_user_post_id;

DROP EXTENSION pg_trgm;
