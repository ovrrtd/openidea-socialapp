ALTER TABLE FRIENDSHIPS DROP CONSTRAINT fk_friendships_user;
ALTER TABLE FRIENDSHIPS DROP CONSTRAINT fk_friendships_added;

DROP TABLE FRIENDSHIPS;