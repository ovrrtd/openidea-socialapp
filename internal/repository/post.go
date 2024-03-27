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
	var conditions []string
	var args []interface{}
	// Add conditions based on filter criteria
	argIndex := 1 // Start index for placeholder arguments
	fmt.Println("filter.Tags ", filter.Tags)
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
	var whereClause = "WHERE ((f.user_id = p.user_id OR f.added_by = p.user_id) OR p.user_id = " + fmt.Sprint(filter.UserID) + ") "
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	// Construct the LIMIT and OFFSET clauses
	limitOffsetClause := fmt.Sprintf("LIMIT $%d ", argIndex)
	argIndex++
	limitOffsetClause += fmt.Sprintf("OFFSET $%d ", argIndex)
	argIndex++

	// Construct the final query for products
	query := `SELECT 
			distinct(p.id),
			p.content_html,
			p.tags,
			p.user_id,
			p.created_at,
			p.updated_at,
			u.id as u_id,
			u.name as u_name,
			u.image_url as u_image_url,
			u.friend_count as u_friend_count,
			u.created_at as u_created_at
	FROM posts p
	LEFT JOIN users u ON p.user_id = u.id 
	LEFT JOIN friendships f ON p.user_id = f.user_id OR p.user_id = f.added_by
	 ` + whereClause + ` ORDER BY 
    p.created_at DESC ` + limitOffsetClause
	// Construct the query to get total product count
	countQuery := "SELECT COUNT(distinct(p.id)) FROM posts p LEFT JOIN users u ON p.user_id = u.id LEFT JOIN friendships f ON p.user_id = f.user_id OR p.user_id = f.added_by " + whereClause
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
	postIds := []int64{}
	for rows.Next() {
		post := entity.Post{}
		user := entity.User{}

		// var commentString []string
		// Add more variables as needed for other columns
		err := rows.Scan(
			&post.ID,
			&post.ContentHtml,
			&post.Tags,
			&post.UserID,
			&post.CreatedAt,
			&post.UpdatedAt,
			&user.ID,
			&user.Name,
			&user.ImageUrl,
			&user.FriendCount,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to scan product")
		}
		postIds = append(postIds, post.ID)
		post.Creator = user
		posts = append(posts, post)
	}

	// query comment by post ids
	if len(postIds) != 0 {
		query = `
		SELECT
			c.content,
			c.post_id,
			c.created_at,
			u.id as u_id,
			u.name as u_name,
			u.image_url as u_image_url,
			u.friend_count as u_friend_count
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id IN (%s)
		ORDER BY c.created_at DESC
	`
		idsStr := ""
		for i, id := range postIds {
			idsStr += fmt.Sprintf("%d", id)
			if i != len(postIds)-1 {
				idsStr += ", "
			}
		}
		rows, err = r.db.QueryContext(ctx, fmt.Sprintf(query, idsStr))

		if err != nil {
			return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to execute query")
		}

		hashPostID := make(map[int64][]entity.Comment)
		for rows.Next() {
			comment := entity.Comment{}
			user := entity.User{}
			err := rows.Scan(
				&comment.Content,
				&comment.PostID,
				&comment.CreatedAt,
				&user.ID,
				&user.Name,
				&user.ImageUrl,
				&user.FriendCount,
			)
			if err != nil {
				return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to scan product")
			}

			comment.Creator = user
			r.logger.Debug().Msgf("user: %+v", user)
			if hashPostID[comment.PostID] == nil {
				hashPostID[comment.PostID] = []entity.Comment{}
			}
			hashPostID[comment.PostID] = append(hashPostID[comment.PostID], comment)
		}
		for i, p := range posts {
			if hashPostID[p.ID] != nil {
				posts[i].Comments = hashPostID[p.ID]
			}
		}
	}

	// Execute the count query to get total product count
	var totalCount int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, errors.Wrap(err, "failed to get total product count")
	}

	meta := common.Meta{
		Total:  totalCount,
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
