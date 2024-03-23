package service

import (
	"context"
	"mime/multipart"
	"socialapp/internal/model/request"
	"socialapp/internal/model/response"
	"socialapp/internal/repository"

	"github.com/rs/zerolog"
)

type Service interface {
	// User
	Register(ctx context.Context, payload request.Register) (*response.Register, int, error)
	Login(ctx context.Context, payload request.Login) (*response.Login, int, error)
	GetUserByID(ctx context.Context, id int64) (*response.User, int, error)
	UpdateAccount(ctx context.Context, payload request.UpdateAccount) (*response.User, int, error)
	LinkEmail(ctx context.Context, payload request.LinkEmail) (*response.User, int, error)
	LinkPhone(ctx context.Context, payload request.LinkPhone) (*response.User, int, error)
	// s3
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, int, error)
}

type Config struct {
	Salt      int
	JwtSecret string
}

type service struct {
	cfg      Config
	log      zerolog.Logger
	userRepo repository.UserRepository
	s3Repo   repository.S3Repository
}

func New(cfg Config, logger zerolog.Logger, userRepo repository.UserRepository, s3Repo repository.S3Repository) Service {
	return &service{
		cfg:      cfg,
		log:      logger,
		userRepo: userRepo,
		s3Repo:   s3Repo,
	}
}