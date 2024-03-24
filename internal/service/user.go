package service

import (
	"context"
	"net/http"
	"regexp"
	"socialapp/internal/helper/common"
	"socialapp/internal/helper/errorer"
	"socialapp/internal/helper/jwt"
	"socialapp/internal/helper/validator"
	"socialapp/internal/model/entity"
	"socialapp/internal/model/request"
	"socialapp/internal/model/response"
	"time"

	jwtV5 "github.com/golang-jwt/jwt/v5"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Register to register a new user by email and password
func (s *service) Register(ctx context.Context, payload request.Register) (*response.Register, int, error) {
	err := validator.ValidateStruct(&payload)

	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	if payload.CredentialType != "email" && payload.CredentialType != "phone" {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}
	ent := entity.User{
		Name: payload.Name,
	}

	if payload.CredentialType == "email" {
		// validate email form
		regex := regexp.MustCompile(common.RegexEmailPattern)
		if !regex.MatchString(payload.CredentialValue) {
			return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidEmail, errorer.ErrInvalidEmail.Error())
		}

		exist, code, err := s.userRepo.FindByEmail(ctx, payload.CredentialValue)

		if err != nil && code != http.StatusNotFound {
			return nil, code, err
		}
		if exist != nil {
			return nil, http.StatusConflict, errors.Wrap(errorer.ErrEmailExist, errorer.ErrEmailExist.Error())
		}

		ent.Email = payload.CredentialValue
	}

	if payload.CredentialType == "phone" {
		// validate phone number
		if pl := len(payload.CredentialValue); pl < 7 || pl > 13 {
			return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidPhone, errorer.ErrInvalidPhone.Error())
		}

		if !common.ValidatePhoneNumber(payload.CredentialValue) {
			return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidPhone, errorer.ErrInvalidPhone.Error())
		}

		exist, code, err := s.userRepo.FindByPhone(ctx, payload.CredentialValue)

		if err != nil && code != http.StatusNotFound {
			return nil, code, err
		}
		if exist != nil {
			return nil, http.StatusConflict, errors.Wrap(errorer.ErrPhoneExist, errorer.ErrPhoneExist.Error())
		}
		ent.Phone = payload.CredentialValue
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), s.cfg.Salt)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, err.Error())
	}
	ent.Password = string(hashedPassword)

	user, code, err := s.userRepo.Register(ctx, ent)

	if err != nil {
		return nil, code, err
	}

	// TODO: generate access token
	userClaims := common.UserClaims{
		Id: user.ID,
		RegisteredClaims: jwtV5.RegisteredClaims{
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	tokenString, err := jwt.GenerateJwt(userClaims, s.cfg.JwtSecret)

	if err != nil {
		return nil, code, errors.Wrap(err, err.Error())
	}

	return &response.Register{
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		AccessToken: tokenString,
	}, code, nil
}

func (s *service) Login(ctx context.Context, payload request.Login) (*response.Login, int, error) {
	err := validator.ValidateStruct(&payload)

	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	if payload.CredentialType != "email" && payload.CredentialType != "phone" {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	user := &entity.User{}

	if payload.CredentialType == "email" {
		// validate email form
		regex := regexp.MustCompile(common.RegexEmailPattern)
		if !regex.MatchString(payload.CredentialValue) {
			return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidEmail, errorer.ErrInvalidEmail.Error())
		}

		usr, code, err := s.userRepo.FindByEmail(ctx, payload.CredentialValue)

		if err != nil {
			return nil, code, err
		}
		user = usr
	}

	if payload.CredentialType == "phone" {
		// validate phone number
		if pl := len(payload.CredentialValue); pl < 7 || pl > 13 {
			return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidPhone, errorer.ErrInvalidPhone.Error())
		}

		if !common.ValidatePhoneNumber(payload.CredentialValue) {
			return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidPhone, errorer.ErrInvalidPhone.Error())
		}

		usr, code, err := s.userRepo.FindByPhone(ctx, payload.CredentialValue)

		if err != nil {
			return nil, code, err
		}
		user = usr

	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, err.Error())
	}

	userClaims := common.UserClaims{
		Id: user.ID,
		RegisteredClaims: jwtV5.RegisteredClaims{
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	tokenString, err := jwt.GenerateJwt(userClaims, s.cfg.JwtSecret)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, err.Error())
	}

	return &response.Login{
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		AccessToken: tokenString,
	}, http.StatusOK, nil
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*response.User, int, error) {
	user, code, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, code, err
	}
	return &response.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		ImageUrl:    user.ImageUrl,
		FriendCount: user.FriendCount,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, code, nil
}

