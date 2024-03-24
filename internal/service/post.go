package service

import (
	"context"
	"net/http"
	"socialapp/internal/helper/common"
	"socialapp/internal/helper/errorer"
	"socialapp/internal/helper/validator"
	"socialapp/internal/model/entity"
	"socialapp/internal/model/request"
	"socialapp/internal/model/response"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (s *service) CreatePost(ctx context.Context, payload request.CreatePost) (int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}
	for _, tag := range payload.Tags {
		if tag == "" {
			return http.StatusBadRequest, errors.Wrap(errorer.ErrBadRequest, errorer.ErrBadRequest.Error())
		}
	}

	// insert post
	return s.postRepo.CreatePost(ctx, entity.Post{
		ContentHtml: payload.ContentHtml,
		UserID:      payload.UserID,
		Tags:        strings.Join(payload.Tags, ","),
	})
}

func (s *service) FindAllPost(ctx context.Context, payload request.FindAllPost) ([]response.GetPosts, *common.Meta, int, error) {
	ent, meta, code, err := s.postRepo.FindAll(ctx, entity.FindAllPostRequest{
		Limit:  payload.Limit,
		Offset: payload.Offset,
		Tags:   payload.Tags,
		Search: payload.Search,
		UserID: payload.UserID,
	})

	if err != nil {
		return nil, nil, code, err
	}

	posts := make([]response.GetPosts, len(ent))

	for i, e := range ent {
		posts[i] = response.GetPosts{
			PostID: strconv.Itoa(int(e.ID)),
			Post: struct {
				PostInHtml string   "json:\"postInHtml\""
				Tags       []string "json:\"tags\""
				CreatedAt  string   "json:\"createdAt\""
			}{
				PostInHtml: e.ContentHtml,
				Tags:       strings.Split(e.Tags, ","),
				CreatedAt:  common.UnixMilliToISO8601(e.CreatedAt),
			},
			Creator: struct {
				UserID      string "json:\"userId\""
				Name        string "json:\"name\""
				ImageURL    string "json:\"imageUrl\""
				FriendCount int    "json:\"friendCount\""
				CreatedAt   string "json:\"createdAt\""
			}{
				UserID:      strconv.Itoa(int(e.Creator.ID)),
				Name:        e.Creator.Name,
				ImageURL:    e.Creator.ImageUrl,
				FriendCount: int(e.Creator.FriendCount),
				CreatedAt:   common.UnixMilliToISO8601(e.CreatedAt),
			},
			Comments: make([]struct {
				Comment string "json:\"comment\""
				Creator struct {
					UserID      string "json:\"userId\""
					Name        string "json:\"name\""
					ImageURL    string "json:\"imageUrl\""
					FriendCount int    "json:\"friendCount\""
				} "json:\"creator\""
				CreatedAt string "json:\"createdAt\""
			}, len(e.Comments)),
		}

		for j, c := range e.Comments {
			posts[i].Comments[j] = struct {
				Comment string "json:\"comment\""
				Creator struct {
					UserID      string "json:\"userId\""
					Name        string "json:\"name\""
					ImageURL    string "json:\"imageUrl\""
					FriendCount int    "json:\"friendCount\""
				} "json:\"creator\""
				CreatedAt string "json:\"createdAt\""
			}{
				Comment: c.Content,
				Creator: struct {
					UserID      string "json:\"userId\""
					Name        string "json:\"name\""
					ImageURL    string "json:\"imageUrl\""
					FriendCount int    "json:\"friendCount\""
				}{
					UserID:      strconv.Itoa(int(c.Creator.ID)),
					Name:        c.Creator.Name,
					ImageURL:    c.Creator.ImageUrl,
					FriendCount: int(c.Creator.FriendCount),
				},
				CreatedAt: common.UnixMilliToISO8601(c.CreatedAt),
			}
		}
	}

	return posts, meta, code, nil
}

func (s *service) CreateComment(ctx context.Context, payload request.CreateComment) (int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}
	// insert comment
	postID, _ := strconv.Atoi(payload.PostID)
	_, code, err := s.postRepo.FindByID(ctx, int64(postID))

	if err != nil {
		return code, err
	}

	return s.postRepo.CreateComment(ctx, entity.Comment{
		Content: payload.Content,
		PostID:  int64(postID),
		UserID:  payload.UserID,
	})
}
