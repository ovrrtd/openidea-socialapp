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

type FriendshipRepository interface {
	CreateFriendship(ctx context.Context, userID int64, addedBy int64) (int, error)
	DeleteFriendship(ctx context.Context, friend1 int64, friend2 int64) (int, error)
	FindAll(ctx context.Context, filter entity.FindAllFriendshipRequest) ([]entity.User, *common.Meta, int, error)
}

func NewFriendshipRepository(logger zerolog.Logger, db *sql.DB) FriendshipRepository {
	return &FriendshipRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

type FriendshipRepositoryImpl struct {
	logger zerolog.Logger
	db     *sql.DB
}

func (r *FriendshipRepositoryImpl) CreateFriendship(ctx context.Context, userID int64, addedBy int64) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	defer tx.Rollback()

	// check friendship
	frd := entity.Friendship{
		UserID:    userID,
		AddedBy:   addedBy,
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}
	err = tx.QueryRowContext(ctx,
		"INSERT INTO friendships (user_id, added_by, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		frd.UserID, frd.AddedBy, frd.CreatedAt, frd.UpdatedAt).Scan(&frd.ID)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	res, err := tx.ExecContext(ctx, "UPDATE users SET friend_count = friend_count + 1 WHERE id = $1 OR id = $2", frd.UserID, frd.AddedBy)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	row, err := res.RowsAffected()

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if row < 1 {
		return http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	return http.StatusOK, nil
}

func (r *FriendshipRepositoryImpl) DeleteFriendship(ctx context.Context, friend1 int64, friend2 int64) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	defer tx.Rollback()

	// Delete friendship
	res, err := tx.ExecContext(ctx, "DELETE FROM friendships WHERE (user_id = $1 AND added_by = $2) OR (user_id = $2 AND added_by = $1)", friend1, friend2)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	row, err := res.RowsAffected()

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if row == 0 {
		return http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
	}

	res, err = tx.ExecContext(ctx, "UPDATE users SET friend_count = friend_count - 1 WHERE id = $1 OR id = $2", friend1, friend2)
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	row, err = res.RowsAffected()

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if row < 1 {
		return http.StatusNotFound, errors.Wrap(errorer.ErrNotFound, errorer.ErrNotFound.Error())
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	return http.StatusOK, nil
}

func (r *FriendshipRepositoryImpl) FindAll(ctx context.Context, filter entity.FindAllFriendshipRequest) ([]entity.User, *common.Meta, int, error) {
	var conditions []string
	var args []interface{}
	// Add conditions based on filter criteria
	argIndex := 1 // Start index for placeholder arguments

	if filter.Search != "" {
		conditions = append(conditions, "LOWER(u.name) LIKE $"+fmt.Sprint(argIndex))
		args = append(args, "%"+strings.ToLower(filter.Search)+"%")
		argIndex++
	}
	fmt.Println(filter.OnlyFriend)
	if filter.OnlyFriend && filter.UserID != 0 {
		conditions = append(conditions, "(f.user_id = $"+fmt.Sprint(argIndex)+" OR f.added_by = $"+fmt.Sprint(argIndex)+")")
		args = append(args, filter.UserID)
		argIndex++
	}

	// Construct the WHERE clause
	var whereClause = "WHERE u.id != " + fmt.Sprint(filter.UserID)
	if len(conditions) > 0 {
		whereClause += " AND " + strings.Join(conditions, " AND ")
	}

	var orderByClause string
	if filter.SortBy != "" {
		orderByClause = "ORDER BY " + "u.created_at"

		if filter.SortBy == "friendCount" {
			orderByClause = "ORDER BY " + "u.friend_count"
		}
		orderBy := "DESC"

		if orderBy != "" && (strings.ToUpper(filter.OrderBy) == "ASC" || strings.ToUpper(filter.OrderBy) == "DESC") {
			orderBy = strings.ToUpper(filter.OrderBy)
		}
		orderByClause += " " + orderBy
	}

	// Construct the LIMIT and OFFSET clauses
	limitOffsetClause := fmt.Sprintf("LIMIT $%d ", argIndex)
	argIndex++
	limitOffsetClause += fmt.Sprintf("OFFSET $%d ", argIndex)
	argIndex++

	query := `
		SELECT
			distinct 
			u.id,
			u.name,
			u.image_url,
			u.friend_count,
			u.created_at
		FROM users u
		LEFT JOIN friendships f ON u.id = f.user_id OR u.id = f.added_by
		` + whereClause + " " + orderByClause + " " + limitOffsetClause

	queryCount := `SELECT COUNT(distinct(u.id)) FROM users u LEFT JOIN friendships f ON u.id = f.user_id OR u.id = f.added_by ` + whereClause

	argsQuery := []interface{}{}
	argsQuery = append(argsQuery, args...)
	argsQuery = append(argsQuery, filter.Limit)
	argsQuery = append(argsQuery, filter.Offset)
	r.logger.Debug().Msgf("query: %s args: %+v", query, argsQuery)
	r.logger.Debug().Msgf("query: %s args: %+v", queryCount, argsQuery)

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, argsQuery...)

	if err != nil {
		return nil, nil, 0, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	defer rows.Close()

	// Process the results
	users := []entity.User{}

	for rows.Next() {
		user := entity.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.ImageUrl, &user.FriendCount, &user.CreatedAt)
		if err != nil {
			return nil, nil, 0, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
		}
		users = append(users, user)
	}

	// Get total count
	var totalCount int
	err = r.db.QueryRowContext(ctx, queryCount, args...).Scan(&totalCount)

	if err != nil {
		return nil, nil, 0, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	meta := common.Meta{
		Total:  totalCount,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}
	return users, &meta, http.StatusOK, nil
}
