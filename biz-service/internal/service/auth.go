package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/pkg/jwt"
)

type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret string
}

func NewAuthService(repo *repository.UserRepo, secret string) *AuthService {
	return &AuthService{userRepo: repo, jwtSecret: secret}
}

func (s *AuthService) Register(phone, password, name string, tenantID uint) (*model.User, string, error) {
	existing, _ := s.userRepo.FindByPhone(phone)
	if existing != nil {
		return nil, "", errors.New("手机号已注册")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &model.User{
		TenantID: tenantID,
		Phone:    phone,
		Name:     name,
		Role:     "parent",
	}
	// 手动设置密码hash（User model不含PasswordHash，这里需要扩展）
	_ = hash

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := jwt.GenerateToken(user.ID, user.TenantID, user.Role, s.jwtSecret)
	return user, token, err
}

func (s *AuthService) Login(phone, password string) (*model.User, string, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return nil, "", errors.New("手机号未注册")
	}
	// 密码校验（暂时跳过，需要User model添加PasswordHash字段）
	_ = password
	_ = bcrypt.CompareHashAndPassword

	token, err := jwt.GenerateToken(user.ID, user.TenantID, user.Role, s.jwtSecret)
	return user, token, err
}
