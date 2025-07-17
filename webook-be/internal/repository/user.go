package repository

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/newton-miku/webook/webook-be/internal/domain"
	"github.com/newton-miku/webook/webook-be/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail  = dao.ErrUserDuplicateEmail
	ErrUserNotFound        = dao.ErrUserNotFound
	ErrUserProfileNotFound = dao.ErrUserProfileNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func (r *UserRepository) UpdateProfile(ctx *gin.Context, u domain.UserProfile) error {
	return r.dao.UpdateProfile(ctx, dao.UserProfile{
		Id:          u.Id,
		UID:         u.UID,
		Nickname:    u.Nickname,
		Email:       u.Email,
		PhoneNumber: u.Phone,
		Summary:     u.Summary,
		Birthday:    u.Birthday,
	})
}

func (r *UserRepository) FindProfileByID(ctx context.Context, uid int64) (domain.UserProfile, error) {
	u, err := r.dao.FindProfileByID(ctx, uid)
	if err != nil {
		return domain.UserProfile{}, err
	}
	return domain.UserProfile{
		Id:       u.Id,
		UID:      u.UID,
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: u.Birthday,
		Phone:    u.PhoneNumber,
		Summary:  u.Summary,
	}, nil
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: []byte(u.Password),
	}, nil
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: string(u.Password),
	})
}
