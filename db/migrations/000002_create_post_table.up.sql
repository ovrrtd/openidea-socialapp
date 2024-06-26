CREATE TABLE IF NOT EXISTS POSTS(
    ID SERIAL PRIMARY KEY,
    CONTENT_HTML VARCHAR(512) NOT NULL,
    USER_ID INTEGER NOT NULL,
    TAGS VARCHAR(512) NOT NULL DEFAULT '',
    CREATED_AT BIGINT NOT NULL,
    UPDATED_AT BIGINT NOT NULL,
    CONSTRAINT fk_posts_user FOREIGN KEY(USER_ID) REFERENCES USERS(id)
);