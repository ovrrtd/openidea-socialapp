CREATE TABLE IF NOT EXISTS FRIENDSHIPS (
    ID SERIAL PRIMARY KEY NOT NULL,
    USER_ID INTEGER NOT NULL,
    ADDED_BY INTEGER NOT NULL,
    UPDATED_AT BIGINT NOT NULL,
    CREATED_AT BIGINT NOT NULL,
    CONSTRAINT fk_friendships_user FOREIGN KEY(USER_ID) REFERENCES USERS(id),
    CONSTRAINT fk_friendships_added FOREIGN KEY(ADDED_BY) REFERENCES USERS(id)
);