package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"github.com/hxbaby/biz-service/internal/model"
	"github.com/hxbaby/biz-service/internal/repository"
	"github.com/hxbaby/biz-service/pkg/jwt"
)

type CustomerService struct {
	repo      *repository.CustomerRepo
	jwtSecret string
}

func NewCustomerService(repo *repository.CustomerRepo, jwtSecret string) *CustomerService {
	return &CustomerService{repo: repo, jwtSecret: jwtSecret}
}

func (s *CustomerService) Register(phone, password, name string) (*model.Customer, string, error) {
	existing, _ := s.repo.FindByPhone(phone)
	if existing != nil {
		return nil, "", errors.New("手机号已注册")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	c := &model.Customer{
		Phone:        phone,
		PasswordHash: string(hash),
		Name:         name,
		Plan:         "free",
		MaxProjects:  1,
		Status:       "active",
	}
	if err := s.repo.Create(c); err != nil {
		return nil, "", err
	}

	token, err := jwt.GenerateToken(c.ID, 0, "customer", s.jwtSecret)
	return c, token, err
}

func (s *CustomerService) Login(phone, password string) (*model.Customer, string, error) {
	c, err := s.repo.FindByPhone(phone)
	if err != nil {
		return nil, "", errors.New("手机号未注册")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("密码错误")
	}
	token, err := jwt.GenerateToken(c.ID, 0, "customer", s.jwtSecret)
	return c, token, err
}