func (s *service) UpdateAccount(ctx context.Context, payload request.UpdateAccount) (*response.User, int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	if !common.ValidateUrl(payload.ImageUrl) {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidImageUrl, errorer.ErrInvalidImageUrl.Error())
	}

	ent, code, err := s.userRepo.FindByID(ctx, payload.ID)

	if err != nil {
		return nil, code, err
	}
	ent.ImageUrl = payload.ImageUrl
	ent.Name = payload.Name

	ent, code, err = s.userRepo.UpdateByID(ctx, *ent)
	if err != nil {
		return nil, code, err
	}

	return &response.User{
		ID:          ent.ID,
		Name:        ent.Name,
		Email:       ent.Email,
		Phone:       ent.Phone,
		ImageUrl:    ent.ImageUrl,
		FriendCount: ent.FriendCount,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
	}, code, nil
}

func (s *service) LinkPhone(ctx context.Context, payload request.LinkPhone) (*response.User, int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	// validate phone number
	if pl := len(payload.Phone); pl < 7 || pl > 13 {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidPhone, errorer.ErrInvalidPhone.Error())
	}

	if !common.ValidatePhoneNumber(payload.Phone) {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidPhone, errorer.ErrInvalidPhone.Error())
	}

	ent, code, err := s.userRepo.FindByID(ctx, payload.ID)

	if err != nil {
		return nil, code, err
	}

	if ent.Phone != "" {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrBadRequest, errorer.ErrBadRequest.Error())
	}

	exist, _, _ := s.userRepo.FindByPhone(ctx, payload.Phone)

	if exist != nil {
		return nil, http.StatusConflict, errors.Wrap(errorer.ErrPhoneExist, errorer.ErrPhoneExist.Error())
	}

	ent.Phone = payload.Phone
	ent, code, err = s.userRepo.UpdateByID(ctx, *ent)
	if err != nil {
		return nil, code, err
	}

	return &response.User{
		ID:          ent.ID,
		Name:        ent.Name,
		Email:       ent.Email,
		Phone:       ent.Phone,
		ImageUrl:    ent.ImageUrl,
		FriendCount: ent.FriendCount,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
	}, code, nil
}

func (s *service) LinkEmail(ctx context.Context, payload request.LinkEmail) (*response.User, int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	regex := regexp.MustCompile(common.RegexEmailPattern)
	if !regex.MatchString(payload.Email) {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidEmail, errorer.ErrInvalidEmail.Error())
	}

	ent, code, err := s.userRepo.FindByID(ctx, payload.ID)
	if err != nil {
		return nil, code, err
	}

	if ent.Email != "" {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrBadRequest, errorer.ErrBadRequest.Error())
	}

	exist, _, _ := s.userRepo.FindByEmail(ctx, payload.Email)
	if exist != nil {
		return nil, http.StatusConflict, errors.Wrap(errorer.ErrEmailExist, errorer.ErrEmailExist.Error())
	}

	ent.Email = payload.Email
	ent, code, err = s.userRepo.UpdateByID(ctx, *ent)
	if err != nil {
		return nil, code, err
	}

	return &response.User{
		ID:          ent.ID,
		Name:        ent.Name,
		Email:       ent.Email,
		Phone:       ent.Phone,
		ImageUrl:    ent.ImageUrl,
		FriendCount: ent.FriendCount,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
	}, code, nil
}
