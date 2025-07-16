package service

import (
	"context"
	"errors"

	"github.com/newton-miku/webook/webook-be/internal/domain"
	"github.com/newton-miku/webook/webook-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")
	ErrUserNotFound          = repository.ErrUserNotFound
)

type UserService struct {
	repo *repository.UserRepository
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

func (svc *UserService) Login(ctx context.Context, user domain.User) error {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		// 用户不存在
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrInvalidUserOrPassword
		}
		return errors.New("登录时系统发生异常")
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password))
	if err != nil {
		// 密码不对
		return ErrInvalidUserOrPassword
	}
	return nil
}
