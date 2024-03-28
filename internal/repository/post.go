package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"socialapp/internal/helper/common"
	"socialapp/internal/helper/errorer"
	"socialapp/internal/model/entity"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type PostRepository interface {
	FindAll(ctx context.Context, filter entity.FindAllPostRequest) ([]entity.Post, *common.Meta, int, error)
	FindByID(ctx context.Context, id int64) (*entity.Post, int, error)
	CreatePost(ctx context.Context, ent entity.Post) (int, error)
	CreateComment(ctx context.Context, ent entity.Comment) (int, error)
}

func NewPostRepository(logger zerolog.Logger, db *sql.DB) PostRepository {
	return &PostRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

type PostRepositoryImpl struct {
	logger zerolog.Logger
	db     *sql.DB
}

func (r *PostRepositoryImpl) FindAll(ctx context.Context, filter entity.FindAllPostRequest) ([]entity.Post, *common.Meta, int, error) {
	type PostRaw struct {
		ID                   int64
		ContentHTML          string
		Tags                 string
		UserID               int64
		CreatedAt            int64
		UpdatedAt            int64
		UserIDPost           int
		UserName             string
		UserImageURL         string
		UserFriendCnt        int64
		UserCreatedAt        int64
		CommentID            *int64
		CommentContent       *string
		CommentPostID        *int64
		CommentCreatedAt     *int64
		CommentUserID        *int64
		CommentUserName      *string
		CommentUserImageURL  *string
		CommentUserFriendCnt *int64
	}
	var conditions []string
	var args []interface{}
	// Add conditions based on filter criteria
	argIndex := 1 // Start index for placeholder arguments
	if len(filter.Tags) > 0 {
		tagConditions := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			tagConditions[i] = "LOWER(p.tags) LIKE $" + fmt.Sprint(argIndex)
			args = append(args, "%"+strings.ToLower(tag)+"%")
			argIndex++
		}
		conditions = append(conditions, "("+strings.Join(tagConditions, " AND ")+")")
	}

	if filter.Search != "" {
		conditions = append(conditions, "LOWER(p.content_html) LIKE $"+fmt.Sprint(argIndex))
		args = append(args, "%"+strings.ToLower(filter.Search)+"%")
		argIndex++
	}

	// Construct the WHERE clause
	var whereClause string
	if len(conditions) > 0 {
		whereClause += "WHERE " + strings.Join(conditions, " AND ")
	}

	// Construct the LIMIT and OFFSET clauses
	limitOffsetClause := fmt.Sprintf("LIMIT $%d ", argIndex)
	argIndex++
	limitOffsetClause += fmt.Sprintf("OFFSET $%d ", argIndex)
	argIndex++

	// Construct the final query for products
	query := `SELECT 
			p.id,
			p.content_html,
			p.tags,
			p.user_id,
			p.created_at,
			p.updated_at,
			u.id as u_id,
			u.name as u_name,
			u.image_url as u_image_url,
			u.friend_count as u_friend_count,
			u.created_at as u_created_at,
			c.id,
			c.content,
			c.post_id,
			c.created_at,
			cu.id as u_id,
			cu.name as u_name,
			cu.image_url as u_image_url,
			cu.friend_count as u_friend_count
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id 
	LEFT JOIN comments c ON p.id = c.post_id
	LEFT JOIN users cu ON c.user_id = cu.id
	 ` + whereClause + ` ORDER BY 
    p.created_at DESC ` + limitOffsetClause
	// Construct the query to get total product count
	posts := []entity.Post{}
	// Execute the main query
	argsQuery := []interface{}{}
	argsQuery = append(argsQuery, args...)
	argsQuery = append(argsQuery, filter.Limit)
	argsQuery = append(argsQuery, filter.Offset)
	rows, err := r.db.QueryContext(ctx, query, argsQuery...)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()
	postIDs := make(map[int64]int)
	for rows.Next() {
		postRaw := PostRaw{}

		// var commentString []string
		// Add more variables as needed for other columns
		err := rows.Scan(
			&postRaw.ID,
			&postRaw.ContentHTML,
			&postRaw.Tags,
			&postRaw.UserID,
			&postRaw.CreatedAt,
			&postRaw.UpdatedAt,
			&postRaw.UserID,
			&postRaw.UserName,
			&postRaw.UserImageURL,
			&postRaw.UserFriendCnt,
			&postRaw.UserCreatedAt,
			&postRaw.CommentID,
			&postRaw.CommentContent,
			&postRaw.CommentPostID,
			&postRaw.CommentCreatedAt,
			&postRaw.CommentUserID,
			&postRaw.CommentUserName,
			&postRaw.CommentUserImageURL,
			&postRaw.CommentUserFriendCnt,
		)
		if err != nil {
			return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to scan product")
		}

		if postIDs[postRaw.ID] == 0 {
			post := entity.Post{
				ID:          postRaw.ID,
				ContentHtml: postRaw.ContentHTML,
				Tags:        postRaw.Tags,
				UserID:      postRaw.UserID,
				CreatedAt:   postRaw.CreatedAt,
				UpdatedAt:   postRaw.UpdatedAt,
				Comments:    []entity.Comment{},
			}
			post.Creator = entity.User{
				ID:          postRaw.UserID,
				Name:        postRaw.UserName,
				ImageUrl:    postRaw.UserImageURL,
				FriendCount: postRaw.UserFriendCnt,
				CreatedAt:   postRaw.UserCreatedAt,
			}
			posts = append(posts, post)
			postIDs[postRaw.ID] = len(posts) - 1
		}

		if postRaw.CommentID != nil {
			posts[postIDs[postRaw.ID]].Comments = append(posts[postIDs[postRaw.ID]].Comments, entity.Comment{
				ID:        *postRaw.CommentID,
				Content:   *postRaw.CommentContent,
				PostID:    *postRaw.CommentPostID,
				CreatedAt: *postRaw.CommentCreatedAt,
				Creator: entity.User{
					ID:          *postRaw.CommentUserID,
					Name:        *postRaw.CommentUserName,
					ImageUrl:    *postRaw.CommentUserImageURL,
					FriendCount: *postRaw.CommentUserFriendCnt,
				},
			})
		}

	}

	meta := common.Meta{
		Total:  len(posts),
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}

	return posts, &meta, http.StatusOK, nil
}

