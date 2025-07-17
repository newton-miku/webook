package service

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/newton-miku/webook/webook-be/internal/domain"
	"github.com/newton-miku/webook/webook-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")
	ErrProfileNotFound       = repository.ErrUserProfileNotFound
	ErrUserNotFound          = repository.ErrUserNotFound
)

type UserService struct {
	repo *repository.UserRepository
}

func (svc *UserService) UpdateProfile(ctx *gin.Context, u domain.UserProfile) error {
	return svc.repo.UpdateProfile(ctx, u)
}

func (svc *UserService) Profile(ctx context.Context, i int64) (domain.UserProfile, error) {
	// 先找用户
	u, err := svc.repo.FindProfileByID(ctx, i)
	if err != nil {
		// 用户档案不存在
		if errors.Is(err, ErrProfileNotFound) {
			return domain.UserProfile{}, ErrProfileNotFound
		}
		return domain.UserProfile{}, errors.New("查询档案时发生异常")
	}
	return u, nil
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = hash
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		// 用户不存在
		if errors.Is(err, ErrUserNotFound) {
			return domain.User{}, ErrInvalidUserOrPassword
		}
		return domain.User{}, errors.New("登录时系统发生异常")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		// 密码不对
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
