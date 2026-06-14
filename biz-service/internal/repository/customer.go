package repository

import (
	"github.com/hxbaby/biz-service/internal/model"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(db *gorm.DB) *CustomerRepo {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) Create(c *model.Customer) error {
	return r.db.Create(c).Error
}

func (r *CustomerRepo) FindByPhone(phone string) (*model.Customer, error) {
	var c model.Customer
	err := r.db.Where("phone = ?", phone).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) FindByID(id uint) (*model.Customer, error) {
	var c model.Customer
	err := r.db.First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CustomerRepo) Update(c *model.Customer) error {
	return r.db.Save(c).Error
}
