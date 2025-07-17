package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail   = errors.New("邮箱已被注册")
	ErrUserProfileDuplicate = errors.New("档案重复")
	ErrUserNotFound         = gorm.ErrRecordNotFound
	ErrUserProfileNotFound  = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	return u, err
}

func (dao *UserDAO) FindByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, "id = ?", id).Error
	return u, err
}

func (dao *UserDAO) FindProfileByID(ctx context.Context, uid int64) (UserProfile, error) {
	var u UserProfile
	err := dao.db.WithContext(ctx).First(&u, "UID = ?", uid).Error
	if errors.Is(err, ErrUserProfileNotFound) {
		u.UID = uid
		err = dao.InsertProfile(ctx, u)
		if err != nil {
			return UserProfile{}, err
		}
		return dao.FindProfileByID(ctx, uid)
	}
	return u, err
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// 对应数据库表结构
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	// 时间统一为UTC+0
	// 创建时间，时间戳
	Ctime int64
	// 更新时间，时间戳
	Utime int64
}

type UserProfile struct {
	Id          int64 `gorm:"primaryKey,autoIncrement"`
	UID         int64 `gorm:"uniqueIndex"`
	Nickname    string
	Email       string
	PhoneNumber string
	Birthday    string
	Summary     string

	Ctime int64
	Utime int64
}

func (dao *UserDAO) UpdateProfile(ctx context.Context, up UserProfile) error {
	now := time.Now().Unix()
	u := UserProfile{}
	err := dao.db.WithContext(ctx).First(&u, "UID = ?", up.UID).Error
	if err != nil {
		return err
	}
	u.Birthday = up.Birthday
	u.Nickname = up.Nickname
	u.Summary = up.Summary
	u.Utime = now
	dao.db.WithContext(ctx).Save(&u)
	return nil
}
func (dao *UserDAO) InsertProfile(ctx context.Context, up UserProfile) error {
	now := time.Now().Unix()
	up.Ctime = now
	up.Utime = now
	up.Birthday = time.Now().Format("2006-01-02")
	if up.UID != 0 {
		// 没有邮箱信息则查询数据
		if up.Email == "" {
			user, err := dao.FindByID(ctx, up.UID)
			if err != nil {
				return err
			}
			up.Email = user.Email
		}
	}
	err := dao.db.WithContext(ctx).Create(&up).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictErr uint16 = 1062
		if mysqlErr.Number == uniqueConflictErr {
			return ErrUserDuplicateEmail
		}
	}
	return nil
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().Unix()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictErr uint16 = 1062
		if mysqlErr.Number == uniqueConflictErr {
			return ErrUserDuplicateEmail
		}
	}
	return dao.InsertProfile(ctx, UserProfile{UID: u.Id, Email: u.Email})
}