func (r *PostRepositoryImpl) CreatePost(ctx context.Context, ent entity.Post) (int, error) {
	ent.CreatedAt = time.Now().UnixMilli()
	ent.UpdatedAt = time.Now().UnixMilli()
	// insert post
	err := r.db.QueryRowContext(
		ctx,
		`
		INSERT INTO posts (content_html, tags, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id
		`,
		&ent.ContentHtml, &ent.Tags, &ent.UserID, &ent.CreatedAt, &ent.UpdatedAt,
	).Scan(&ent.ID)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return http.StatusOK, nil
}

// create comment
func (r *PostRepositoryImpl) CreateComment(ctx context.Context, ent entity.Comment) (int, error) {
	ent.CreatedAt = time.Now().UnixMilli()
	ent.UpdatedAt = time.Now().UnixMilli()

	var tgtUserId int64
	// get post by id
	err := r.db.QueryRowContext(
		ctx,
		"SELECT user_id FROM posts WHERE id = $1",
		&ent.PostID,
	).Scan(&tgtUserId)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if tgtUserId != ent.UserID {
		// check if user is friend
		var id int64
		_ = r.db.QueryRowContext(
			ctx,
			"SELECT id FROM friendships WHERE (user_id = $1 AND added_by = $2) OR (user_id = $2 AND added_by = $1)",
			&ent.UserID, &tgtUserId,
		).Scan(&id)

		if id == 0 {
			return http.StatusBadRequest, errors.Wrap(errorer.ErrBadRequest, "user is not friend")
		}
	}
	// insert comment
	err = r.db.QueryRowContext(
		ctx,
		`
		INSERT INTO comments (content, created_at, updated_at, post_id, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id
		`,
		&ent.Content, &ent.CreatedAt, &ent.UpdatedAt, &ent.PostID, &ent.UserID,
	).Scan(&ent.ID)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return http.StatusOK, nil
}

// find by id
func (r *PostRepositoryImpl) FindByID(ctx context.Context, id int64) (*entity.Post, int, error) {
	var post entity.Post
	row := r.db.QueryRowContext(ctx, "SELECT id, content_html, tags, user_id, created_at, updated_at FROM posts WHERE id = $1", id)
	err := row.Scan(&post.ID, &post.ContentHtml, &post.Tags, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	return &post, http.StatusOK, nil
}
