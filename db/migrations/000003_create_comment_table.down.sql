ALTER TABLE COMMENTS DROP CONSTRAINT fk_comments_user;
ALTER TABLE COMMENTS DROP CONSTRAINT fk_comments_post;

DROP TABLE COMMENTS;
