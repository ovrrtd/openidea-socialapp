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

	"github.com/pkg/errors"
)

func (s *service) FindAllFriendships(ctx context.Context, filter request.FindAllFriendships) ([]response.FindAllFriendships, *common.Meta, int, error) {
	err := validator.ValidateStruct(&filter)
	if err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	ent, meta, code, err := s.friendshipRepo.FindAll(ctx, entity.FindAllFriendshipRequest{
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		UserID:     filter.UserID,
		OnlyFriend: filter.OnlyFriend,
		Search:     filter.Search,
		SortBy:     filter.SortBy,
		OrderBy:    filter.OrderBy,
	})

	if err != nil {
		return nil, nil, code, err
	}

	ret := make([]response.FindAllFriendships, len(ent))
	for i, e := range ent {
		ret[i] = response.FindAllFriendships{
			ID:          strconv.Itoa(int(e.ID)),
			Name:        e.Name,
			ImageUrl:    e.ImageUrl,
			FriendCount: e.FriendCount,
			CreatedAt:   common.UnixMilliToISO8601(e.CreatedAt),
		}
	}
	return ret, meta, code, nil
}

func (s *service) CreateFriendship(ctx context.Context, payload request.CreateFriendship) (int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}
	userID, err := strconv.Atoi(payload.UserID)
	if err != nil {
		return http.StatusNotFound, errors.Wrap(errors.New("invalid user id"), "invalid user id")
	}
	addedBy := payload.AddedBy
	s.log.Debug().Msgf("userID: %d, addedBy: %d", userID, addedBy)
	if userID == 0 || addedBy == 0 {
		return http.StatusBadRequest, errors.Wrap(errors.New("invalid user id"), "invalid user id")
	}

	if int64(userID) == addedBy {
		return http.StatusBadRequest, errors.Wrap(errors.New("can not add yourself"), "can not add yourself")
	}

	return s.friendshipRepo.CreateFriendship(ctx, int64(userID), addedBy)
}

func (s *service) DeleteFriendship(ctx context.Context, payload request.DeleteFriendship) (int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	friend1, _ := strconv.Atoi(payload.Friend1)
	friend2 := payload.Friend2
	if friend1 == 0 || friend2 == 0 {
		return http.StatusBadRequest, errors.Wrap(errors.New("invalid user id"), "invalid user id")
	}

	if int64(friend1) == friend2 {
		return http.StatusBadRequest, errors.Wrap(errors.New("can not delete yourself"), "can not delete yourself")
	}

	return s.friendshipRepo.DeleteFriendship(ctx, int64(friend1), friend2)
}
